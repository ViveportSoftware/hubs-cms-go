package main

import (
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/jobs"
	"hubs-cms-go/logger"
	"hubs-cms-go/router"
	"log"
	"net/http"
	"time"
)

// Setup application initialization
func Setup() {
	config.Setup()
	logger.Setup(config.EnvVariable.LogLevel)
	cache.Setup()
	client.Setup()
	router.Setup()
	jobs.Setup()
}

// @title Package Management Service Swagger
// @version 1.0.0
// @description this service is used to get packages info

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	Setup()

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.EnvVariable.Port),
		Handler:      router.Router,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
	}

}
