package user

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	error_responses "admin-api/pkg/responses"
	"admin-api/pkg/share"
)

type UserService interface {
	Show(UserShowRequest) (*UserResponse, *error_responses.ErrorResponse)
	ShowOne(id int64) (*UserResponse, *error_responses.ErrorResponse)
	Create(req *UserCreateRequest, uCtx share.UserContext) *error_responses.ErrorResponse
	Update(id int64, req *UpdateUserRequest, uCtx share.UserContext) (*User, *error_responses.ErrorResponse)
	Delete(id int64, uCtx share.UserContext) *error_responses.ErrorResponse
	GetCreateForm() any
	GetUpdateForm(id int64) (*UserResponse, *error_responses.ErrorResponse)
}

type UserServiceImpl struct {
	Repo  UserRepo
	Redis *redis.Client
	DB    *sqlx.DB
}

func NewUserServiceImpl(db *sqlx.DB, rdb *redis.Client) *UserServiceImpl {
	return &UserServiceImpl{
		Repo:  NewUserRepoImpl(db, rdb),
		Redis: rdb,
		DB:    db,
	}
}

func (s *UserServiceImpl) Show(userRequest UserShowRequest) (*UserResponse, *error_responses.ErrorResponse) {
	return s.Repo.Show(userRequest)
}

func (s *UserServiceImpl) ShowOne(id int64) (*UserResponse, *error_responses.ErrorResponse) {
	return s.Repo.ShowOne(id)
}

func (s *UserServiceImpl) Create(req *UserCreateRequest, uCtx share.UserContext) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}

	var user User
	if err := user.New(req, &uCtx, s.DB); err != nil {
		return msg.NewErrorResponse("user_create_failed", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return msg.NewErrorResponse("password_hash_failed", err)
	}
	user.Password = string(hashedPassword)

	if err := s.Repo.Create(&user); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) Update(id int64, req *UpdateUserRequest, uCtx share.UserContext) (*User, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Authorization: fetch target user
	target, e := s.Repo.ShowOne(id)
	if e != nil {
		return nil, e
	}
	if len(target.Users) == 0 {
		return nil, msg.NewErrorResponse("user_not_found", fmt.Errorf("user %d not found", id))
	}
	targetUser := target.Users[0]

	// Check: caller must have equal or higher role than target
	if uCtx.RoleID > targetUser.RoleID {
		return nil, msg.NewErrorResponse("access_denied", fmt.Errorf("cannot modify higher-privilege user"))
	}

	// Check: can't promote someone above your own level
	if req.RoleID != nil && uCtx.RoleID > *req.RoleID {
		return nil, msg.NewErrorResponse("access_denied", fmt.Errorf("cannot assign higher role than your own"))
	}

	// Check: can't modify yourself to deactivate
	if id == uCtx.UserID && req.StatusID != nil && *req.StatusID != 1 {
		return nil, msg.NewErrorResponse("access_denied", fmt.Errorf("cannot deactivate yourself"))
	}

	updates := map[string]any{}

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}
	if req.RoleName != nil {
		updates["role_name"] = *req.RoleName
	}
	if req.StatusID != nil {
		updates["status_id"] = *req.StatusID
	}

	if len(updates) == 0 {
		return nil, msg.NewErrorResponse("no_updates_provided", fmt.Errorf("empty"))
	}

	updates["updated_by"] = uCtx.UserID

	return s.Repo.Update(id, updates)
}

func (s *UserServiceImpl) Delete(id int64, uCtx share.UserContext) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}

	// Can't delete yourself
	if id == uCtx.UserID {
		return msg.NewErrorResponse("access_denied", fmt.Errorf("cannot delete yourself"))
	}

	// Authorization: fetch target user
	target, e := s.Repo.ShowOne(id)
	if e != nil {
		return e
	}
	if len(target.Users) == 0 {
		return msg.NewErrorResponse("user_not_found", fmt.Errorf("user %d not found", id))
	}
	targetUser := target.Users[0]

	// Check: caller must have equal or higher role than target
	if uCtx.RoleID > targetUser.RoleID {
		return msg.NewErrorResponse("access_denied", fmt.Errorf("cannot delete higher-privilege user"))
	}

	return s.Repo.Delete(id, uCtx.UserID)
}

func (s *UserServiceImpl) GetCreateForm() any {
	return map[string]any{
		"fields": map[string]any{
			"user_name":  map[string]string{"type": "text", "required": "true", "min": "4"},
			"password":   map[string]string{"type": "password", "required": "true", "min": "6"},
			"first_name": map[string]string{"type": "text"},
			"last_name":  map[string]string{"type": "text"},
			"email":      map[string]string{"type": "email"},
			"role_id":    map[string]string{"type": "number", "required": "true"},
			"role_name":  map[string]string{"type": "text", "required": "true"},
		},
	}
}

func (s *UserServiceImpl) GetUpdateForm(id int64) (*UserResponse, *error_responses.ErrorResponse) {
	return s.Repo.ShowOne(id)
}
