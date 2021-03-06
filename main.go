package main

import (
	"github.com/muhammadardie/echo-cms/auth"
	DB "github.com/muhammadardie/echo-cms/db"
	_ "github.com/muhammadardie/echo-cms/docs" // docs generated by Swag CLI
	"github.com/muhammadardie/echo-cms/middleware"
	"github.com/muhammadardie/echo-cms/routes"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
	"os"
)

// @title Swagger Example API
// @version 1.0
// @description Conduit API
// @title Conduit API

// @host 127.0.0.1:8080
// @BasePath /api

// @schemes http https
// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func init() {
	DB.InitRedis()
}

func main() {
	r := middleware.New()

	// swagger
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	g := r.Group("/api")
	g.POST("/login", auth.Login)
	g.POST("/logout", auth.Logout)
	g.POST("/token/refresh", auth.Refresh)
	g.Use(middleware.TokenAuthMiddleware)

	routes.Register(g)

	r.Logger.Fatal(r.Start(":" + os.Getenv("PORT")))
}
