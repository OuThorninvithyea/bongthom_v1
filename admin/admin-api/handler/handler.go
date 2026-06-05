package handler

import (
	// Community packages
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// Internal packages
	"admin-api/internal/admin/auth"
	"admin-api/internal/admin/user"
	"admin-api/internal/admin/websocket"
	"admin-api/pkg/middlewares"
)

type ServiceHandlers struct {
	Admin *AdminService
	Front *FrontService
}

func NewServiceHandlers(a *fiber.App, db *sqlx.DB, rdb *redis.Client, wsmgr *websocket.WebSocketManager) *ServiceHandlers {
	return &ServiceHandlers{
		Admin: NewAdminService(a, db, rdb, wsmgr),
		Front: NewFrontService(a),
	}
}

type AdminService struct {
	Auth *auth.AuthRoute
	User *user.UserRoute
}

func NewAdminService(a *fiber.App, db *sqlx.DB, rdb *redis.Client, wsmgr *websocket.WebSocketManager) *AdminService {
	authRoute := auth.NewAuthRoute(a, db, rdb) //
	middlewares.NewJwtMiddleware(a, db, rdb)   // checking HTTP incoming request, { Header, body } as plain TEXT
	userRoute := user.NewUserRoute(a, db, rdb) //
	return &AdminService{
		Auth: authRoute,
		User: userRoute,
	}
}

type FrontService struct {
}

func NewFrontService(a *fiber.App) *FrontService {
	return &FrontService{}
}
