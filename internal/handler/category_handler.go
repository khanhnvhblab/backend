package handler

import (
	"net/http"
	"strconv"
	"todolist/backend/internal/model"
	"todolist/backend/internal/repository"
	"todolist/backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List godoc
// @Summary      List categories
// @Description  Return all categories for the authenticated user.
// @Tags         categories
// @Security     BearerAuth
// @Produce      json
// @Param        page   query     int  false  "Page number (default 1)"
// @Param        limit  query     int  false  "Items per page (default 20, max 100)"
// @Success      200  {object}  categoryListResp
// @Failure      401  {object}  errResp
// @Router       /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	cats, total, err := h.svc.GetAll(c.Request.Context(), userID, repository.CategoryFilter{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  cats,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Create godoc
// @Summary      Create a category
// @Description  Create a new category. Color defaults to #6366f1 if not provided.
// @Tags         categories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.CreateCategoryRequest  true  "Category payload"
// @Success      201   {object}  categoryResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Router       /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": cat})
}

// GetByID godoc
// @Summary      Get a category
// @Description  Return a single category by ID.
// @Tags         categories
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  categoryResp
// @Failure      400  {object}  errResp
// @Failure      401  {object}  errResp
// @Failure      404  {object}  errResp
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	cat, err := h.svc.GetByID(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cat})
}

// Update godoc
// @Summary      Update a category
// @Description  Update a category's name or color.
// @Tags         categories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "Category ID"
// @Param        body  body      model.UpdateCategoryRequest  true  "Fields to update"
// @Success      200   {object}  categoryResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Failure      404   {object}  errResp
// @Router       /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cat})
}

// Delete godoc
// @Summary      Delete a category
// @Description  Permanently delete a category.
// @Tags         categories
// @Security     BearerAuth
// @Param        id   path  string  true  "Category ID"
// @Success      204
// @Failure      401  {object}  errResp
// @Failure      404  {object}  errResp
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
