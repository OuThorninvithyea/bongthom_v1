package share

import (

	// Community pacakge
	"time"
)

type UserFromCreate struct {
	FirstName   string      `json:"first_name" db:"first_name"`
	LastName    string      `json:"last_name" db:"last_name"`
	Nationality Nationality `json:"nationality" db:"nationality"`
}

type UserFormUpdate struct {
	ID            int64       `json:"id" db:"id"`
	FirstName     string      `json:"first_name" db:"first_name"`
	LastName      string      `json:"last_name" db:"last_name"`
	NationalityID int         `json:"national_id" db:"nationality_id"`
	Nationality   Nationality `json:"nationality" db:"nationality"`
}

type Nationality struct {
	ID    int64  `json:"id" db:"id"`
	Value string `json:"value" db:"value"`
}

type UserContext struct {
	UserID       int64     `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name" db:"user_name"`
	LoginSession string    `json:"login_session" db:"login_session"`
	Exp          time.Time `json:"exp" db:"exp"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	IP           string    `json:"ip" db:"ip"`
	RoleID       int       `json:"role_id" db:"role_id"`
}

type Paging struct {
	Page    int `json:"page" query:"page" validate:"required,min=1"`
	Perpage int `json:"per_page" query:"per_page" validate:"required,min=1"`
}

// Sort products price and asc to desc
type Sort struct {
	Property  string `json:"property" validate:"required"`                 // column database tables
	Direction string `json:"direction" validate:"required,oneof=asc desc"` // keyword database  ORDER BY asc desc
}

// Fitler  or Search products
type Filter struct {
	Property string      `json:"property" validate:"required" query:"property"` // column database tables
	Value    interface{} `json:"value" validate:"required" query:"value"`       // value data in row table
}

type Status struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// StatusData contains predefined status values.
var StatusData = []Status{
	{ID: 1, Name: "Active"},
	{ID: 2, Name: "Inactive"},
	{ID: 3, Name: "Suspended"},
	{ID: 4, Name: "Deleted"},
}
