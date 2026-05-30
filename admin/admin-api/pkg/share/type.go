package share

import (

	// Community pacakge
	"time"
)

type UserContext struct {
	UserID       float64   `json:"user_id"`
	UserName     string    `json:"user_name"`
	LoginSession string    `json:"login_session"`
	Exp          time.Time `json:"exp"`
	UserAgent    string    `json:"user_agent"`
	Ip           string    `json:"ip"`
	RoleID       int       `json:"role_id"`
}

type Paging struct {
	Page    int `json:"page" query:"page" validate:"required,min=1"`
	Perpage int `json:"per_page" query:"per_page" validate:"required,min=1"`
}

// Sort products price and asc to desc 
type Sort struct {
	Property  string `json:"property" validate:"required"`  // column database tables
	Direction string `json:"direction" validate:"required,oneof=asc desc"` // keyword database  ORDER BY asc desc
}

// Fitler  or Search products 
type Filter struct {
	Property string      `json:"property" validate:"required" query:"property"` // column database tables
	Value    interface{} `json:"value" validate:"required" query:"value"` // value data in row table 
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StatusData contains predefined status values.
var StatusData = []Status{
	{ID: 1, Name: "Active"},
	{ID: 2, Name: "Inactive"},
	{ID: 3, Name: "Suspended"},
	{ID: 4, Name: "Deleted"},
}
