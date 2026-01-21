package main

import (
	"log"
	"net/http"

	"github.com/sudankdk/cga/configs"
	"github.com/sudankdk/cga/internal/handler"
	my_redis "github.com/sudankdk/cga/internal/redis"
	"github.com/sudankdk/cga/internal/router"
	"github.com/sudankdk/cga/internal/service"
)

func main() {
	cfg := configs.Load()
	my_redis.InitRedis(cfg.Redis)
	notificationService := service.NewNotificationService()
	handler := handler.NewHandler(*notificationService)

	routers := router.NewRouter(handler) //expects handler
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: routers,
	}
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
