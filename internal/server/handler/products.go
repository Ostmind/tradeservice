package handlers

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"tradeservice/internal/models"
)

type ProductManager interface {
	Add(ctx context.Context, name string, productId string) (id string, err error)
	Get(ctx context.Context) ([]models.Category, error)
	Set(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}

type ProductController struct {
	manager ProductManager
	logger  *slog.Logger
}

func NewProductHandler(manager ProductManager, log *slog.Logger) *ProductController {
	return &ProductController{manager, log}
}

func (ctr ProductController) Get(c echo.Context) error {

	ctr.logger.Debug("Get Request for Products")

	res, err := ctr.manager.Get(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr ProductController) Add(c echo.Context) error {

	ctr.logger.Debug("Post Request for Products")

	productName := c.Param("productName")

	productId := c.Param("productId")

	res, err := ctr.manager.Add(c.Request().Context(), productName, productId)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr ProductController) Delete(c echo.Context) error {

	ctr.logger.Debug("Delete Request for Products")

	productId := c.Param("productId")

	err := ctr.manager.Delete(c.Request().Context(), productId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (ctr ProductController) Set(c echo.Context) error {

	ctr.logger.Debug("Patch Request for Products")

	productId := c.Param("productId")

	productName := c.Param("productName")

	err := ctr.manager.Set(c.Request().Context(), productId, productName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
