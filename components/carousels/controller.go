package carousels

import (
	"context"
	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var ctx = context.Background()

const colName = "carousels"

// Get Carousels godoc
// @Summary Get recent carousel
// @Contentription get most recent carousel
// @ID get-carousels
// @Tags Carousels
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Carousels}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /carousels [get]
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

	result := make([]Carousels, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Carousels godoc
// @Summary Find carousel by ID
// @Description Find carousel by ID
// @ID find-carousels
// @Tags Carousels
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the carousel to get"
// @Success 200 {object} utils.HttpSuccess{data=Carousels}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /carousels/{id} [get]
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

	var record Carousels

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Carousels godoc
// @Summary Create carousel
// @Description Create a carousel content
// @ID create-carousels
// @Tags Carousels
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param image formData file true "Blog image"
// @Param tagline formData string true "Blog tagline"
// @Param tagdesc formData string true "Blog tagdesc"
// @Success 200 {object} Carousels
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /carousels [post]
func Create(c echo.Context) error {
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
	path := "./uploaded_files/carousel/"
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
	carouselsRecord := &Carousels{
		ID:        primitive.NewObjectID(),
		Tagline:   c.FormValue("tagline"),
		Tagdesc:   c.FormValue("tagdesc"),
		Image:     filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Collection(colName).InsertOne(ctx, carouselsRecord)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(carouselsRecord, "Saved"))
}

// Update Blog godoc
// @Summary Update carousel
// @Description Update carousel
// @ID update-carousel
// @Tags Carousels
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of carousel to get"
// @Param image formData file false "Blog image"
// @Param tagline formData string false "Blog tagline"
// @Param tagdesc formData string false "Blog tagdesc"
// @Success 200 {object} Carousels
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /carousels/{id} [put]
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

	changes := &Carousels{
		Tagline: c.FormValue("tagline"),
		Tagdesc: c.FormValue("tagdesc"),
		Image:   "",
	}

	/* check image exist first */
	file, err := c.FormFile("image")
	// if no error then there is valid image request
	if err == nil {
		/* delete existing file if exist */
		path := "./uploaded_files/carousel/"
		var record Carousels

		if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
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

	update, err := db.Collection(colName).UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(update, "Updated"))

}

// Delete Blog godoc
// @Summary Delete a carousel
// @Description Delete a carousel
// @ID delete-carousel
// @Tags Carousels
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the carousel"
// @Success 200 {object} Carousels
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /carousels/{id} [delete]
func Destroy(c echo.Context) error {
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
	path := "./uploaded_files/carousel/"
	var record Carousels

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if record.Image != "" {
		err := os.Remove(path + record.Image)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	/* delete record */
	result, err := db.Collection(colName).DeleteOne(ctx, selector)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Deleted"))
}
