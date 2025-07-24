package products_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tradeservice/internal/models"
	"tradeservice/internal/server/handler/products"
	mockproducts "tradeservice/internal/server/handler/products/mockProducts"

	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

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

func createContext(method, path string, params map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(""))
	rec := httptest.NewRecorder()
	echo := e.NewContext(req, rec)

	keys := make([]string, 0, len(params))
	vals := make([]string, 0, len(params))

	for k, v := range params {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	echo.SetParamNames(keys...)

	echo.SetParamValues(vals...)

	return echo, rec
}

func TestCategoriesController_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productList := []models.ProductDto{
		{ID: "1", Name: "cat1"},
		{ID: "2", Name: "cat2"},
	}

	mockManager.EXPECT().GetProduct(gomock.Any()).Return(productList, nil)

	echo, rec := createContext(http.MethodGet, "/products", nil)

	err := handler.GetProduct(echo)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_GetProduct_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	mockManager.EXPECT().GetProduct(gomock.Any()).Return(nil, models.ErrDB)

	echo, rec := createContext(http.MethodGet, "/products", nil)

	err := handler.GetProduct(echo)
	require.NoError(t, err)
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

	echo, rec := createContext(http.MethodPost, "/products/create/:productName", map[string]string{
		"categoryName": productName,
	})

	err := handler.AddProduct(echo)
	require.NoError(t, err)
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

	echo, rec := createContext(http.MethodPost, "/products/create/:productName", map[string]string{
		"productName": productName,
	})

	err := handler.AddProduct(echo)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategoriesController_DeleteProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"

	mockManager.EXPECT().DeleteProduct(gomock.Any(), productID).Return(nil)

	echo, rec := createContext(http.MethodDelete, "/product/:productId", map[string]string{
		"productId": productID,
	})

	err := handler.DeleteProduct(echo)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_DeleteProduct_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"

	mockManager.EXPECT().DeleteProduct(gomock.Any(), productID).Return(models.ErrNotFound)

	echo, rec := createContext(http.MethodDelete, "/product/:productId", map[string]string{
		"productId": productID,
	})

	err := handler.DeleteProduct(echo)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCategoriesController_SetProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"
	productName := "updated"

	mockManager.EXPECT().SetProduct(gomock.Any(), productID, productName).Return(nil)

	echo, rec := createContext(http.MethodPost, "/update/:productName/:productId", map[string]string{
		"productId":   productID,
		"productName": productName,
	})

	err := handler.SetProduct(echo)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_SetProduct_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := newTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"
	productName := "updated"

	mockManager.EXPECT().SetProduct(gomock.Any(), productID, productName).Return(models.ErrNotFound)

	echo, rec := createContext(http.MethodPost, "/update/:productName/:productId", map[string]string{
		"productId":   productID,
		"productName": productName,
	})

	err := handler.SetProduct(echo)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
