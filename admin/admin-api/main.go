package main

import (
	// Commnuity packages
	"fmt"
	// Internal packages
	config "admin-api/configs"
	database "admin-api/configs/databases"
	redisConfig "admin-api/configs/redis"
	"admin-api/handler"
	"admin-api/internal/admin/websocket"
	"admin-api/pkg/logs"
	"admin-api/pkg/translate"
	"admin-api/router"
)

func main() {
	// Initalize config
	app_configs := config.NewConfig()

	// Initalize Databse e.g '20+connections opened' are ready to use
	db_pool := database.GetDB()
	defer db_pool.Close()

	// Initalize Redis
	fmt.Println("Connecting to Redis...")
	rdb := redisConfig.NewRedisClient()
	fmt.Println("Redis connected")

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
