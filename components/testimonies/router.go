package testimonies

import (
	"github.com/labstack/echo/v4"
)

func TestimoniesRegister(g *echo.Group) {
	testimonies := g.Group("/testimonies")
	testimonies.GET("", Get)
	testimonies.POST("", Create)
	testimonies.GET("/:id", Find)
	testimonies.PUT("/:id", Update)
	testimonies.DELETE("/:id", Destroy)
}
