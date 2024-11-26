package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhammadardie/echo-cms/components/abouts"
	"github.com/muhammadardie/echo-cms/components/blogs"
	"github.com/muhammadardie/echo-cms/components/carousels"
	"github.com/muhammadardie/echo-cms/components/companies"
	"github.com/muhammadardie/echo-cms/components/contacts"
	"github.com/muhammadardie/echo-cms/components/galleries"
	"github.com/muhammadardie/echo-cms/components/headers"
	"github.com/muhammadardie/echo-cms/components/services"
	"github.com/muhammadardie/echo-cms/components/socmeds"
	"github.com/muhammadardie/echo-cms/components/teams"
	"github.com/muhammadardie/echo-cms/components/testimonies"
	"github.com/muhammadardie/echo-cms/components/users"
)

func Register(g *echo.Group) {
	g.GET("/csrf", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("csrf"))
	})
	abouts.AboutsRegister(g)
	blogs.BlogsRegister(g)
	carousels.CarouselsRegister(g)
	companies.CompaniesRegister(g)
	contacts.ContactsRegister(g)
	galleries.GalleriesRegister(g)
	headers.HeadersRegister(g)
	services.ServicesRegister(g)
	socmeds.SocmedsRegister(g)
	teams.TeamsRegister(g)
	testimonies.TestimoniesRegister(g)
	users.UsersRegister(g)
}

func RegisterPublic(r *echo.Echo) {
	publicGroup := r.Group("/api/public")

	// Public read-only routes
	publicGroup.GET("/abouts", abouts.Get)
	publicGroup.GET("/abouts/:id", abouts.Find)

	publicGroup.GET("/blogs", blogs.Get)
	publicGroup.GET("/blogs/:id", blogs.Find)

	publicGroup.GET("/carousels", carousels.Get)
	publicGroup.GET("/carousels/:id", carousels.Find)

	publicGroup.GET("/companies", companies.Get)
	publicGroup.GET("/companies/:id", companies.Find)

	publicGroup.GET("/contacts", contacts.Get)
	publicGroup.GET("/contacts/:id", contacts.Find)

	publicGroup.GET("/galleries", galleries.Get)
	publicGroup.GET("/galleries/:id", galleries.Find)

	publicGroup.GET("/headers", headers.Get)
	publicGroup.GET("/headers/:id", headers.Find)
	publicGroup.GET("/headers/page/:pagename", headers.FindByPage)

	publicGroup.GET("/services", services.Get)
	publicGroup.GET("/services/:id", services.Find)

	publicGroup.GET("/socmeds", socmeds.Get)
	publicGroup.GET("/socmeds/:id", socmeds.Find)

	publicGroup.GET("/teams", teams.Get)
	publicGroup.GET("/teams/:id", teams.Find)

	publicGroup.GET("/testimonies", testimonies.Get)
	publicGroup.GET("/testimonies/:id", testimonies.Find)
}
