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

type TodoHandler struct {
	svc service.TodoService
}

func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// List godoc
// @Summary      List todos
// @Description  Return paginated list of todos for the authenticated user.
// @Tags         todos
// @Security     BearerAuth
// @Produce      json
// @Param        status      query     string  false  "pending | in_progress | done"
// @Param        priority    query     string  false  "low | medium | high"
// @Param        category_id query     string  false  "Filter by category ID"
// @Param        search      query     string  false  "Search in title (case-insensitive)"
// @Param        sort_by     query     string  false  "created_at | deadline"
// @Param        order       query     string  false  "asc | desc"
// @Param        page        query     int     false  "Page number (default 1)"
// @Param        limit       query     int     false  "Items per page (default 20, max 100)"
// @Success      200  {object}  todoListResp
// @Failure      401  {object}  errResp
// @Router       /todos [get]
func (h *TodoHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := repository.TodoFilter{
		Status:     c.Query("status"),
		Priority:   c.Query("priority"),
		CategoryID: c.Query("category_id"),
		Search:     c.Query("search"),
		SortBy:     c.DefaultQuery("sort_by", "created_at"),
		Order:      c.DefaultQuery("order", "desc"),
		Page:       page,
		Limit:      limit,
	}

	todos, total, err := h.svc.GetAll(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  todos,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Create godoc
// @Summary      Create a todo
// @Description  Create a new todo for the authenticated user.
// @Tags         todos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      model.CreateTodoRequest  true  "Todo payload"
// @Success      201   {object}  todoResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Router       /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	var req model.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": todo})
}

// GetByID godoc
// @Summary      Get a todo
// @Description  Return a single todo by ID.
// @Tags         todos
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Todo ID"
// @Success      200  {object}  todoResp
// @Failure      400  {object}  errResp
// @Failure      401  {object}  errResp
// @Failure      404  {object}  errResp
// @Router       /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	todo, err := h.svc.GetByID(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

// Update godoc
// @Summary      Update a todo
// @Description  Update a todo's fields (title, description, priority, category, deadline).
// @Tags         todos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "Todo ID"
// @Param        body  body      model.UpdateTodoRequest  true  "Fields to update"
// @Success      200   {object}  todoResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Failure      404   {object}  errResp
// @Router       /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	var req model.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

// UpdateStatus godoc
// @Summary      Update todo status
// @Description  Update only the status of a todo. Sets completed_at when status becomes done.
// @Tags         todos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                         true  "Todo ID"
// @Param        body  body      model.UpdateTodoStatusRequest  true  "New status"
// @Success      200   {object}  todoResp
// @Failure      400   {object}  errResp
// @Failure      401   {object}  errResp
// @Failure      404   {object}  errResp
// @Router       /todos/{id}/status [patch]
func (h *TodoHandler) UpdateStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	var req model.UpdateTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.svc.UpdateStatus(c.Request.Context(), userID, c.Param("id"), req.Status)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

// Delete godoc
// @Summary      Delete a todo
// @Description  Permanently delete a todo.
// @Tags         todos
// @Security     BearerAuth
// @Param        id   path  string  true  "Todo ID"
// @Success      204
// @Failure      401  {object}  errResp
// @Failure      404  {object}  errResp
// @Router       /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	if err := h.svc.Delete(c.Request.Context(), userID, c.Param("id")); err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
