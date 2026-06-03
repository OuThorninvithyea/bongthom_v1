package user

import (

	// Commnuity pacakges
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	// Internal pacakges
	constants "admin-api/pkg/constants"
	response "admin-api/pkg/http"
	"admin-api/pkg/translate"
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

	users, total, e := h.Service.List(paging.Page, paging.PerPage)
	if e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(msg, constants.Generic_error, e.Err),
		)
	}

	msg, e_msg := translate.TranslateWithError(c, "users_retrieved")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponseWithPaing(msg, constants.Generic_success, users, paging.Page, paging.PerPage, total),
	)
}

// ShowOne — GET /users/:id
func (h *UserHandler) ShowOne(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		msg, _ := translate.TranslateWithError(c, "invalid_user_id")
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(msg, constants.Generic_invalid, err),
		)
	}

	user, e := h.Service.GetByID(id)
	if e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusNotFound).JSON(
				response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
			)
		}
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewResponseError(msg, constants.Generic_notFound, e.Err),
		)
	}

	msg, e_msg := translate.TranslateWithError(c, "user_retrieved")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(msg, constants.Generic_success, user),
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
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fe := ve[0]
			msg, _ := translate.TranslateWithError(c, "validation_"+fe.Tag(),
				map[string]any{
					"Field": fe.Field(),
					"Param": fe.Param(),
				})
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(msg, constants.Generic_invalid, err),
			)
		}
		return err
	}

	var createdBy int64 = 1 // TODO: from JWT

	user, e := h.Service.Create(req, createdBy)
	if e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
			)
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(msg, constants.Generic_error, e.Err),
		)
	}

	msg, e_msg := translate.TranslateWithError(c, "user_created")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusCreated).JSON(
		response.NewResponse(msg, constants.Generic_success, user),
	)
}

// UpdateUser — PUT /users/update/:id
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		msg, _ := translate.TranslateWithError(c, "invalid_user_id")
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(msg, constants.Generic_invalid, err),
		)
	}

	req := &UpdateUserRequest{}
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError("Invalid request body", constants.Generic_invalid, err),
		)
	}

	var updatedBy int64 = 1 // TODO: from JWT

	user, e := h.Service.Update(id, req, updatedBy)
	if e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
			)
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(msg, constants.Generic_error, e.Err),
		)
	}

	msg, e_msg := translate.TranslateWithError(c, "user_updated")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(msg, constants.Generic_success, user),
	)
}

// DeleteUser — DELETE /users/delete/:id
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		msg, _ := translate.TranslateWithError(c, "invalid_user_id")
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(msg, constants.Generic_invalid, err),
		)
	}

	var deletedBy int64 = 1 // TODO: from JWT

	if e := h.Service.Delete(id, deletedBy); e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(msg, constants.Generic_error, e.Err),
		)
	}

	msg, e_msg := translate.TranslateWithError(c, "user_deleted")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(msg, constants.Generic_success, nil),
	)
}

// GetUserFormCreate — GET /users/form/create
func (h *UserHandler) GetUserFormCreate(c fiber.Ctx) error {
	form := h.Service.GetCreateForm()
	msg, e_msg := translate.TranslateWithError(c, "form_create_retrieved")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(msg, constants.Generic_success, form),
	)
}

// GetUserFormUpdate — GET /users/form/update/:id
func (h *UserHandler) GetUserFormUpdate(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		msg, _ := translate.TranslateWithError(c, "invalid_user_id")
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(msg, constants.Generic_invalid, err),
		)
	}

	user, e := h.Service.GetUpdateForm(id)
	if e != nil {
		msg, e_msg := translate.TranslateWithError(c, e.MessageID)
		if e_msg != nil {
			return c.Status(fiber.StatusNotFound).JSON(
				response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
			)
		}
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewResponseError(msg, constants.Generic_notFound, e.Err),
		)
	}

	msg, e_msg := translate.TranslateWithError(c, "form_update_retrieved")
	if e_msg != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewResponseError(e_msg.Err.Error(), constants.Translate_Failed, e_msg.Err),
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(msg, constants.Generic_success, user),
	)
}
