package socmeds

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

const colName = "socmeds"

// Get Socmeds godoc
// @Summary Get recent socmeds
// @Description Get most recent socmeds
// @ID get-socmeds
// @Tags Socmeds
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Socmeds}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /socmeds [get]
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

	result := make([]Socmeds, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Socmeds godoc
// @Summary Find info socmeds by ID
// @Description Find info socmeds by ID
// @ID find-socmeds
// @Tags Socmeds
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the socmeds to get"
// @Success 200 {object} utils.HttpSuccess{data=Socmeds}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /socmeds/{id} [get]
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

	var record Socmeds

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Socmeds godoc
// @Summary Create an info for page socmeds
// @Description Create an info for page socmeds
// @ID create-socmeds
// @Tags Socmeds
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param name body string true "Socmeds name"
// @Param icon body string true "Socmeds icon"
// @Param desc body string true "Socmeds desc"
// @Success 200 {object} Socmeds
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /socmeds [post]
func Create(c echo.Context) error {
	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	// Parse the JSON body into the socmeds struct
	socmed := new(Socmeds)
	if err := c.Bind(socmed); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	// Validate required fields
	if err := c.Validate(socmed); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Set additional fields
	socmed.ID = primitive.NewObjectID()
	socmed.CreatedAt = time.Now()
	socmed.UpdatedAt = time.Now()

	_, err = db.Collection(colName).InsertOne(ctx, socmed)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(socmed, "Saved"))
}

// Update Socmeds godoc
// @Summary Update an info for page socmeds
// @Description Update an info for page socmeds
// @ID update-socmeds
// @Tags Socmeds
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of socmeds to get"
// @Param name body string true "Socmeds name"
// @Param icon body string true "Socmeds icon"
// @Param url body string true "Socmeds url"
// @Success 200 {object} Socmeds
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /socmeds/{id} [put]
func Update(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	changes := new(Socmeds)

	if err := c.Bind(changes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	updateFields := bson.M{
		"name": changes.Name,
		"icon": changes.Icon,
		"url":  changes.Url,
	}

	selector := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	result, err := db.Collection(colName).UpdateOne(ctx, selector, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Updated"))
}

// Delete Socmeds godoc
// @Summary Delete an socmeds info
// @Description Delete an socmeds info
// @ID delete-socmeds
// @Tags Socmeds
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the socmeds"
// @Success 200 {object} Socmeds
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /socmeds/{id} [delete]
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
