package abouts

import (
	"net/http"
	"context"
	"log"
	"os"
	"io"
	"time"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func Get(c echo.Context) error {
	var ctx = context.Background()

	db, err := DB.Connect()
    if err != nil {
        log.Fatal(err.Error())
    }

    csr, err := db.Collection("abouts").Find(ctx, bson.M{})
    if err != nil {
        log.Fatal(err.Error())
    }
    
    defer csr.Close(ctx)

    result := make([]Abouts, 0)
	if err = csr.All(ctx, &result); err != nil {
	    log.Fatal(err)
	}

	return c.JSON(http.StatusOK, result)
}

func Create(c echo.Context) error {
	var ctx = context.Background()

	/* upload image first */
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	db, err := DB.Connect()
    if err != nil {
        log.Fatal(err.Error())
    }

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	path := "./uploaded_files/about/"
	dst, err := os.OpenFile(path + file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	/* end upload image */

	/* store record to db */
	abouts := &Abouts{
    	ID:        primitive.NewObjectID(),
	    Title:     c.FormValue("title"),
	    Desc:      c.FormValue("desc"),
	    Image:     file.Filename,
	    CreatedAt: time.Now(),
	    UpdatedAt: time.Now(),
    }

	_, err = db.Collection("abouts").InsertOne(ctx, abouts)

    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err)
    }

	return c.JSON(http.StatusOK, abouts)
}