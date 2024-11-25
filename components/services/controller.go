package services

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

const colName = "services"

// Get Services godoc
// @Summary Get recent service
// @Description Get most recent service
// @ID get-services
// @Tags Services
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Services}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /services [get]
func Get(c echo.Context) error {
	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	csr, err := db.Collection(colName).Find(ctx, bson.M{})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	defer csr.Close(ctx)

	result := make([]Services, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Services godoc
// @Summary Find info services by ID
// @Description Find info services by ID
// @ID find-services
// @Tags Services
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the service to get"
// @Success 200 {object} utils.HttpSuccess{data=Services}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /services/{id} [get]
func Find(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	selector := bson.M{"_id": id}

	var record Services

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Services godoc
// @Summary Create an info for page service
// @Description Create an info for page service
// @ID create-services
// @Tags Services
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param title body string true "Services title"
// @Param icon body string true "Services icon"
// @Param desc body string true "Services desc"
// @Success 200 {object} Services
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /services [post]
func Create(c echo.Context) error {
	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	// Parse the JSON body into the services struct
	service := new(Services)
	if err := c.Bind(service); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	// Validate required fields
	if err := c.Validate(service); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Set additional fields
	service.ID = primitive.NewObjectID()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	// Insert the user into the database
	_, err = db.Collection(colName).InsertOne(ctx, service)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save service")
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(service, "Saved"))
}

// Update Services godoc
// @Summary Update an info for page service
// @Description Update an info for page service
// @ID update-service
// @Tags Services
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of service to get"
// @Param title body string true "Services title"
// @Param icon body string true "Services icon"
// @Param desc body string true "Services desc"
// @Success 200 {object} Services
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /services/{id} [put]
func Update(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	changes := new(Services)

	if err := c.Bind(changes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	updateFields := bson.M{
		"title": changes.Title,
		"icon":  changes.Icon,
		"desc":  changes.Desc,
	}

	selector := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	result, err := db.Collection(colName).UpdateOne(ctx, selector, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update service")
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Updated"))
}

// Delete Services godoc
// @Summary Delete an service info
// @Description Delete an service info
// @ID delete-service
// @Tags Services
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the service"
// @Success 200 {object} Services
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /services/{id} [delete]
func Destroy(c echo.Context) error {
	ctx := context.Background()
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	selector := bson.M{"_id": id}

	/* delete record */
	result, err := db.Collection(colName).DeleteOne(ctx, selector)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Deleted"))
}
