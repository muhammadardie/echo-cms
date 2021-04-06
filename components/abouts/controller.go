package abouts

import (
	"context"
	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Get Abouts godoc
// @Summary Get recent info about
// @Description get most recent info about
// @ID get-abouts
// @Tags Abouts
// @Accept  json
// @Produce  json
// @Success 200 {object} Abouts
// @Router /abouts [get]
func Get(c echo.Context) error {
	var ctx = context.Background()

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	csr, err := db.Collection("abouts").Find(ctx, bson.M{})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	defer csr.Close(ctx)

	result := make([]Abouts, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// Find Abouts godoc
// @Summary Find info abouts by ID
// @Description Find info abouts by ID
// @ID find-abouts
// @Tags Abouts
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the about to get"
// @Success 200 {object} Abouts
// @Router /abouts/{id} [get]
func Find(c echo.Context) error {
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

	var record Abouts

	if err = db.Collection("abouts").FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, record)
}

// Create About godoc
// @Summary Create an info for page about
// @Description Create an info for page about
// @ID create-abouts
// @Tags Abouts
// @Accept  json
// @Produce  json
// @Param about body Abouts true "Info about you want to create"
// @Success 200 {object} Abouts
// @Success 200 {object} []string{title=string,image=string,desc=string}
// @Router /abouts [post]
func Create(c echo.Context) error {
	var ctx = context.Background()

	/* upload image first */
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	db, err := DB.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	path := "./uploaded_files/about/"
	name := file.Filename
	extension := filepath.Ext(name)
	filename := xid.New().String() + extension

	dst, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, 0666)
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
		Image:     filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Collection("abouts").InsertOne(ctx, abouts)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, abouts)
}

// Update About godoc
// @Summary Update an info for page about
// @Description Update an info for page about
// @ID update-about
// @Tags Abouts
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the about to get"
// @Param about body Abouts true "Info about you want to create"
// @Success 200 {object} Abouts
// @Router /abouts/{id} [put]
func Update(c echo.Context) error {
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

	changes := &Abouts{
		Title: c.FormValue("title"),
		Desc:  c.FormValue("desc"),
		Image: "",
	}

	/* check image exist first */
	file, err := c.FormFile("image")
	// if no error then there is valid image request
	if err == nil {
		/* delete existing file if exist */
		path := "./uploaded_files/about/"
		var record Abouts

		if err = db.Collection("abouts").FindOne(ctx, selector).Decode(&record); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if record.Image != "" {
			err := os.Remove(path + record.Image)

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}

		/* upload new file */
		name := file.Filename
		extension := filepath.Ext(name)
		filename := xid.New().String() + extension

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		changes.Image = filename
	}

	_, err = db.Collection("abouts").UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, changes)

}

// Delete Abouts godoc
// @Summary Delete an about info
// @Description Delete an about info
// @ID delete-about
// @Tags Abouts
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the about"
// @Success 200 {object} Abouts
// @Router /abouts/{id} [delete]
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

	/* delete any exist image */
	path := "./uploaded_files/about/"
	var record Abouts

	if err = db.Collection("abouts").FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if record.Image != "" {
		err := os.Remove(path + record.Image)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	/* delete record */
	result, err := db.Collection("abouts").DeleteOne(ctx, selector)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, result)
}
