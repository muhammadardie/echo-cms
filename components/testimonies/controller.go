package testimonies

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

const colName = "testimonies"

// Get Testimony godoc
// @Summary Get recent testimony
// @Contentription get most recent testimony
// @ID get-testimonies
// @Tags Testimonies
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Testimonies}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /testimonies [get]
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

	result := make([]Testimonies, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Testimony godoc
// @Summary Find testimony by ID
// @Description Find testimony by ID
// @ID find-testimonies
// @Tags Testimonies
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the testimony to get"
// @Success 200 {object} utils.HttpSuccess{data=Testimonies}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /testimonies/{id} [get]
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

	var record Testimonies

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Testimony godoc
// @Summary Create testimony
// @Description Create a testimony content
// @ID create-testimonies
// @Tags Testimonies
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param avatar formData file true "Testimony avatar"
// @Param comment formData string true "Testimony comment"
// @Param username formData string true "Testimony username"
// @Success 200 {object} Testimonies
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /testimonies [post]
func Create(c echo.Context) error {
	/* upload image first */
	file, err := c.FormFile("avatar")
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
	path := "./uploaded_files/testimony/"
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
	testimoniesRecord := &Testimonies{
		ID:        primitive.NewObjectID(),
		Username:  c.FormValue("username"),
		Comment:   c.FormValue("comment"),
		Avatar:    filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Collection(colName).InsertOne(ctx, testimoniesRecord)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(testimoniesRecord, "Saved"))
}

// Update Testimony godoc
// @Summary Update testimony
// @Description Update testimony
// @ID update-testimony
// @Tags Testimony
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of testimony to get"
// @Param avatar formData file false "Testimony avatar"
// @Param comment formData string false "Testimony comment"
// @Param username formData string false "Testimony username"
// @Success 200 {object} Testimonies
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /testimonies/{id} [put]
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

	changes := &Testimonies{
		Username: c.FormValue("username"),
		Comment:  c.FormValue("comment"),
		Avatar:   "",
	}

	/* check image exist first */
	file, err := c.FormFile("avatar")
	// if no error then there is valid image request
	if err == nil {
		/* delete existing file if exist */
		path := "./uploaded_files/testimony/"
		var record Testimonies

		if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if record.Avatar != "" {
			err := os.Remove(path + record.Avatar)

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

		changes.Avatar = filename
	}

	update, err := db.Collection(colName).UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(update, "Updated"))

}

// Delete Testimony godoc
// @Summary Delete a testimony
// @Description Delete a testimony
// @ID delete-testimony
// @Tags Testimony
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the testimony"
// @Success 200 {object} Testimonies
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /testimonies/{id} [delete]
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
	path := "./uploaded_files/testimony/"
	var record Testimonies

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if record.Avatar != "" {
		err := os.Remove(path + record.Avatar)

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
