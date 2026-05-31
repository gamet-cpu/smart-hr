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

// GetAllVacancy godoc
// @Summary      Список вакансий
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array}  map[string]interface{}
// @Router       /vacancy [get]
func (v *VacancyHandler) GetAllVacancy(c *gin.Context) {
	users, err := v.repo.GetAllVacancy(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetVacancyByID godoc
// @Summary      Вакансия по ID
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /vacancy/{id} [get]
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

// CreateVacancy godoc
// @Summary      Создать вакансию (только company)
// @Tags         vacancy
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object{title=string,department_id=int,city_id=int,is_remote=bool,salary_min=int,salary_max=int,currency=string,employment_type=string,description=string,status=string} true "Данные вакансии"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /vacancy [post]
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

// UpdateVacancy godoc
// @Summary      Обновить вакансию (только company)
// @Tags         vacancy
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "ID вакансии"
// @Param        body body object{title=string,department_id=int,city_id=int,is_remote=bool,salary_min=int,salary_max=int,currency=string,employment_type=string,description=string,status=string} true "Поля для обновления"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /vacancy/{id} [put]
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

// DeleteVacancy godoc
// @Summary      Удалить вакансию (только company)
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /vacancy/{id} [delete]
func (v *VacancyHandler) DeleteVacancy(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can delete vacancy"})
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

// SoftDeleteVacancy godoc
// @Summary      Удаление васансии перенос в архив
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /vacancy/archive/{id} [delete]
func (v *VacancyHandler) SoftDelete(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company can delete vacancy"})
		return
	}

	vacancy, err := v.repo.GetVacancyByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
		return
	}

	err = v.repo.SoftDelete(
		context.Background(),
		id,
		*vacancy,
	)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vacancy soft deleted"})
}

// GetArchiveVacancy godoc
// @Summary      Список архивированных вакансий
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array}  map[string]interface{}
// @Router       /vacancy/archive [get]
func (v *VacancyHandler) GetArchiveVacancy(c *gin.Context) {
	users, err := v.repo.GetArchiveVacancy(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetArchiveVacancyByID godoc
// @Summary      Архивная Вакансия по ID
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /vacancy/archive/{id} [get]
func (v *VacancyHandler) GetArchiveVacancyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	vacancy, err := v.repo.GetArchiveVacancyByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

// SoftDeleteVacancy godoc
// @Summary      Деархивирование вакансии
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /vacancy/archive/{id} [put]
func (v *VacancyHandler) UnArchiveVacancy(c *gin.Context) {
	role := c.GetString("role")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if role != "company" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company unarhive vacancy"})
		return
	}

	vacancy, err := v.repo.GetArchiveVacancyByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
		return
	}

	err = v.repo.UnArchiveVacancy(
		context.Background(),
		id,
		*vacancy,
	)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vacancy unarchived"})
}

// TopVacancy godoc
// @Summary      Топ вакансий
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /vacancy/topvacancy [get]
func (v *VacancyHandler) TopVacancy(c *gin.Context) {
	vacancy, err := v.repo.TopVacancy(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

// ActualVacancy godoc
// @Summary      Актуальные вакансии
// @Tags         vacancy
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /vacancy/actualvacancy [get]
func (v *VacancyHandler) ActualVacancy(c *gin.Context) {
	vacancy, err := v.repo.ActualVacancy(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}
