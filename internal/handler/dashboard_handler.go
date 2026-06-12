package handler

import (
	"net/http"
	"todolist/backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DashboardHandler struct {
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetStats godoc
// @Summary      Get dashboard stats
// @Description  Return todo counts by status, overdue and due-soon for the authenticated user.
// @Tags         dashboard
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dashboardStatsResp
// @Failure      401  {object}  errResp
// @Router       /dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := c.MustGet("user_id").(bson.ObjectID)

	stats, err := h.svc.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(service.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
