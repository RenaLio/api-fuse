package main

import (
	v1 "awesomeGoProject/api/v1"
	"awesomeGoProject/pkg/adapter"
	config2 "awesomeGoProject/pkg/config"
	load_balancer "awesomeGoProject/pkg/load-balancer"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
	"time"
)

func main() {
	fmt.Println("加载配置中...")
	config, err := config2.Load("./config/config.yaml")
	//fmt.Println(config)
	if err != nil {
		panic(err)
	}
	filePath := "example.txt"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	file.WriteString(fmt.Sprintf("加载配置成功:%s\n", time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println("加载配置成功")
	adapters := make([]adapter.Adapter, 0)
	adCh := make(chan adapter.Adapter)
	go func() {
		for i := range adCh {
			//fmt.Println("获取模型成功", i.GetName(), i.GetId(), i.GetModels())
			fmt.Fprintf(file, "获取模型成功: %s\t\t %d\t\t %v\n", i.GetName(), i.GetId(), i.GetModels())
			adapters = append(adapters, i)
		}
	}()
	wg := sync.WaitGroup{}
	var id uint64
	for i := range config.Providers {
		for j := range config.Providers[i].Data {
			data := config.Providers[i].Data[j]
			for c := range data.APIKeys {
				id++
				wg.Add(1)
				go func(i, j, c int, id uint64) {
					defer wg.Done()
					//fmt.Println(data.Name, data.URL, data.APIKeys[c])
					ad := NewAdapter(config.Providers[i].Name, id, data.Name, data.URL, data.APIKeys[c], data.Models)
					err := ad.FetchModels()
					if err != nil {
						fmt.Println("获取模型错误", i, j, c, id, data.APIKeys[c], err)
						fmt.Fprintf(file, "获取模型错误: %s\t\t %d\t\t %s\t\t %v\n", ad.GetName(), ad.GetId(), data.APIKeys[c], err.Error())
					} else {
						adCh <- ad
					}
				}(i, j, c, id)
			}
		}
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
	fmt.Println("获取模型任务完成")
	close(adCh)
	file.Close()
	loadBalance := load_balancer.NewClientLoadBalancer(adapters, config.ModelMapping, config.CheckModels)
	fmt.Println("构建负载均衡器成功")
	go func() {
		fmt.Println("定时检查任务已创建")
		loadBalance.StartHealthCheck()
	}()
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/v1/chat/completions", HandleChat(loadBalance))
	r.GET("/v1/models", HandleModels(loadBalance))

	r.Run() // 默认监听8080端口
}

func HandleChat(lb *load_balancer.LoadBalancer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req v1.ChatCompletionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		lb.HandleRequest(ctx.Writer, ctx.Request, &req)
	}
}

func HandleModels(lb *load_balancer.LoadBalancer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		models := lb.GetActiveModels()
		resp := v1.ModelResp{
			Object: "list",
		}
		data := make([]v1.Model, 0)
		for _, v := range models {
			data = append(data, v1.Model{
				ID:      v,
				Object:  "model",
				Created: 0,
				OwnedBy: "owner",
			})
		}
		resp.Data = data
		ctx.JSON(200, resp)
	}
}
