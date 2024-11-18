package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/muhammadardie/echo-cms/auth"
	"github.com/muhammadardie/echo-cms/utils"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = ErrorHandler

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = string(log.ERROR) // Default to ERROR level if not set
	}

	// Set the log level based on the environment variable
	switch logLevel {
	case "DEBUG":
		e.Logger.SetLevel(log.DEBUG)
	case "INFO":
		e.Logger.SetLevel(log.INFO)
	case "WARN":
		e.Logger.SetLevel(log.WARN)
	case "ERROR":
		e.Logger.SetLevel(log.ERROR)
	default:
		e.Logger.SetLevel(log.ERROR) // Default to ERROR if invalid value
	}

	e.Logger.SetLevel(log.ERROR)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	return e
}

func ErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				report.Message = fmt.Sprintf("%s is required",
					err.Field())
			case "email":
				report.Message = fmt.Sprintf("%s is not valid email",
					err.Field())
			case "gte":
				report.Message = fmt.Sprintf("%s value must be greater than %s",
					err.Field(), err.Param())
			case "lte":
				report.Message = fmt.Sprintf("%s value must be lower than %s",
					err.Field(), err.Param())
			}

			break
		}
	}

	message := fmt.Sprintf("%v", report.Message)

	c.Logger().Error(report)
	c.JSON(report.Code, utils.NewError(report.Code, message))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	const unauthorizedMessage = "Error: Access Token is not valid or has expired"

	return func(c echo.Context) error {
		err := auth.TokenValid(c) // check jwt still valid
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, unauthorizedMessage)
		}

		tokenAuth, err := auth.ExtractTokenMetadata(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, unauthorizedMessage)
		}

		_, err = auth.FetchAuth(tokenAuth) // check jwt still exist in redis
		if err != nil {
			return c.JSON(http.StatusUnauthorized, unauthorizedMessage)
		}

		return next(c)
	}
}
