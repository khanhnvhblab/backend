package handler

import (
	"net/http"
	"todolist/backend/internal/model"
	"todolist/backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetMe godoc
// @Summary      Get current user profile
// @Description  Return the profile of the authenticated user.
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  userResp
// @Failure      401  {object}  errResp
// @Failure      404  {object}  errResp
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	user, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateMe godoc
// @Summary      Update current user profile
// @Description  Update the name of the authenticated user.
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.UpdateUserRequest  true  "Fields to update"
// @Success      200   {object}  userResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Failure      404   {object}  errResp
// @Router       /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Update(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteMe godoc
// @Summary      Delete current user account
// @Description  Permanently delete the authenticated user's account.
// @Tags         users
// @Security     BearerAuth
// @Success      204
// @Failure      401  {object}  errResp
// @Failure      404  {object}  errResp
// @Router       /users/me [delete]
func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	if err := h.svc.Delete(c.Request.Context(), userID); err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
