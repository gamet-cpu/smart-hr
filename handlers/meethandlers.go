package handler

import (
	"context"
	"net/http"
	"smart-hr/repository"

	"github.com/gin-gonic/gin"
)

type MeetHandler struct {
	repo *repository.MeetRepository
}

func NewMeetHandler(m *repository.MeetRepository) *MeetHandler {
	return &MeetHandler{repo: m}
}

// GetAllMeets godoc
// @Summary      Список встреч
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array}  map[string]interface{}
// @Router       /meets [get]
func (m *MeetHandler) GetAllMeets(c *gin.Context) {
	userId := c.GetInt("user_id")
	role := c.GetString("role")
	meets, err := m.repo.GetAllMeets(context.Background(), userId, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meets)
}
