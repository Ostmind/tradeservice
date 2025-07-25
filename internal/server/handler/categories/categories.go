package categories

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"tradeservice/internal/models"

	"github.com/labstack/echo/v4"
)

//go:generate mockgen -source=categories.go -destination=mock/categoriesrepository.go

type CategoryManager interface {
	AddCategory(ctx context.Context, name string, productID string) (ID string, err error)
	GetCategory(ctx context.Context) ([]models.CategoryDto, error)
	SetCategory(ctx context.Context, ID string, name string) error
	DeleteCategory(ctx context.Context, ID string) error
}

type CategoriesController struct {
	manager CategoryManager
	logger  *slog.Logger
}

func NewCategoriesHandler(manager CategoryManager, log *slog.Logger) *CategoriesController {
	return &CategoriesController{manager, log}
}

func (ctr CategoriesController) GetCategory(echo echo.Context) error {
	ctr.logger.Debug("Get Request for Categories")

	res, err := ctr.manager.GetCategory(echo.Request().Context())
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr CategoriesController) AddCategory(echo echo.Context) error {
	ctr.logger.Debug("Post Request for Categories")

	categoryName := echo.Param("categoryName")

	productID := echo.Param("productId")

	res, err := ctr.manager.AddCategory(echo.Request().Context(), categoryName, productID)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr CategoriesController) DeleteCategory(echo echo.Context) error {
	ctr.logger.Debug("Delete Request for Categories")

	categoryID := echo.Param("categoryId")

	err := ctr.manager.DeleteCategory(echo.Request().Context(), categoryID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}

func (ctr CategoriesController) SetCategory(echo echo.Context) error {
	ctr.logger.Debug("Patch Request for Categories")

	categoryID := echo.Param("categoryId")

	categoryName := echo.Param("categoryName")

	err := ctr.manager.SetCategory(echo.Request().Context(), categoryID, categoryName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}
