package users

import (
	"github.com/labstack/echo/v4"
)

func UsersRegister(g *echo.Group) {
	users := g.Group("/users")
	users.GET("", Get)
	users.POST("", Create)
	users.GET("/:id", Find)
	users.PUT("/:id", Update)
	users.DELETE("/:id", Destroy)
}
