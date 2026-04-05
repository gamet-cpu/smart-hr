package handler

import (
	"context"
	"net/http"
	"strconv"

	"smart-hr/repository"

	"github.com/gin-gonic/gin"
)

type VacancyHandler struct {
	repo *repository.VacancyRepository
}

func NewVacancyHandler(r *repository.VacancyRepository) *VacancyHandler {
	return &VacancyHandler{repo: r}
}

func (v *VacancyHandler) GetAllVacancy(c *gin.Context) {
	users, err := v.repo.GetAllVacancy(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (v *VacancyHandler) GetVacancyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	vacancy, err := v.repo.GetVacancyByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

func (v *VacancyHandler) CreateVacancy(c *gin.Context) {
	err := v.repo.CreateVacancy(context.Background())
	if err == "not company"{
		c.JSON(http.StatusBadRequest, gin.H("error": "You can't create vacancy"))
	}
	c.JSON(http.StatusOK, nil)
}
