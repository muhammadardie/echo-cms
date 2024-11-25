package users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	result := make([]PublicUsers, 0)
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

	// Connect to the database
	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to DB")
	}

	// Parse the JSON body into the Users struct
	users := new(Users)
	if err := c.Bind(users); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	// Validate required fields
	if err := c.Validate(users); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Check for unique email
	if users.Email != "" {
		emailFilter := bson.M{
			"email": users.Email,
		}
		existingUser := db.Collection(colName).FindOne(ctx, emailFilter)
		if existingUser.Err() == nil { // Email already exists
			return echo.NewHTTPError(http.StatusConflict, "Email already exists")
		}
	}

	// Check if the password is provided
	if users.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Password is required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}
	users.Password = string(hashedPassword)

	// Set additional fields
	users.ID = primitive.NewObjectID()
	users.CreatedAt = time.Now()
	users.UpdatedAt = time.Now()

	// Insert the user into the database
	_, err = db.Collection(colName).InsertOne(ctx, users)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save user")
	}

	// Respond with the created user (excluding sensitive data like password)
	users.Password = "" // Avoid returning the password in the response
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

	changes := new(UpdateUser)

	if err := c.Bind(changes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}

	if err := c.Validate(changes); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Check for unique email
	if changes.Email != "" {
		emailFilter := bson.M{
			"email": changes.Email,
			"_id":   bson.M{"$ne": id}, // Exclude the current user
		}
		existingUser := db.Collection(colName).FindOne(ctx, emailFilter)
		if existingUser.Err() == nil { // Email already exists
			return echo.NewHTTPError(http.StatusConflict, "Email already exists")
		}
	}

	updateFields := bson.M{
		"username": changes.Username,
		"email":    changes.Email,
	}

	// Check if password is provided
	if changes.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changes.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
		}
		updateFields["password"] = string(hashedPassword)
	}

	selector := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	result, err := db.Collection(colName).UpdateOne(ctx, selector, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Updated successfully"))
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
