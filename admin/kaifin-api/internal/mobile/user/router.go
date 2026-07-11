package user

import (

	// Commnuity pacagkes
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type UserRoute struct {
	handler *UserHandler
}

func NewUserRoute(a *fiber.App, db *sqlx.DB, rdb *redis.Client) *UserRoute {
	h := NewUserHandler(db, rdb)

	v1 := a.Group("/api/v1/admin")
	users := v1.Group("/users")

	users.Get("/", h.Show)
	users.Get("/:id", h.ShowOne)
	users.Post("/create", h.Create)
	users.Put("/update/:id", h.Update)
	users.Delete("/delete/:id", h.Delete)

	users.Get("/form/create", h.GetUserFormCreate)
	users.Get("/form/update/:id", h.GetUserFormUpdate)

	return &UserRoute{
		handler: h,
	}
}
