package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"smart-hr/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(r *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{name=string,email=string,password=string,role=string} true "Данные пользователя"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	err = h.repo.Create(context.Background(),
		req.Name,
		req.Email,
		string(hash),
		req.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

// Login godoc
// @Summary      Вход / получение токена
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{email=string,password=string} true "Логин и пароль"
// @Success      200 {object} map[string]string "token"
// @Failure      401 {object} map[string]string
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, hash, role, err := h.repo.GetByEmail(context.Background(), req.Email)
	if err != nil ||
		bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	t, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}

// GetMe godoc
// @Summary      Текущий пользователь
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]string
// @Router       /user/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	id := c.GetInt("user_id")

	user, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAll godoc
// @Summary      Все пользователи
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array}  map[string]interface{}
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Update godoc
// @Summary      Обновить профиль
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object{name=string,company_name=string,phone=string,description=string} true "Поля для обновления"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /user [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.GetInt("user_id")

	var req struct {
		Name        string `json:"name"`
		CompanyName string `json:"company_name"`
		Phone       string `json:"phone"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.Update(
		context.Background(),
		id,
		req.Name,
		req.CompanyName,
		req.Phone,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete godoc
// @Summary      Удалить аккаунт
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]string
// @Router       /user [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.GetInt("user_id")

	err := h.repo.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Respond Vacancy godoc
// @Summary      Откликтнутся на вакансию
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]string
// @Router       /user/{id} [put]
func (h *UserHandler) RespondVacancy(c *gin.Context) {
	userId := c.GetInt("user_id")
	idStr := c.Param("id")
	vacId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
	}
	err = h.repo.RespondVacancy(context.Background(), userId, vacId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "responded"})
}

// Responded Vacancy godoc
// @Summary      Список Соискателей
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID вакансии"
// @Success      200 {object} map[string]string
// @Router       /users/responses/{id} [get]
func (h *UserHandler) GetResponses(c *gin.Context) {
	idStr := c.Param("id")
	vacId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
	}
	users, err := h.repo.GetResponses(context.Background(), vacId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// MyRespondedVacancy godoc
// @Summary      Мои отклики
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]string
// @Router       /users/myresponses [get]
func (h *UserHandler) MyResponses(c *gin.Context) {
	userId := c.GetInt("user_id")
	vacancy, err := h.repo.MyResponses(context.Background(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vacancy)
}
