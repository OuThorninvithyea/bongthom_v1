package user

import (

	// Community pacakges
	sql "admin-api/pkg/sql"
	"admin-api/pkg/logs"
	"admin-api/pkg/share"
	"admin-api/pkg/utls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// User maps to tbl_users
type User struct {
	ID           int64      `json:"id" db:"id"`
	FirstName    *string    `json:"first_name" db:"first_name"`
	LastName     *string    `json:"last_name" db:"last_name"`
	UserName     string     `json:"user_name" db:"user_name"`
	Email        *string    `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"`
	RoleName     string     `json:"role_name" db:"role_name"`
	RoleID       int        `json:"role_id" db:"role_id"`
	IsAdmin      bool       `json:"is_admin" db:"is_admin"`
	LoginSession *string    `json:"-" db:"login_session"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CurrencyID   *int       `json:"currency_id" db:"currency_id"`
	LanguageID   *int       `json:"language_id" db:"language_id"`
	StatusID     *int       `json:"status_id" db:"status_id"`
	Order        *int       `json:"order" db:"order"`
	CreatedBy    *int64     `json:"-" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy    *int64     `json:"-" db:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	DeletedBy    *int64     `json:"-" db:"deleted_by"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

type UserShowRequest struct {
	PageOption share.Paging   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts      []share.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters    []share.Filter `json:"filters.omitempty" query:"filters"`
	Search     string         `json:"q,omitempty" query:"q"`
	CurrencyID int            `json:"currency_id,omitempty" query:"currency_id"`
}

func (u *UserShowRequest) bind(c fiber.Ctx, v *utls.Validator) error {
	if err := c.Bind().Query(u); err != nil {
		return err
	}

	for i := range u.Filters {
		// request takes &filters[index][value]= int,
		value := c.Query(fmt.Sprintf("filters[%d][value]", i))
		if intValue, err := strconv.Atoi(value); err == nil {
			u.Filters[i].Value = intValue
		} else if boolValue, err := strconv.ParseBool(value); err == nil {
			u.Filters[i].Value = boolValue
		} else {
			u.Filters[i].Value = value
		}
	}
	if u.Search == "" {
		u.Search = c.Query("q")
	}
	if u.CurrencyID == 0 {
		if v := strings.TrimSpace(c.Query("currency_id")); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				u.CurrencyID = n
			}
		}
	}

	if err := v.Validate(u); err != nil {
		return err
	}
	return nil
}

// CreateUserRequest — fields needed to create a user
type UserCreateRequest struct {
	UserName  string  `json:"user_name" validate:"required,min=4"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	RoleID    int     `json:"role_id" validate:"required"`
	RoleName  string  `json:"role_name" validate:"required"`
	CreatedBy *int64  `json:"-" db:"created_by"`
}

func (u *User) New(userReq *UserCreateRequest, uCtx *share.UserContext, db_pool *sqlx.DB) error {
	if uCtx.RoleID > userReq.RoleID {
		return fmt.Errorf("failed you role can not create this user")
	}

	login_session, err := uuid.NewV7()
	if err != nil {
		logs.NewCustomLog("get_uuid_failed", err.Error(), "error")
		return err
	}
	sessionString := login_session.String()

	app_timezone := os.Getenv("TIME_ZONE")
	location, err := time.LoadLocation(app_timezone)
	if err != nil {
		return fmt.Errorf("failed to load location: %w", err)
	}
	local_now := time.Now().In(location)

	is_username, err := sql.IsExits("tbl_users", "user_name", userReq.UserName, db_pool)
	if err != nil {
		return err
	}
	if is_username {
		return fmt.Errorf("username:`%s` already exists", userReq.UserName)
	}

	createdByID, err := sql.GetUserIdByField("tbl_users", "user_name", uCtx.UserName, db_pool)
	if err != nil {
		return err
	}

	orderValue, err := sql.GetSeqNextVal("tbl_users_id_seq", db_pool)
	if err != nil {
		return fmt.Errorf("failed to generate order value: %w", err)
	}

	u.FirstName = userReq.FirstName
	u.LastName = userReq.LastName
	u.UserName = strings.ToUpper(strings.TrimSpace(userReq.UserName))
	u.Password = userReq.Password
	u.Email = userReq.Email
	u.LoginSession = &sessionString
	statusID := 1
	u.StatusID = &statusID
	u.Order = orderValue
	createdByInt64 := int64(*createdByID)
	u.CreatedBy = &createdByInt64
	u.CreatedAt = local_now
	u.RoleID = userReq.RoleID
	u.RoleName = userReq.RoleName

	return nil
}

// UpdateUserRequest — partial update
type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	RoleID    *int    `json:"role_id"`
	RoleName  *string `json:"role_name"`
	StatusID  *int    `json:"status_id"`
}

// UserResponse — safe fields returned to client
type UserResponse struct {
	Users []User `json:"users"`
	Total int
}

func (r *UserCreateRequest) bind(c fiber.Ctx, v *utls.Validator) error {
	if err := c.Bind().Body(&r); err != nil {
		return err
	}
	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

type PagingRequest struct {
	Page    int `query:"page" validate:"min=1"`
	PerPage int `query:"per_page" validate:"min=1,max=100"`
}
