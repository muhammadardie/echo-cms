package headers

import (
	"github.com/labstack/echo/v4"
)

func HeadersRegister(g *echo.Group) {
	headers := g.Group("/headers")
	headers.GET("/page/:pagename", FindByPage)
	headers.GET("", Get)
	headers.POST("", Create)
	headers.GET("/:id", Find)
	headers.PUT("/:id", Update)
	headers.DELETE("/:id", Destroy)
}
