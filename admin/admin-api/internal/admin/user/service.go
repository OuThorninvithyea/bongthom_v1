package user

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	error_responses "admin-api/pkg/responses"
)

type UserService interface {
	List(page, perPage int) ([]User, int, *error_responses.ErrorResponse)
	GetByID(id int64) (*User, *error_responses.ErrorResponse)
	Create(req *CreateUserRequest, createdBy int64) (*User, *error_responses.ErrorResponse)
	Update(id int64, req *UpdateUserRequest, updatedBy int64) (*User, *error_responses.ErrorResponse)
	Delete(id int64, deletedBy int64) *error_responses.ErrorResponse
	GetCreateForm() any
	GetUpdateForm(id int64) (*User, *error_responses.ErrorResponse)
}

type UserServiceImpl struct {
	Repo  UserRepo
	Redis *redis.Client
}

func NewUserServiceImpl(db *sqlx.DB, rdb *redis.Client) *UserServiceImpl {
	return &UserServiceImpl{
		Repo:  NewUserRepoImpl(db),
		Redis: rdb,
	}
}

func (s *UserServiceImpl) List(page, perPage int) ([]User, int, *error_responses.ErrorResponse) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return s.Repo.FindAll(page, perPage)
}

func (s *UserServiceImpl) GetByID(id int64) (*User, *error_responses.ErrorResponse) {
	return s.Repo.FindByID(id)
}

func (s *UserServiceImpl) Create(req *CreateUserRequest, createdBy int64) (*User, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Check duplicate username
	existing, _ := s.Repo.FindByUserName(req.UserName)
	if existing != nil {
		return nil, msg.NewErrorResponse("user_name_taken", fmt.Errorf("username exists"))
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, msg.NewErrorResponse("password_hash_failed", err)
	}

	user := &User{
		UserName:  req.UserName,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		RoleID:    req.RoleID,
		RoleName:  req.RoleName,
		CreatedBy: &createdBy,
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) Update(id int64, req *UpdateUserRequest, updatedBy int64) (*User, *error_responses.ErrorResponse) {
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
		msg := error_responses.ErrorResponse{}
		return nil, msg.NewErrorResponse("no_updates_provided", fmt.Errorf("empty"))
	}

	updates["updated_by"] = updatedBy

	return s.Repo.Update(id, updates)
}

func (s *UserServiceImpl) Delete(id int64, deletedBy int64) *error_responses.ErrorResponse {
	return s.Repo.SoftDelete(id, deletedBy)
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

func (s *UserServiceImpl) GetUpdateForm(id int64) (*User, *error_responses.ErrorResponse) {
	return s.Repo.FindByID(id)
}
