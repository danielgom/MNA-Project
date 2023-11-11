package routes

import (
	"MNA-project/pkg/internal/pet"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"

	"MNA-project/pkg/context"
	"MNA-project/pkg/internal/pet/service"
	"MNA-project/pkg/util/errors"
)

// Handler is an instance of our pet handler API.
type Handler struct {
	petSvc service.PetService
}

// NewHandler returns a PetHandler instance.
func NewHandler(svc service.PetService) *Handler {
	return &Handler{petSvc: svc}
}

// Register adds all routes related to user service.
func (h *Handler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("v1/pets")
	authGroup.POST("", handler(h.Save))
	authGroup.GET("", handler(h.GetPetsByUser))
	authGroup.GET("/:id", handler(h.GetPetByID))
	authGroup.PUT("/:id", handler(h.UpdateById))
	authGroup.DELETE("/:id", handler(h.DeletePetByID))
}

// Save is used to create save a pet.
// @Summary Save Pet
// @Description Save a new user pet
// @Tags pet
// @Param user body pet.RegisterRequest true "Save request"
// @Success 201 {object} pet.GeneralResponse
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/pets [post]
func (h *Handler) Save(c context.Context) error {
	var req pet.RegisterRequest

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, signErr := h.petSvc.Save(c.Request().Context(), &req, c.GetUserID())
		if signErr != nil {
			return nil, signErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

// GetPetsByUser gets all pets by user.
// @Summary Get pets info by user
// @Description Get pets info by user
// @Tags pet
// @Success 200 {object} pet.Pets
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/pets [get]
func (h *Handler) GetPetsByUser(c context.Context) error {
	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		userID := c.GetUserID()

		res, getErr := h.petSvc.GetAllByUser(c.Request().Context(), userID)
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: pet.Pets{Pets: res},
		}, nil
	})
}

// GetPetByID gets pet information.
// @Summary Gets pet information by ID
// @Description Gets pet information
// @Tags pet
// @Param petID path int true "Pet ID"
// @Success 200 {object} pet.GeneralResponse
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/pets/{petID} [get]
func (h *Handler) GetPetByID(c context.Context) error {
	petID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		userID := c.GetUserID()
		res, getErr := h.petSvc.GetById(c.Request().Context(), int64(petID), userID)
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// UpdateById updates pet.
// @Summary Updates pet information
// @Description Updates multiple fields from pet in the DB
// @Tags pet
// @Param user body pet.UpdateRequest true "Save request"
// @Param petID path int true "pet ID"
// @Success 200 {object} pet.GeneralResponse
// @Failure 404 {object} errors.CommonError
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/pets/{petID} [put]
func (h *Handler) UpdateById(c context.Context) error {
	var req pet.UpdateRequest
	petID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		req.UserID = c.GetUserID()
		res, delErr := h.petSvc.UpdateById(c.Request().Context(), int64(petID), &req)
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// DeletePetByID deletes a pet by ID.
// @Summary Deletes a pet
// @Description Deletes pet by ID
// @Tags pet
// @Param userID path int true "User ID"
// @Success 200
// @Failure 403 {object} errors.CommonError
// @Failure 500 {object} errors.CommonError
// @Router /v1/pets/{userID} [delete]
func (h *Handler) DeletePetByID(c context.Context) error {
	petID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		userID := c.GetUserID()
		delErr := h.petSvc.DeleteById(c.Request().Context(), int64(petID), userID)
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status: http.StatusOK,
		}, nil
	})
}
