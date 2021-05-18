package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/muhammadardie/echo-cms/components/users"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var ctx = context.Background()

type Token struct {
	ID            primitive.ObjectID 	`json:"_id"`
	Username      string   				`json:"username"`
	Email     	  string   				`json:"email"`
	AccessToken   string   				`json:"access_token"`
	RefreshToken  string   				`json:"refresh_token"`
}

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body users.UserLogin true "Credentials to use"
// @Success 200 {object} utils.HttpSuccess{data=string{_id=string,username=string,email=string,access_token=string,refresh_token=string}}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Failure 500 {object} utils.HttpError
// @Router /login [post]
func Login(c echo.Context) error {
	ctx := context.Background()
	user := new(users.UserLogin)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate input
	if err := c.Validate(user); err != nil {
		return err
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	selector := bson.M{"email": user.Email}
	var dbUser users.Users

	if err = db.Collection("users").FindOne(ctx, selector).Decode(&dbUser); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}

	// Comparing the password with the hash
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)
	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password")
	}

	ts, err := CreateToken(dbUser.ID.Hex())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	saveErr := CreateAuth(dbUser.ID.Hex(), ts)
	if saveErr != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	tokens := &Token{
		ID: dbUser.ID,
		Username: dbUser.Username,
		Email: dbUser.Email,
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(tokens, "Successfully logged in"))
}

// Logout godoc
// @Summary Logout for existing user
// @Description Logout for existing user
// @ID logout
// @Tags Auth
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess
// @Failure 401 {object} utils.HttpError
// @Router /logout [post]
func Logout(c echo.Context) error {
	au, err := ExtractTokenMetadata(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil { //if any goes wrong
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess("", "Successfully logged out"))
}

func CreateAuth(userid string, td *TokenDetails) error {
	client := DB.InitRedis()
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(ctx, td.TokenUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(ctx, td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func FetchAuth(authD *AccessDetails) (string, error) {
	client := DB.InitRedis()
	userid, err := client.Get(ctx, authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}

	return userid, nil
}

func DeleteAuth(accessUuid string) error {
	client := DB.InitRedis()
	err := client.Del(ctx, accessUuid).Err()
	if err != nil {
		return err
	}

	return nil
}
