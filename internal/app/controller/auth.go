package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shennawardana23/example-mcp-pub/internal/app/model"
	"github.com/shennawardana23/example-mcp-pub/internal/app/service"
)

// AuthController handles authentication-related requests
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	response, err := ctrl.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error:   "authentication_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Me returns the current user's information
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/me [get]
func (ctrl *AuthController) Me(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	user, err := ctrl.authService.GetUserByUsername(c.Request.Context(), username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "user_retrieval_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
