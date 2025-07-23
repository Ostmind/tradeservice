package products_test

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tradeservice/internal/models"
	"tradeservice/internal/server/handler/products"
	mockproducts "tradeservice/internal/server/handler/products/mockProducts"
)

// простой тестовый логгер, для передачи в контроллер
func newTestLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func createContext(method, path string, params map[string]string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	keys := make([]string, 0, len(params))
	vals := make([]string, 0, len(params))
	for k, v := range params {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	c.SetParamNames(keys...)
	c.SetParamValues(vals...)

	return c, rec
}

func TestCategoriesController_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productList := []models.ProductDto{
		{Id: "1", Name: "cat1"},
		{Id: "2", Name: "cat2"},
	}

	mockManager.EXPECT().GetProduct(gomock.Any()).Return(productList, nil)

	c, rec := createContext(http.MethodGet, "/products", nil, "")

	err := handler.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_GetProduct_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	mockManager.EXPECT().GetProduct(gomock.Any()).Return(nil, errors.New("db error"))

	c, rec := createContext(http.MethodGet, "/products", nil, "")

	err := handler.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCategoriesController_AddProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productName := "newcat"
	newID := "42"

	mockManager.EXPECT().AddProduct(gomock.Any(), productName).Return(newID, nil)

	c, rec := createContext(http.MethodPost, "/products/create/:productName", map[string]string{
		"categoryName": productName,
	}, "")

	err := handler.AddProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, newID, strings.TrimSpace(rec.Body.String()))
}

func TestCategoriesController_AddProduct_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productName := "dupProd"

	mockManager.EXPECT().AddProduct(gomock.Any(), productName).Return("", models.ErrUnique)

	c, rec := createContext(http.MethodPost, "/products/create/:productName", map[string]string{
		"productName": productName,
	}, "")

	err := handler.AddProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategoriesController_DeleteProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productId := "42"

	mockManager.EXPECT().DeleteProduct(gomock.Any(), productId).Return(nil)

	c, rec := createContext(http.MethodDelete, "/product/:productId", map[string]string{
		"productId": productId,
	}, "")

	err := handler.DeleteProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_DeleteProduct_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productId := "42"

	mockManager.EXPECT().DeleteProduct(gomock.Any(), productId).Return(models.ErrNotFound)

	c, rec := createContext(http.MethodDelete, "/product/:productId", map[string]string{
		"productId": productId,
	}, "")

	err := handler.DeleteProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCategoriesController_SetProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productId := "42"
	productName := "updated"

	mockManager.EXPECT().SetProduct(gomock.Any(), productId, productName).Return(nil)

	c, rec := createContext(http.MethodPost, "/update/:productName/:productId", map[string]string{
		"productId":   productId,
		"productName": productName,
	}, "")

	err := handler.SetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_SetProduct_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productId := "42"
	productName := "updated"

	mockManager.EXPECT().SetProduct(gomock.Any(), productId, productName).Return(models.ErrNotFound)

	c, rec := createContext(http.MethodPost, "/update/:productName/:productId", map[string]string{
		"productId":   productId,
		"productName": productName,
	}, "")

	err := handler.SetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
