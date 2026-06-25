package handler

import (
	"net/http"

	"enterprise-order-management-api/internal/dto"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	users service.UserService
}

func NewUserHandler(users service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(c echo.Context) error {
	query := dto.UserListQuery{
		Page:   parseIntQuery(c, "page", 1),
		Limit:  parseIntQuery(c, "limit", 10),
		Search: c.QueryParam("search"),
	}

	res, meta, err := h.users.List(c.Request().Context(), query)
	if err != nil {
		return err
	}
	return response.Paginated(c, res, meta)
}

func (h *UserHandler) FindByID(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	res, err := h.users.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) Update(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.users.Update(c.Request().Context(), id, req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) UpdateMe(c echo.Context) error {
	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.users.UpdateProfile(c.Request().Context(), appmiddleware.CurrentUserID(c), req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	if err := h.users.ChangePassword(c.Request().Context(), appmiddleware.CurrentUserID(c), req); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "Password changed successfully")
}

func (h *UserHandler) UploadAvatar(c echo.Context) error {
	file, err := c.FormFile("avatar")
	if err != nil {
		return apperror.BadRequest("Avatar file is required")
	}
	res, err := h.users.UploadAvatar(c.Request().Context(), appmiddleware.CurrentUserID(c), file)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) UploadProfileVideo(c echo.Context) error {
	file, err := c.FormFile("video")
	if err != nil {
		return apperror.BadRequest("Video file is required")
	}
	res, err := h.users.UploadProfileVideo(c.Request().Context(), appmiddleware.CurrentUserID(c), file)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.users.Delete(c.Request().Context(), id, appmiddleware.CurrentUserID(c)); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "User deleted successfully")
}
