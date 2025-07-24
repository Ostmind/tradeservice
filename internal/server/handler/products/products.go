package products

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"tradeservice/internal/models"

	"github.com/labstack/echo/v4"
)

//go:generate mockgen -source=products.go -destination=mockProducts/productsrepository.go

type ProductManager interface {
	AddProduct(ctx context.Context, name string) (id string, err error)
	GetProduct(ctx context.Context) ([]models.ProductDto, error)
	SetProduct(ctx context.Context, id string, name string) error
	DeleteProduct(ctx context.Context, id string) error
}

type ProductController struct {
	manager ProductManager
	logger  *slog.Logger
}

func NewProductHandler(manager ProductManager, log *slog.Logger) *ProductController {
	return &ProductController{manager, log}
}

func (ctr ProductController) GetProduct(echo echo.Context) error {
	ctr.logger.Debug("Get Request for Products")

	res, err := ctr.manager.GetProduct(echo.Request().Context())
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr ProductController) AddProduct(echo echo.Context) error {
	ctr.logger.Debug("Post Request for Products")

	productName := echo.Param("productName")

	res, err := ctr.manager.AddProduct(echo.Request().Context(), productName)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr ProductController) DeleteProduct(echo echo.Context) error {
	ctr.logger.Debug("Delete Request for Products")

	productID := echo.Param("productId")

	err := ctr.manager.DeleteProduct(echo.Request().Context(), productID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}

func (ctr ProductController) SetProduct(echo echo.Context) error {
	ctr.logger.Debug("Patch Request for Products")

	productID := echo.Param("productId")

	productName := echo.Param("productName")

	err := ctr.manager.SetProduct(echo.Request().Context(), productID, productName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}
