package main

import (
	// Commnuity packages 
	"fmt" 
	
	// Internal packages
	config "admin-api/configs"
	"admin-api/configs/databases"
	"admin-api/handler"
	"admin-api/internal/admin/websocket"
	"admin-api/router"
)

func main() {
	// Initalize config
	app_configs := config.NewConfig()

	// Initalize Databse e.g '20+connections opened' are ready to use 
	db_pool := database.GetDB()
	defer db_pool.Close()

	// Initalize Websocket Manager 
	ws_manager := websocket.NewWebSocketManager()
	
	// SetupRouter 
	app := router.New()
	
	// Initalize service handlers e.g 'admin', 'front'
	h := handler.NewServiceHandlers(app, db_pool, ws_manager)
	_ = h
	// Start Http Server (entering even loop)
	app.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
}
