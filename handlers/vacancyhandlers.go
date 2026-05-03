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
	role := c.GetString("role")

	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can create vacancy"})
		return
	}
	var req struct {
		Title          string `json:"title"`
		DepartmentID   int    `json:"department_id"`
		CityID         int    `json:"city_id"`
		IsRemote       bool   `json:"is_remote"`
		SalaryMin      int    `json:"salary_min"`
		SalaryMax      int    `json:"salary_max"`
		Currency       string `json:"currency"`
		EmploymentType string `json:"employment_type"`
		Description    string `json:"description"`
		Status         string `json:"status"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := v.repo.CreateVacancy(
		context.Background(),
		req.Title,
		req.DepartmentID,
		req.CityID,
		req.IsRemote,
		req.SalaryMin,
		req.SalaryMax,
		req.Currency,
		req.EmploymentType,
		req.Description,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vacancy created"})
}

func (v *VacancyHandler) UpdateVacancy(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can update vacancy"})
		return
	}
	var req struct {
		Title          string `json:"title"`
		DepartmentID   int    `json:"department_id"`
		CityID         int    `json:"city_id"`
		IsRemote       bool   `json:"is_remote"`
		SalaryMin      int    `json:"salary_min"`
		SalaryMax      int    `json:"salary_max"`
		Currency       string `json:"currency"`
		EmploymentType string `json:"employment_type"`
		Description    string `json:"description"`
		Status         string `json:"status"`
	}

	if err = c.BindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = v.repo.UpdateVacancy(
		context.Background(),
		id,
		req.Title,
		req.DepartmentID,
		req.CityID,
		req.IsRemote,
		req.SalaryMin,
		req.SalaryMax,
		req.Currency,
		req.EmploymentType,
		req.Description,
		req.Status,
	)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vacancy updated"})
}

func (v *VacancyHandler) DeleteVacancy(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can update vacancy"})
		return
	}
	err = v.repo.DeleteVacancy(
		context.Background(),
		id,
	)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vacancy deleted"})
}
