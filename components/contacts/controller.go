package contacts

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

const colName = "contacts"

// Get Contacts godoc
// @Summary Get recent contact
// @Description Get most recent contact
// @ID get-contacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Contacts}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /contacts [get]
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

	result := make([]Contacts, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Contacts godoc
// @Summary Find info contacts by ID
// @Description Find info contacts by ID
// @ID find-contacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the contact to get"
// @Success 200 {object} utils.HttpSuccess{data=Contacts}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /contacts/{id} [get]
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

	var record Contacts

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Contacts godoc
// @Summary Create an info for page contact
// @Description Create an info for page contact
// @ID create-contacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param address body string true "Contacts address"
// @Param phone body string true "Contacts phone"
// @Param mail body string true "Contacts mail"
// @Success 200 {object} Contacts
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /contacts [post]
func Create(c echo.Context) error {
	/* store record to db */
	contacts := &Contacts{
		ID:        primitive.NewObjectID(),
		Address:   c.FormValue("address"),
		Phone:     c.FormValue("phone"),
		Mail:      c.FormValue("mail"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	_, err = db.Collection(colName).InsertOne(ctx, contacts)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(contacts, "Saved"))
}

// Update Contacts godoc
// @Summary Update an info for page contact
// @Description Update an info for page contact
// @ID update-contact
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of contact to get"
// @Param address body string true "Contacts address"
// @Param phone body string true "Contacts phone"
// @Param mail body string true "Contacts mail"
// @Success 200 {object} Contacts
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /contacts/{id} [put]
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

	changes := &Contacts{
		Address: c.FormValue("address"),
		Phone:   c.FormValue("phone"),
		Mail:    c.FormValue("mail"),
	}

	update, err := db.Collection(colName).UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(update, "Updated"))
}

// Delete Contacts godoc
// @Summary Delete an contact info
// @Description Delete an contact info
// @ID delete-contact
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the contact"
// @Success 200 {object} Contacts
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /contacts/{id} [delete]
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
