package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"tradeservice/internal/models"
)

type CategoryManager interface {
	Add(ctx context.Context, name string, productId string) (id string, err error)
	Get(ctx context.Context) ([]models.Category, error)
	Set(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}

type CategoriesController struct {
	manager CategoryManager
	logger  *slog.Logger
}

func NewCategoriesHandler(manager CategoryManager, log *slog.Logger) *CategoriesController {
	return &CategoriesController{manager, log}
}

func (ctr CategoriesController) Get(c echo.Context) error {

	ctr.logger.Debug("Get Request for Categories")

	res, err := ctr.manager.Get(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr CategoriesController) Add(c echo.Context) error {

	ctr.logger.Debug("Post Request for Categories")

	categoryName := c.Param("categoryName")

	productId := c.Param("productId")

	res, err := ctr.manager.Add(c.Request().Context(), categoryName, productId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr CategoriesController) Delete(c echo.Context) error {

	ctr.logger.Debug("Delete Request for Categories")

	categoryId := c.Param("categoryId")

	err := ctr.manager.Delete(c.Request().Context(), categoryId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (ctr CategoriesController) Set(c echo.Context) error {

	ctr.logger.Debug("Patch Request for Categories")

	categoryId := c.Param("categoryId")

	categoryName := c.Param("categoryName")

	err := ctr.manager.Set(c.Request().Context(), categoryId, categoryName)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
