package handler

import (
	"context"
	"net/http"
	"smart-hr/model"
	"smart-hr/repository"
	"strconv"

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
// @Tags         meets
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

// GetMeetsById godoc
// @Summary      Встреча по Id
// @Tags         meets
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID встречи"
// @Success      200 {array}  map[string]interface{}
// @Router       /meets/{id} [get]
func (m *MeetHandler) GetMeetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	meets, err := m.repo.GetMeetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meets)
}

// Create meet godoc
// @Summary      Назначить встречу (только company)
// @Tags         meets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object{recruter_id=int, user_id=int, meet_date=string, start_time=string, meet_leight=int, meet_type=string, link=string, additional=string} true "Данные встречи"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /meets [post]
func (m *MeetHandler) CreateMeet(c *gin.Context) {
	userId := c.GetInt("user_id")
	role := c.GetString("role")

	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can create meets"})
		return
	}
	var meet model.Meet

	if err := c.BindJSON(&meet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := m.repo.CreateMeet(
		context.Background(),
		userId,
		meet.UserID,
		meet.MeetDate,
		meet.StartTime,
		meet.MeetLeight,
		meet.MeetType,
		meet.Link,
		meet.Additional,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meet created"})
}

// Update meet godoc
// @Summary      Изменить встречу (только company)
// @Tags         meets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "ID встречи"
// @Param        body body object{meet_date=string, start_time=string, meet_leight=int, meet_type=string, link=string, additional=string} true "Данные встречи"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /meets/{id} [put]
func (m *MeetHandler) UpdateMeet(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can update meets"})
		return
	}
	var meet struct {
		MeetDate   string `json:"meet_date"`
		StartTime  string `json:"start_time"`
		MeetLeight int    `json:"meet_leight"`
		MeetType   string `json:"meet_type"`
		Link       string `json:"link"`
		Additional string `json:"additional"`
	}

	if err = c.BindJSON(&meet); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = m.repo.UpdateMeet(
		context.Background(),
		id,
		meet.MeetDate,
		meet.StartTime,
		meet.MeetLeight,
		meet.MeetType,
		meet.Link,
		meet.Additional,
	)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vacancy updated"})
}

// DeleteMeet godoc
// @Summary      Удалить встречу (только company)
// @Tags         meets
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID встречи"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /meets/{id} [delete]
func (m *MeetHandler) DeleteMeet(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can delete meet"})
		return
	}
	err = m.repo.DeleteMeet(
		context.Background(),
		id,
	)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meet deleted"})
}
