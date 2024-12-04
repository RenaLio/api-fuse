package load_balancer

import (
	v1 "awesomeGoProject/api/v1"
	"awesomeGoProject/pkg/adapter"
	"fmt"
	"golang.org/x/exp/rand"
	"net/http"
	"sync"
	"time"
)

const nTimes = 8

type Server = adapter.Adapter

type LoadBalancer struct {
	Providers     []Server
	ModelCountMap map[string]int      // 实际模型名称-> 可用的适配器数量
	ModelMapping  map[string][]string // 模型别名 -> 实际模型名称
	Index         int
	Mutex         sync.RWMutex
	CheckModels   []string
}

func NewClientLoadBalancer(providers []Server, modelMapping map[string][]string, checkModes []string) *LoadBalancer {
	modelMapCount := make(map[string]int)
	for i := range providers {
		for _, modelId := range providers[i].GetModels() {
			modelMapCount[modelId]++
		}
	}
	for modelId, _ := range modelMapping {
		modelMapping[modelId] = append(modelMapping[modelId], modelId)
	}

	return &LoadBalancer{
		Providers:     providers,
		ModelCountMap: modelMapCount,
		ModelMapping:  modelMapping,
		CheckModels:   checkModes,
	}
}

// NextProvider 返回下一个可用的Provider，传入的是原始请求的输入model
func (lb *LoadBalancer) NextProvider(modelId string) (Server, string) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	models, ok := lb.ModelMapping[modelId]
	if !ok {
		models = []string{modelId}
	}

	for i := 0; i < len(lb.Providers); i++ {
		provider := lb.Providers[lb.Index]
		for j := range models {
			if provider.IsModelActive(models[j]) {
				lb.Index = (lb.Index + 1) % len(lb.Providers)
				return provider, models[j]
			}
		}
		lb.Index = (lb.Index + 1) % len(lb.Providers)
	}

	return nil, ""
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request, req *v1.ChatCompletionRequest) {
	//var srv Server
	var err error

	var count int
	lb.Mutex.RLock()
	if val, ok := lb.ModelMapping[req.Model]; ok {
		for _, v := range val {
			count += lb.ModelCountMap[v]
		}
	} else {
		count = lb.ModelCountMap[req.Model]
	}
	lb.Mutex.RUnlock()
	if count < 1 {
		http.Error(w, "No active providers available for the specified model", http.StatusServiceUnavailable)
		return
	}
	model := req.Model
	for _ = range nTimes {
		srv, modelId := lb.NextProvider(model)
		if srv == nil {
			break
		}
		req.Model = modelId
		err = srv.Chat(w, r, req)
		keyId := GetKeyId(srv)
		if err == nil {
			fmt.Printf("请求成功|| 请求模型:%s\t\t 实际模型:%s\t\t 服务商名称：%s\t\t ID:%v\t\t Key:%s\n", model, modelId, srv.GetName(), srv.GetId(), keyId)
			return
		} else {
			fmt.Printf("失败重试|| 请求模型:%s\t\t 实际模型:%s\t\t 服务商名称：%s\t\t ID:%v\t\t Key:%s\n", model, modelId, srv.GetName(), srv.GetId(), keyId)
			lb.BanProviderMode(model)
		}
	}
	http.Error(w, "No active providers available for the specified model", http.StatusServiceUnavailable)
}

func (lb *LoadBalancer) BanProviderMode(modelId string) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()
	count := lb.ModelCountMap[modelId]
	if count > 0 {
		lb.ModelCountMap[modelId] = count - 1
	}
}

func (lb *LoadBalancer) StartHealthCheck() {
	task := make(chan string, 1)
	go func() {
		for {
			select {
			case _ = <-task:
				fmt.Println("一轮健康测试已开始")
				lb.DoHealthCheck()
				fmt.Println("一轮健康测试已结束")
				num := rand.Intn(60)
				total := 60*60 + num
				time.Sleep(time.Second * time.Duration(total))
				task <- "start"
			}
		}
	}()
	task <- "start"
}

func (lb *LoadBalancer) DoHealthCheck() {
	list := make([]string, 0)
	tempList := make([]string, 0)
	modelSet := make(map[string]bool)
	for _, model := range lb.CheckModels {
		if val, ok := lb.ModelMapping[model]; ok {
			tempList = append(tempList, val...)
		} else {
			tempList = append(tempList, model)
		}
	}
	for _, model := range tempList {
		if _, ok := modelSet[model]; !ok {
			list = append(list, model)
			modelSet[model] = true
		}
	}
	for _, model := range list {
		fmt.Println("开始检测任务-->", model)
		for _, provider := range lb.Providers {
			if provider.HasModel(model) {
				ok := provider.HealthCheck(model)
				if ok {
					if !provider.IsModelActive(model) {
						fmt.Printf("健康检测||失败->成功|| 请求模型:%s\t\t 服务商名称:%s\t\t  Key:%s\n", model, provider.GetName(), GetKeyId(provider))
						provider.SetModelActive(model, true)
						lb.Mutex.Lock()
						if _, exist := lb.ModelMapping[model]; exist {
							lb.ModelCountMap[model]++
						}
						lb.Mutex.Unlock()
					} else {
						fmt.Printf("健康检测||成功->成功|| 请求模型:%s\t\t 服务商名称:%s\t\t  Key:%s\n", model, provider.GetName(), GetKeyId(provider))
					}
				} else {
					if provider.IsModelActive(model) {
						fmt.Printf("健康检测||成功->失败|| 请求模型:%s\t\t 服务商名称:%s\t\t  Key:%s\n", model, provider.GetName(), GetKeyId(provider))
						provider.SetModelActive(model, false)
						lb.Mutex.Lock()
						if _, exist := lb.ModelMapping[model]; exist {
							if lb.ModelCountMap[model] > 0 {
								lb.ModelCountMap[model]--
							}
						}
						lb.Mutex.Unlock()
					} else {
						fmt.Printf("健康检测||失败->失败|| 请求模型:%s\t\t 服务商名称:%s\t\t  Key:%s\n", model, provider.GetName(), GetKeyId(provider))
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
		fmt.Println("检测任务完成-->", model)
	}
}

func (lb *LoadBalancer) GetActiveModels() []string {
	lb.Mutex.RLock()
	defer lb.Mutex.RUnlock()
	models := make([]string, 0)
	for model, count := range lb.ModelCountMap {
		if count > 0 {
			models = append(models, model)
		}
	}
	for model, list := range lb.ModelMapping {
		flag := false
		for _, v := range list {
			if lb.ModelCountMap[v] > 0 {
				flag = true
				break
			}
		}
		if flag {
			models = append(models, model)
		}
	}
	return models
}
