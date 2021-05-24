package galleries

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

const colName = "galleries"

// Get Galleries godoc
// @Summary Get recent info about
// @Description Get most recent info about
// @ID get-galleries
// @Tags Galleries
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} utils.HttpSuccess{data=[]Galleries}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /galleries [get]
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

	result := make([]Galleries, 0)
	if err = csr.All(ctx, &result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(result, ""))
}

// Find Galleries godoc
// @Summary Find info galleries by ID
// @Description Find info galleries by ID
// @ID find-galleries
// @Tags Galleries
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of the about to get"
// @Success 200 {object} utils.HttpSuccess{data=Galleries}
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /galleries/{id} [get]
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

	var record Galleries

	if err = db.Collection(colName).FindOne(ctx, selector).Decode(&record); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(record, ""))
}

// Create Gallery godoc
// @Summary Create a gallery
// @Description Create a gallery
// @ID create-galleries
// @Tags Galleries
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param image formData file true "Galleries image"
// @Param url formData string false "Galleries url"
// @Param title formData string true "Galleries title"
// @Param desc formData string false "Galleries description"
// @Success 200 {object} Galleries
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /galleries [post]
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
	path := "./uploaded_files/gallery/"
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
	// end upload image

	/* store record to db */
	galleries := &Galleries{
		ID:        primitive.NewObjectID(),
		Title:     c.FormValue("title"),
		Url:       c.FormValue("url"),
		Desc:      c.FormValue("desc"),
		Image:     filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Collection(colName).InsertOne(ctx, galleries)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, utils.NewSuccess(galleries, "Saved"))
}

// Update About godoc
// @Summary Update an info for page about
// @Description Update an info for page about
// @ID update-about
// @Tags Galleries
// @Accept  mpfd
// @Produce  json
// @Security Bearer
// @Param id path string true "ID of blog to get"
// @Param image formData file false "Galleries image"
// @Param url formData string false "Galleries url"
// @Param title formData string false "Galleries title"
// @Param desc formData string false "Galleries description"
// @Success 200 {object} Galleries
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /galleries/{id} [put]
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

	changes := &Galleries{
		Title: c.FormValue("title"),
		Url:   c.FormValue("url"),
		Desc:  c.FormValue("desc"),
		Image: "",
	}

	/* check image exist first */
	file, err := c.FormFile("image")
	// if no error then there is valid image request
	if err == nil {
		/* delete existing file if exist */
		path := "./uploaded_files/gallery/"
		var record Galleries

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

// Delete Gallerie godoc
// @Summary Delete a gallery
// @Description Delete a gallery
// @ID delete-gallery
// @Tags Galleries
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the gallery"
// @Success 200 {object} Galleries
// @Failure 400 {object} utils.HttpError
// @Failure 401 {object} utils.HttpError
// @Router /galleries/{id} [delete]
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
	path := "./uploaded_files/gallery/"
	var record Galleries

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
