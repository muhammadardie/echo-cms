package contacts

import (
	"github.com/labstack/echo/v4"
)

func ContactsRegister(g *echo.Group) {
	contacts := g.Group("/contacts")
	contacts.GET("", Get)
	contacts.POST("", Create)
	contacts.GET("/:id", Find)
	contacts.PUT("/:id", Update)
	contacts.DELETE("/:id", Destroy)

}
