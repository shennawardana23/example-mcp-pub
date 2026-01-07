package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shennawardana23/example-mcp-pub/internal/app/model"
	"github.com/shennawardana23/example-mcp-pub/internal/app/service"
)

// ServiceController handles service-related requests
type ServiceController struct {
	service service.ServiceService
}

// NewServiceController creates a new service controller
func NewServiceController(service service.ServiceService) *ServiceController {
	return &ServiceController{service: service}
}

// Create handles service creation
// @Summary Create a new service
// @Description Create a new service in the catalog
// @Tags services
// @Accept json
// @Produce json
// @Param service body model.ServiceRequest true "Service details"
// @Success 201 {object} model.ServiceResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/services [post]
func (ctrl *ServiceController) Create(c *gin.Context) {
	var req model.ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	username, _ := c.Get("username")
	response, err := ctrl.service.Create(c.Request.Context(), &req, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "creation_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID handles service retrieval by ID
// @Summary Get service by ID
// @Description Get a service by its ID
// @Tags services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} model.ServiceResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/v1/services/{id} [get]
func (ctrl *ServiceController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid service ID",
		})
		return
	}

	response, err := ctrl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Error:   "not_found",
			Message: "Service not found",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// List handles service listing with filtering
// @Summary List services
// @Description List services with optional filtering and pagination
// @Tags services
// @Produce json
// @Param search query string false "Search term"
// @Param type query string false "Service type"
// @Param status query string false "Service status"
// @Param owner query string false "Service owner"
// @Param tag query string false "Tag filter"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} model.PaginatedResponse
// @Router /api/v1/services [get]
func (ctrl *ServiceController) List(c *gin.Context) {
	var req model.ServiceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid query parameters",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	response, err := ctrl.service.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update handles service updates
// @Summary Update service
// @Description Update an existing service
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param service body model.ServiceRequest true "Service details"
// @Success 200 {object} model.ServiceResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/services/{id} [put]
func (ctrl *ServiceController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid service ID",
		})
		return
	}

	var req model.ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	username, _ := c.Get("username")
	response, err := ctrl.service.Update(c.Request.Context(), id, &req, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete handles service deletion
// @Summary Delete service
// @Description Delete a service by ID
// @Tags services
// @Param id path int true "Service ID"
// @Success 204
// @Failure 404 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/services/{id} [delete]
func (ctrl *ServiceController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid service ID",
		})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "deletion_failed",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
