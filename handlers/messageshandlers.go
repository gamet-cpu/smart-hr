package handler

import (
	"context"
	"net/http"
	"smart-hr/repository"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	repo *repository.MessageRepository
}

func NewMessageHandler(m *repository.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: m}
}

func (m *MessageHandler) SentReq(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Topic       string `json:"topic"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := m.repo.SentReq(context.Background(),
		req.Name,
		req.Email,
		req.Topic,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "message sent"})
}
