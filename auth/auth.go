package auth

import (
	"context"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/muhammadardie/echo-cms/components/users"
	DB "github.com/muhammadardie/echo-cms/db"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var ctx = context.Background()

func Login(c echo.Context) error {
	ctx := context.Background()
	user := new(users.Users)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate input
	if err := c.Validate(user); err != nil {
		return err
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	selector := bson.M{"email": user.Email}
	var dbUser users.Users

	if err = db.Collection("users").FindOne(ctx, selector).Decode(&dbUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Users not found")
	}

	// Comparing the password with the hash
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)
	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid password")
	}

	ts, err := CreateToken(dbUser.ID.Hex())
	if err != nil {
 		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	saveErr := CreateAuth(dbUser.ID.Hex(), ts)
	if saveErr != nil {
	    return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	tokens := map[string]string{
	    "access_token":  ts.AccessToken,
	    "refresh_token": ts.RefreshToken,
	}

	return c.JSON(http.StatusOK, tokens)
}

func Logout(c echo.Context) error {
  au, err := ExtractTokenMetadata(c)

  if err != nil {
    return c.JSON(http.StatusUnauthorized, err.Error())
  }

  delErr := DeleteAuth(au.AccessUuid)
  if delErr != nil { //if any goes wrong
    return c.JSON(http.StatusUnauthorized, delErr)
  }

  return c.JSON(http.StatusOK, "Successfully logged out")
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
	if err        != nil {
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