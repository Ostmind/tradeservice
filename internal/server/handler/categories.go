package handlers

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"tradeservice/internal/models"
)

type CategoryManager interface {
	AddCategory(ctx context.Context, name string, productId string) (id string, err error)
	GetCategory(ctx context.Context) ([]models.CategoryDto, error)
	SetCategory(ctx context.Context, id string, name string) error
	DeleteCategory(ctx context.Context, id string) error
}

type CategoriesController struct {
	manager CategoryManager
	logger  *slog.Logger
}

func NewCategoriesHandler(manager CategoryManager, log *slog.Logger) *CategoriesController {
	return &CategoriesController{manager, log}
}

func (ctr CategoriesController) GetCategory(c echo.Context) error {

	ctr.logger.Debug("Get Request for Categories")

	res, err := ctr.manager.GetCategory(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr CategoriesController) AddCategory(c echo.Context) error {

	ctr.logger.Debug("Post Request for Categories")

	categoryName := c.Param("categoryName")

	productId := c.Param("productId")

	res, err := ctr.manager.AddCategory(c.Request().Context(), categoryName, productId)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr CategoriesController) DeleteCategory(c echo.Context) error {

	ctr.logger.Debug("Delete Request for Categories")

	categoryId := c.Param("categoryId")

	err := ctr.manager.DeleteCategory(c.Request().Context(), categoryId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (ctr CategoriesController) SetCategory(c echo.Context) error {

	ctr.logger.Debug("Patch Request for Categories")

	categoryId := c.Param("categoryId")

	categoryName := c.Param("categoryName")

	err := ctr.manager.SetCategory(c.Request().Context(), categoryId, categoryName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
