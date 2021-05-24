package users

import (
	"context"
	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var ctx = context.Background()

const colName = "users"

// Get Users godoc
// @Summary Get recent user
// @Description Get most recent user
// @ID get-users
// @Tags Users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Users}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /users [get]
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

	result := make([]Users, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Users godoc
// @Summary Find info users by ID
// @Description Find info users by ID
// @ID find-users
// @Tags Users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the user to get"
// @Success 200 {object} utils.HttpSuccess{data=Users}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /users/{id} [get]
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

	var record Users

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Users godoc
// @Summary Create an info for page user
// @Description Create an info for page user
// @ID create-users
// @Tags Users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param username body string true "Users username"
// @Param password body string true "Users password"
// @Param email body string true "Users email"
// @Success 200 {object} Users
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Failure 500 {object} utils.HttpError
// @Router /users [post]
func Create(c echo.Context) error {

	/* check password */
	passValue := c.FormValue("password")
	if passValue == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Password is required")
	}

	/* hash password */
	password := []byte(passValue)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	/* store record to db */
	users := &Users{
		ID:        primitive.NewObjectID(),
		Username:  c.FormValue("username"),
		Password:  string(hashedPassword),
		Email:     c.FormValue("email"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.Validate(users); err != nil {
		return err
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Fail to connect DB")
	}

	_, err = db.Collection(colName).InsertOne(ctx, users)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(users, "Saved"))
}

// Update Users godoc
// @Summary Update an info for page user
// @Description Update an info for page user
// @ID update-user
// @Tags Users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of user to get"
// @Param username body string true "Users username"
// @Param password body string true "Users password"
// @Param email body string true "Users email"
// @Success 200 {object} Users
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Failure 500 {object} utils.HttpError
// @Router /users/{id} [put]
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

	/* hash password */
	password := []byte(c.FormValue("password"))
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	changes := &Users{
		Username: c.FormValue("username"),
		Password: string(hashedPassword),
		Email:    c.FormValue("email"),
	}

	if err := c.Validate(changes); err != nil {
		return err
	}

	update, err := db.Collection(colName).UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(update, "Updated"))
}

// Delete Users godoc
// @Summary Delete an user info
// @Description Delete an user info
// @ID delete-user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the user"
// @Success 200 {object} Users
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /users/{id} [delete]
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
