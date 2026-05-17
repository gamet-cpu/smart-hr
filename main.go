package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "smart-hr/docs" // ← этого не было!
	handler "smart-hr/handlers"
	"smart-hr/middleware"
	"smart-hr/repository"
)

// @title           Smart HR API
// @version         1.0
// @description     REST API для платформы Smart HR
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	dbpool, err := pgxpool.New(context.Background(),
		viper.GetString("db.url"))
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(dbpool)
	userHandler := handler.NewUserHandler(userRepo)
	vacancyRepo := repository.NewVacancyRepository(dbpool)
	vacancyHandler := handler.NewVacancyHandler(vacancyRepo)
	messageRepo := repository.NewMessageRepository(dbpool)
	messageHandler := handler.NewMessageHandler(messageRepo)
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/message", messageHandler.SentReq)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user/me", userHandler.GetMe)
		auth.GET("/users", userHandler.GetAll)
		auth.PUT("/user", userHandler.Update)
		auth.DELETE("/user", userHandler.Delete)
		auth.PUT("/user/:id", userHandler.RespondVacancy)
		auth.GET("/users/responses/:id", userHandler.GetResponses)
		auth.GET("/vacancy", vacancyHandler.GetAllVacancy)
		auth.GET("/vacancy/:id", vacancyHandler.GetVacancyByID)
		auth.POST("/vacancy", vacancyHandler.CreateVacancy)
		auth.PUT("/vacancy/:id", vacancyHandler.UpdateVacancy)
		auth.DELETE("/vacancy/:id", vacancyHandler.DeleteVacancy)
		auth.DELETE("/vacancy/archive/:id", vacancyHandler.SoftDelete)
		auth.GET("/vacancy/archive", vacancyHandler.GetArchiveVacancy)
		auth.GET("/vacancy/archive/:id", vacancyHandler.GetArchiveVacancyByID)
		auth.PUT("/vacancy/archive/:id", vacancyHandler.UnArchiveVacancy)
	}

	r.Run(":" + viper.GetString("port"))
}
