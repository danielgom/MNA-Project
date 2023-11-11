// Package routes will be the responsible for adding all routes from the service.
package routes

import (
	"MNA-project/pkg/internal/user/service"
	"fmt"
	"net/http"
	"strconv"

	"MNA-project/pkg/context"
	"MNA-project/pkg/internal/user"
	"MNA-project/pkg/util/errors"

	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

var (
	errInvalidID = fmt.Errorf("invalid id")
)

// Handler is an instance of our user handler API.
type Handler struct {
	UsrSvc service.UserService
}

// NewHandler returns a UserHandler instance.
func NewHandler(svc service.UserService) *Handler {
	return &Handler{UsrSvc: svc}
}

// Register adds all routes related to user service.
func (h *Handler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("v1/users")
	authGroup.POST("/signup", handler(h.SignUp))
	authGroup.GET("/:id", handler(h.GetUserById))
	authGroup.GET("", handler(h.GetAllUsers))
	authGroup.PUT("/:id", handler(h.UpdateById))
	authGroup.DELETE("/:id", handler(h.DeleteUserById))
	authGroup.POST("/login", handler(h.Login))
}

// SignUp is used to create a new user.
// @Summary SignUp User
// @Description Register a User with email and password
// @Tags user
// @Param user body user.RegisterRequest true "Register request"
// @Success 201 {object} user.RegisterResponse
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/users/signup [post]
func (h *Handler) SignUp(c context.Context) error {
	var req user.RegisterRequest

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, signErr := h.UsrSvc.SignUp(c.Request().Context(), &req)
		if signErr != nil {
			return nil, signErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

// GetUserById gets user by id.
// @Summary Gets user info by ID
// @Description Gets user info with ID
// @Tags user
// @Param userID path int true "User ID"
// @Success 200 {object} user.Response
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/users/{userID} [get]
func (h *Handler) GetUserById(c context.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("invalid parameter", errInvalidID)
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.UsrSvc.GetById(c.Request().Context(), int64(userID))
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// GetAllUsers gets all users saved in the db.
// @Summary Gets all users
// @Description Gets users saved
// @Tags user
// @Success 200 {object} user.Users
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/users [get]
func (h *Handler) GetAllUsers(c context.Context) error {
	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.UsrSvc.GetAll(c.Request().Context())
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: user.Users{Users: res},
		}, nil
	})
}

// UpdateById updates user.
// @Summary Updates user information
// @Description Updates multiple fields from user in the DB
// @Tags user
// @Param user body user.UpdateRequest true "Save request"
// @Param userID path int true "User ID"
// @Success 200 {object} user.Response
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/users/{userID} [put]
func (h *Handler) UpdateById(c context.Context) error {
	var req user.UpdateRequest

	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("invalid parameter", errInvalidID)
	}

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, delErr := h.UsrSvc.UpdateById(c.Request().Context(), int64(userID), &req)
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// DeleteUserById deletes a user by ID.
// @Summary Deletes user
// @Description Deletes user by ID
// @Tags user
// @Param userID path int true "User ID"
// @Success 200
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/users/{userID} [delete]
func (h *Handler) DeleteUserById(c context.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("invalid parameter", errInvalidID)
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		delErr := h.UsrSvc.DeleteById(c.Request().Context(), int64(userID))
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status: http.StatusOK,
		}, nil
	})
}

// Login returns a JWT based on the user that has been logged in.
// @Summary User login
// @Description Logins a user and returns a JWT
// @Tags user
// @Param user body user.LoginRequest true "Login request"
// @Success 201 {object} user.LoginResponse
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/users/login [post]
func (h *Handler) Login(c context.Context) error {
	var req user.LoginRequest

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, logErr := h.UsrSvc.Login(c.Request().Context(), &req)
		if logErr != nil {
			return nil, logErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}
