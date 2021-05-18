package socmeds

import (
	"context"
	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
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
	/* store record to db */
	socmeds := &Socmeds{
		ID:        primitive.NewObjectID(),
		Name:      c.FormValue("name"),
		Icon:      c.FormValue("icon"),
		Url:       c.FormValue("url"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	_, err = db.Collection(colName).InsertOne(ctx, socmeds)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(socmeds, "Saved"))
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

	selector := bson.M{"_id": id}

	changes := &Socmeds{
		Name: c.FormValue("name"),
		Icon: c.FormValue("icon"),
		Url:  c.FormValue("url"),
	}

	update, err := db.Collection(colName).UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(update, "Updated"))
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
