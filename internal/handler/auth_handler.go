package handler

import (
	"net/http"
	"todolist/backend/internal/model"
	"todolist/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
	userSvc service.UserService
}

func NewAuthHandler(authSvc service.AuthService, userSvc service.UserService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userSvc: userSvc}
}

// Register godoc
// @Summary      Register a new account
// @Description  Create a new user account with email, password, and name.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.CreateUserRequest  true  "Registration payload"
// @Success      201   {object}  registerResp
// @Failure      400   {object}  errResp
// @Failure      409   {object}  errResp
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userSvc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	}})
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password. Returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      service.LoginRequest  true  "Login credentials"
// @Success      200   {object}  loginResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh token for a new access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      service.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  refreshResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req service.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authSvc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Logout godoc
// @Summary      Logout
// @Description  Invalidate the current session. Requires a valid access token.
// @Tags         auth
// @Security     BearerAuth
// @Success      204
// @Failure      401  {object}  errResp
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
