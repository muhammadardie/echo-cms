package headers

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	DB "github.com/muhammadardie/echo-cms/db"
	"github.com/muhammadardie/echo-cms/utils"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

const colName = "headers"

// Get Headers godoc
// @Summary Get recent header
// @Contentription get most recent header
// @ID get-headers
// @Tags Headers
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Headers}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /headers [get]
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

	result := make([]Headers, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Headers godoc
// @Summary Find header by ID
// @Description Find header by ID
// @ID find-headers
// @Tags Headers
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the header to get"
// @Success 200 {object} utils.HttpSuccess{data=Headers}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /headers/{id} [get]
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

	var record Headers

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Headers godoc
// @Summary Create header
// @Description Create a header content
// @ID create-headers
// @Tags Headers
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param image formData file true "Header image"
// @Param page formData string true "Header page"
// @Param tagline formData string true "Header tagline"
// @Param tagdesc formData string true "Header tagdesc"
// @Success 200 {object} Headers
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /headers [post]
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
	path := "./uploaded_files/header/"
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
	headersRecord := &Headers{
		ID:        primitive.NewObjectID(),
		Page:      c.FormValue("page"),
		Tagline:   c.FormValue("tagline"),
		Tagdesc:   c.FormValue("tagdesc"),
		Image:     filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Collection(colName).InsertOne(ctx, headersRecord)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(headersRecord, "Saved"))
}

// Update Header godoc
// @Summary Update header
// @Description Update header
// @ID update-header
// @Tags Headers
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of header to get"
// @Param image formData file true "Header image"
// @Param page formData string true "Header page"
// @Param tagline formData string true "Header tagline"
// @Param tagdesc formData string true "Header tagdesc"
// @Success 200 {object} Headers
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /headers/{id} [put]
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

	changes := &Headers{
		Page:      c.FormValue("page"),
		Tagline:   c.FormValue("tagline"),
		Tagdesc:   c.FormValue("tagdesc"),
		Image:     "",
		UpdatedAt: time.Now(),
	}

	/* check image exist first */
	file, err := c.FormFile("image")
	// if no error then there is valid image request
	if err == nil {
		/* delete existing file if exist */
		path := "./uploaded_files/header/"
		var record Headers

		if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if record.Image != "" {
			filePath := path + record.Image

			// Ensure the file path exists before attempting deletion
			if _, err := os.Stat(filePath); err == nil {
				// File exists, proceed to delete
				err := os.Remove(filePath)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Failed to delete old file: "+err.Error())
				}
			} else if os.IsNotExist(err) {
				// File does not exist, skip deletion
			} else {
				// Other errors
				return echo.NewHTTPError(http.StatusInternalServerError, "Error checking file: "+err.Error())
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

// Delete Header godoc
// @Summary Delete a header
// @Description Delete a header
// @ID delete-header
// @Tags Headers
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the header"
// @Success 200 {object} Headers
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /headers/{id} [delete]
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
	path := "./uploaded_files/header/"
	var record Headers

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if record.Image != "" {
		filePath := path + record.Image

		// Check if the file exists
		if _, err := os.Stat(filePath); err == nil {
			// File exists, proceed to delete
			err := os.Remove(filePath)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}
	}

	/* delete record */
	result, err := db.Collection(colName).DeleteOne(ctx, selector)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, "Deleted"))
}
