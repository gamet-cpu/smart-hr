package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	handler "smart-hr/handlers"
	"smart-hr/middleware"
	"smart-hr/repository"
)

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

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user/me", userHandler.GetMe)
		auth.GET("/users", userHandler.GetAll)
		auth.PUT("/user", userHandler.Update)
		auth.DELETE("/user", userHandler.Delete)
	}

	r.Run(":" + viper.GetString("port"))
}
