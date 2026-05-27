package handler

import (
	// Community packages
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"

	// Internal packages
	"admin-api/internal/admin/auth"
	"admin-api/internal/admin/websocket"
)

type ServiceHandlers struct {
	Admin *AdminService
	Front *FrontService
}

func NewServiceHandlers(a *fiber.App, db *sqlx.DB, wsmgr *websocket.WebSocketManager) *ServiceHandlers {
	return &ServiceHandlers{
		Admin: NewAdminService(a, db, wsmgr),
		Front: NewFrontService(a),
	}
}

type AdminService struct {
	Auth *auth.AuthRoute // moduel of its project
}

func NewAdminService(a *fiber.App, db *sqlx.DB, wsmgr *websocket.WebSocketManager) *AdminService {
	authRoute := auth.NewAuthRoute(a, db)
	return &AdminService{
		Auth: authRoute,
	}
}

type FrontService struct {
}

func NewFrontService(a *fiber.App) *FrontService {
	return &FrontService{}
}
