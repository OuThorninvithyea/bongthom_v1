package user

import (

	// Commnuity Pacakges
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// Internal pacakges
	constants "admin-api/pkg/constants"
	response "admin-api/pkg/http"
	"admin-api/pkg/utls"
)

type UserHandler struct {
	Service UserService
}

func NewUserHandler(db *sqlx.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		Service: NewUserServiceImpl(db, rdb),
	}
}

// Show — GET /users
func (h *UserHandler) Show(c fiber.Ctx) error {
	var paging PagingRequest
	if err := c.Bind().Query(&paging); err != nil {
		paging.Page = 1
		paging.PerPage = 20
	}
	if paging.Page < 1 {
		paging.Page = 1
	}
	if paging.PerPage < 1 {
		paging.PerPage = 20
	}

	users, total, err := h.Service.List(paging.Page, paging.PerPage)
	if err != nil {
		errMsg := "Internal server error"
		if err.Err != nil {
			errMsg = err.Err.Error()
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(errMsg, constants.Generic_error, err.Err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponseWithPaing("Users retrieved", constants.Generic_success, users, paging.Page, paging.PerPage, total),
	)
}

// ShowOne — GET /users/:id
func (h *UserHandler) ShowOne(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid user ID", constants.Generic_invalid, err),
		)
	}

	user, uerr := h.Service.GetByID(id)
	if uerr != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewResponseError("User not found", constants.Generic_notFound, uerr.Err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse("User retrieved", constants.Generic_success, user),
	)
}

// CreateUser — POST /users/create
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	req := &CreateUserRequest{}
	v := utls.NewValidator()

	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid request body", constants.Generic_invalid, err),
		)
	}

	if err := v.Validate(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewValidatorError(ve),
			)
		}
		return err
	}

	// TODO: get createdBy from JWT when middleware is wired
	var createdBy int64 = 1

	user, uerr := h.Service.Create(req, createdBy)
	if uerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Failed to create user", constants.Generic_error, uerr.Err),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.NewResponse("User created", constants.Generic_success, user),
	)
}

// UpdateUser — PUT /users/update/:id
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid user ID", constants.Generic_invalid, err),
		)
	}

	req := &UpdateUserRequest{}
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid request body", constants.Generic_invalid, err),
		)
	}

	var updatedBy int64 = 1 // TODO: from JWT

	user, uerr := h.Service.Update(id, req, updatedBy)
	if uerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(uerr.MessageID, constants.Generic_error, uerr.Err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse("User updated", constants.Generic_success, user),
	)
}

// DeleteUser — DELETE /users/delete/:id
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid user ID", constants.Generic_invalid, err),
		)
	}

	var deletedBy int64 = 1 // TODO: from JWT

	if derr := h.Service.Delete(id, deletedBy); derr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(derr.MessageID, constants.Generic_error, derr.Err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse("User deleted", constants.Generic_success, nil),
	)
}

// GetUserFormCreate — GET /users/form/create
func (h *UserHandler) GetUserFormCreate(c fiber.Ctx) error {
	form := h.Service.GetCreateForm()
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse("Create user form", constants.Generic_success, form),
	)
}

// GetUserFormUpdate — GET /users/form/update/:id
func (h *UserHandler) GetUserFormUpdate(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid user ID", constants.Generic_invalid, err),
		)
	}

	user, uerr := h.Service.GetUpdateForm(id)
	if uerr != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewResponseError("User not found", constants.Generic_notFound, uerr.Err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse("Update user form", constants.Generic_success, user),
	)
}
