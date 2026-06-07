package user

import (

	// Community pacakges
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// Interntal pacakges

	config "admin-api/configs"
	"admin-api/pkg/redis_util"
	error_responses "admin-api/pkg/responses"
	custom_sql "admin-api/pkg/sql"
)

type UserRepo interface {
	Show(u UserShowRequest) (*UserResponse, *error_responses.ErrorResponse)
	ShowOne(id int64) (*UserResponse, *error_responses.ErrorResponse)
	GetByUserName(userName string) (*User, *error_responses.ErrorResponse)
	Create(user *User) *error_responses.ErrorResponse
	Update(id int64, updates map[string]any) (*User, *error_responses.ErrorResponse)
	Delete(id int64, deletedBy int64) *error_responses.ErrorResponse
}

type UserRepoImpl struct {
	db       *sqlx.DB
	redis    *redis.Client
	cacheTTL time.Duration
}

func NewUserRepoImpl(db *sqlx.DB, rdb *redis.Client) UserRepo {
	cfg := config.InitRedis()
	ttl := time.Duration(cfg.RedisExpire) * time.Second
	return &UserRepoImpl{
		db:       db,
		redis:    rdb,
		cacheTTL: ttl,
	}
}

func (r *UserRepoImpl) Show(userRequest UserShowRequest) (*UserResponse, *error_responses.ErrorResponse) {
	// Calculatings for skipping user in
	var per_page = userRequest.PageOption.Perpage
	var page = userRequest.PageOption.Page
	var offset = (page - 1) * per_page
	var limit_clause = fmt.Sprintf(" LIMIT %d OFFSET %d", per_page, offset)
	var sql_orderby = custom_sql.BuildSQLSort(userRequest.Sorts)

	sql_filters, args_filters := custom_sql.BuildSQLFilter(userRequest.Filters)
	if len(args_filters) > 0 {
		sql_filters = " AND " + sql_filters
	}

	if searchClause, searchArgs := custom_sql.BuildSQLSearch(
		[]string{"u.user_name", "u.first_name", "u.last_name", "u.user_alias", "u.email"},
		userRequest.Search, len(args_filters)+1,
	); searchClause != "" {
		sql_filters += " AND " + searchClause
		args_filters = append(args_filters, searchArgs...)
	}

	msg := error_responses.ErrorResponse{}

	// Total count with same filters (no limit/offset/order)
	var total int
	countQuery := fmt.Sprintf(
		`SELECT COUNT(*) FROM (
		SELECT user_name, first_name, last_name, email, role_name, role_id, is_admin,
		login_session, last_login, currency_id, language_id, status_id, created_at, updated_at
		FROM tbl_users u
		WHERE deleted_at IS NULL %s
	) AS t`,
		sql_filters)
	err := r.db.Get(&total, countQuery, args_filters...)
	if err != nil {
		return nil, msg.NewErrorResponse("database_error", err)
	}
	var users []User
	query := fmt.Sprintf(
		`SELECT id, user_name, first_name, last_name, email, role_name, role_id, is_admin,
		 login_session, last_login, currency_id, language_id, status_id, created_at, updated_at
		 FROM tbl_users u
		 WHERE deleted_at IS NULL
		%s %s %s`, sql_filters, sql_orderby, limit_clause)

	err = r.db.Select(&users, query, args_filters...)
	if err != nil {
		return nil, msg.NewErrorResponse("database_error", err)
	}
	return &UserResponse{Users: users, Total: total}, nil
}

func (r *UserRepoImpl) ShowOne(id int64) (*UserResponse, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	// Read-through cache
	cacheKey := fmt.Sprintf("user_info_id:%d", id)
	if r.redis != nil {
		rdb := redis_util.NewRedisUtil(r.redis, r.cacheTTL)
		var cached User
		if err := rdb.GetCacheKey(cacheKey, &cached, context.Background()); err == nil {
			return &UserResponse{Users: []User{cached}, Total: 1}, nil
		}
	}

	var user User
	err := r.db.Get(&user,
		`SELECT * FROM tbl_users WHERE id = $1 AND deleted_at IS NULL LIMIT 1`, id,
	)
	if err != nil {
		return nil, msg.NewErrorResponse("user_not_found", err)
	}

	// Populate cache
	if r.redis != nil {
		rdb := redis_util.NewRedisUtil(r.redis, r.cacheTTL)
		_ = rdb.SetCacheKey(cacheKey, &user, context.Background())
	}

	return &UserResponse{
		Users: []User{user}, Total: 1,
	}, nil
}

func (r *UserRepoImpl) GetByUserName(userName string) (*User, *error_responses.ErrorResponse) {
	var user User
	err := r.db.Get(&user,
		`SELECT * FROM tbl_users WHERE user_name = $1 LIMIT 1`, userName,
	)
	if err != nil {
		return nil, nil // not found is OK — caller checks
	}
	return &user, nil
}

func (r *UserRepoImpl) Create(user *User) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}

	// Insert into PostgreSQL and capture the generated ID
	query := `
		INSERT INTO tbl_users (
			first_name, last_name, user_name, email, password,
			role_name, role_id, login_session, status_id, "order",
			created_by, created_at
		) VALUES (
			:first_name, :last_name, :user_name, :email, :password,
			:role_name, :role_id, :login_session, :status_id, :order,
			:created_by, :created_at
		) RETURNING id`

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return msg.NewErrorResponse("database_error", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&user.ID); err != nil {
			return msg.NewErrorResponse("database_error", err)
		}
	}

	// Set redis data
	key := fmt.Sprintf("user_info_id:%d", user.ID)
	rdb := redis_util.NewRedisUtil(r.redis, r.cacheTTL)
	rdb.SetCacheKey(key, user, context.Background())

	return nil
}

func (r *UserRepoImpl) Update(id int64, updates map[string]any) (*User, *error_responses.ErrorResponse) {
	msg := error_responses.ErrorResponse{}

	setClauses := []string{}
	args := []any{}
	i := 1
	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	if len(setClauses) == 0 {
		return nil, msg.NewErrorResponse("no_updates_provided", fmt.Errorf("empty update"))
	}

	query := fmt.Sprintf(
		`UPDATE tbl_users SET %s WHERE id = $%d AND deleted_at IS NULL RETURNING *`,
		strings.Join(setClauses, ", "), i,
	)
	args = append(args, id)

	var user User
	err := r.db.Get(&user, query, args...)
	if err != nil {
		return nil, msg.NewErrorResponse("database_error", err)
	}
	return &user, nil
}

func (r *UserRepoImpl) Delete(id int64, deletedBy int64) *error_responses.ErrorResponse {
	msg := error_responses.ErrorResponse{}

	result, err := r.db.Exec(
		`UPDATE tbl_users SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`,
		deletedBy, id,
	)
	if err != nil {
		return msg.NewErrorResponse("database_error", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return msg.NewErrorResponse("user_not_found", fmt.Errorf("user %d not found", id))
	}
	return nil
}
