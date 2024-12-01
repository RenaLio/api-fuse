package main

import "fmt"

func main() {
	fmt.Println("gogo")
	//lb, err := NewClientLoadBalancer("config.json")
	//if err != nil {
	//	logrus.Fatalf("Failed to create load balancer: %v", err)
	//}
	//
	//r := gin.Default()
	//
	//// 增加CORS中间件
	//r.Use(cors.Default())
	//
	//// 增加限流中间件
	//rate, err := limiter.NewRateFromFormatted("100-H")
	//if err != nil {
	//	logrus.Fatalf("Failed to create rate limiter: %v", err)
	//}
	//store := memory.NewStore()
	//middleware := stdlib.NewMiddleware(limiter.New(store, rate))
	//r.Use(func(c *gin.Context) {
	//	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		c.Next()
	//	}))
	//	handler.ServeHTTP(c.Writer, c.Request)
	//})
	//
	//r.POST("/v1/completions", func(c *gin.Context) {
	//	lb.HandleRequest(c.Writer, c.Request)
	//})
	//
	//r.Run() // 默认监听8080端口
}
