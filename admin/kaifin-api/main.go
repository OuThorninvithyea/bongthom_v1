package main

import (
	// Commnuity packages
	"fmt"
	
	// Internal packages
	config "kaifin-api/configs"
	database "kaifin-api/configs/databases"
	 "kaifin-api/configs/redis"
	"kaifin-api/handler"
	"kaifin-api/internal/mobile/websocket"
	"kaifin-api/pkg/logs"
	"kaifin-api/pkg/translate"
	"kaifin-api/router"
)

func main() {
	// Initalize config
	app_configs := config.NewConfig()

	// Initalize Databse e.g '20+connections opened' are ready to use
	db_pool := database.GetDB()
	defer db_pool.Close()

	// Initalize Redis
	rdb := redis.NewRedisClient()

	// Initalize Websocket Manager
	ws_manager := websocket.NewWebSocketManager()

	// SetupRouter
	app := router.New()

	// Initialize translate
	if err := translate.Init(); err != nil {
		logs.NewCustomLog("FailedInitializeI18n", err.Err.Error(), "error")
	}


	// Initalize service handlers e.g 'admin', 'front'
	h := handler.NewServiceHandlers(app, db_pool, rdb, ws_manager)
	_ = h

	// Start Http Server (entering even loop)
	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
}
