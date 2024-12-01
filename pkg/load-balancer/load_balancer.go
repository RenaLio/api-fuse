package load_balancer

import (
	v1 "awesomeGoProject/api/v1"
	"awesomeGoProject/pkg/adapter"
	"fmt"
	"net/http"
	"strconv"
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
	//data, err := ioutil.ReadFile(configPath)
	//if err != nil {
	//	logrus.Fatalf("Failed to read config file: %v", err)
	//	return nil, err
	//}
	//
	//var config struct {
	//	ModelMapping map[string][]string `json:"model_mapping"`
	//	Providers    []BaseProvider      `json:"providers"`
	//}
	//if err := json.Unmarshal(data, &config); err != nil {
	//	logrus.Fatalf("Failed to unmarshal config: %v", err)
	//	return nil, err
	//}
	//
	//providers := make([]ClientProvider, len(config.Providers))
	//for i := range config.Providers {
	//	config.Providers[i].Active = make(map[string]bool)
	//	if len(config.Providers[i].Models) == 0 {
	//		err := config.Providers[i].FetchModels()
	//		if err != nil {
	//			logrus.Errorf("Failed to fetch models for provider %s: %v", config.Providers[i].Name, err)
	//			return nil, err
	//		}
	//	}
	//	providers[i] = &config.Providers[i]
	//}
	//
	//lb := &ClientLoadBalancer{
	//	Providers:           providers,
	//	ModelCountMap: make(map[string]int),
	//	ModelMapping:        config.ModelMapping,
	//	Index:               0,
	//}
	//
	//go lb.startHealthCheck()

	//return nil, nil
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
		key := srv.GetAPIKey()
		keyId := "不足8位，id->" + strconv.Itoa(int(srv.GetId()))
		if len(key) >= 8 {
			keyId = key[len(key)-8:]
		}
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
				time.Sleep(60 * 4 * time.Minute)
				task <- "start"
			}
		}
	}()
	task <- "start"
}

func (lb *LoadBalancer) DoHealthCheck() {
	list := make([]string, 0)
	for _, model := range lb.CheckModels {
		if val, ok := lb.ModelMapping[model]; ok {
			list = append(list, val...)
		} else {
			list = append(list, model)
		}
	}
	for _, model := range list {
		fmt.Println("开始检测任务-->", model)
		for _, provider := range lb.Providers {
			if provider.HasModel(model) {
				ok := provider.HealthCheck(model)
				if ok {
					if !provider.IsModelActive(model) {
						provider.SetModelActive(model, true)
						lb.Mutex.Lock()
						if _, exist := lb.ModelMapping[model]; exist {
							lb.ModelCountMap[model]++
						}
						lb.Mutex.Unlock()
					}
				} else {
					if provider.IsModelActive(model) {
						provider.SetModelActive(model, false)
						lb.Mutex.Lock()
						if _, exist := lb.ModelMapping[model]; exist {
							lb.ModelCountMap[model]--
						}
						lb.Mutex.Unlock()
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
