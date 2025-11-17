package main

import (
	"log"
	"lottery_system/handlers"
	"lottery_system/storage"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	store := storage.NewStorage()
	if err := store.Load(); err != nil {
		log.Fatalf("加载存储失败: %v", err)
	}
	handlers.Init(store)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := store.CheckAndAutoDraw(); err != nil {
				log.Printf("自动开奖检查失败: %v", err)
			}
		}
	}()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	api := r.Group("/api")
	{
		api.POST("/participants", handlers.AddParticipant)
		api.GET("/participants", handlers.GetParticipants)
		api.POST("/draw", handlers.DrawWinner)
		api.GET("/winner", handlers.GetWinner)
		api.POST("/reset", handlers.Reset)
		api.POST("/prizes", handlers.AddPrize)
		api.GET("/prizes", handlers.GetPrizes)
		api.POST("/draw-time", handlers.SetDrawTime)
		api.GET("/draw-time", handlers.GetDrawTime)
	}
	log.Println("正在运行")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
