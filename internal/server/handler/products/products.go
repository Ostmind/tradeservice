package products

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"tradeservice/internal/models"
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

func (ctr ProductController) GetProduct(c echo.Context) error {

	ctr.logger.Debug("Get Request for Products")

	res, err := ctr.manager.GetProduct(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr ProductController) AddProduct(c echo.Context) error {

	ctr.logger.Debug("Post Request for Products")

	productName := c.Param("productName")

	res, err := ctr.manager.AddProduct(c.Request().Context(), productName)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr ProductController) DeleteProduct(c echo.Context) error {

	ctr.logger.Debug("Delete Request for Products")

	productId := c.Param("productId")

	err := ctr.manager.DeleteProduct(c.Request().Context(), productId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (ctr ProductController) SetProduct(c echo.Context) error {

	ctr.logger.Debug("Patch Request for Products")

	productId := c.Param("productId")

	productName := c.Param("productName")

	err := ctr.manager.SetProduct(c.Request().Context(), productId, productName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
