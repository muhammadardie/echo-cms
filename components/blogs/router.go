package blogs

import (
	"github.com/labstack/echo/v4"
)

func BlogsRegister(g *echo.Group) {
	blogs := g.Group("/blogs")
	blogs.GET("", Get)
	blogs.POST("", Create)
	blogs.GET("/:id", Find)
	blogs.PUT("/:id", Update)
	blogs.DELETE("/:id", Destroy)
}
