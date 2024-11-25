package contacts

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
	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	// Parse the JSON body into the contacts struct
	contact := new(Contacts)
	if err := c.Bind(contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	// Validate required fields
	if err := c.Validate(contact); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Set additional fields
	contact.ID = primitive.NewObjectID()
	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()

	_, err = db.Collection(colName).InsertOne(ctx, contact)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(contact, "Saved"))
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

	changes := new(Contacts)

	if err := c.Bind(changes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	updateFields := bson.M{
		"address": changes.Address,
		"phone":   changes.Phone,
		"mail":    changes.Mail,
	}

	selector := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	result, err := db.Collection(colName).UpdateOne(ctx, selector, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Updated"))
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
