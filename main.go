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
	meetRepo := repository.NewMeetRepository(dbpool)
	meetHandler := handler.NewMeetHandler(meetRepo)
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/message", messageHandler.SentReq)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user/me", userHandler.GetMe)                                //Получение данных о пользователе
		auth.GET("/users", userHandler.GetAll)                                 //Получение данных о всех пользователях
		auth.PUT("/user", userHandler.Update)                                  //Обновление данных пользователя
		auth.DELETE("/user", userHandler.Delete)                               //Удаление пользователя
		auth.PUT("/user/:id", userHandler.RespondVacancy)                      //Откликтнутся на вакансию
		auth.GET("/users/myresponses", userHandler.MyResponses)                //Вакансии на которые откликнулся пользователь
		auth.GET("/users/responses/:id", userHandler.GetResponses)             //Пользователи откликнувшиеся на вакансию
		auth.GET("/vacancy", vacancyHandler.GetAllVacancy)                     //Список всех вакансий
		auth.GET("/vacancy/:id", vacancyHandler.GetVacancyByID)                //Конкретная вакансия по ID
		auth.POST("/vacancy", vacancyHandler.CreateVacancy)                    //Создать вакансию
		auth.PUT("/vacancy/:id", vacancyHandler.UpdateVacancy)                 //Обновить данные о вакансии
		auth.DELETE("/vacancy/:id", vacancyHandler.DeleteVacancy)              //Удалить вакансию полностью
		auth.DELETE("/vacancy/archive/:id", vacancyHandler.SoftDelete)         //Удаление вакансии с переносом в архив
		auth.GET("/vacancy/archive", vacancyHandler.GetArchiveVacancy)         //Список архивированных вакансий
		auth.GET("/vacancy/archive/:id", vacancyHandler.GetArchiveVacancyByID) //Архивированная вакансия по ID
		auth.PUT("/vacancy/archive/:id", vacancyHandler.UnArchiveVacancy)      //Деархивация вакансии
		auth.GET("/vacancy/topvacancy", vacancyHandler.TopVacancy)             //Самые активные 5 вакансий
		auth.GET("/vacancy/actualvacancy", vacancyHandler.ActualVacancy)       //Новые вакансии
		auth.GET("/meets", meetHandler.GetAllMeets)                            //Все встречи пользователя
		auth.GET("/meets/:id", meetHandler.GetMeetById)                        //Встреча по id
		auth.POST("/meets", meetHandler.CreateMeet)                            //Назначить встречу
		auth.PUT("/meets/:id", meetHandler.UpdateMeet)                         //Изменить встречу
		auth.DELETE("/meets/:id", meetHandler.DeleteMeet)                      //Удалить встречу
	}

	r.Run(":" + viper.GetString("port"))
}
