package products_test

import (
	"net/http"
	"strings"
	"testing"
	"tradeservice/internal/models"
	"tradeservice/internal/server/handler/products"
	mockproducts "tradeservice/internal/server/handler/products/mockProducts"
	"tradeservice/internal/server/utils"

	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCategoriesController_GetProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productList := []models.ProductDto{
		{ID: "1", Name: "cat1"},
		{ID: "2", Name: "cat2"},
	}

	mockManager.EXPECT().GetProduct(gomock.Any()).Return(productList, nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodGet, "/products", nil)

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.GetProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_GetProduct_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	mockManager.EXPECT().GetProduct(gomock.Any()).Return(nil, models.ErrDB)

	rec, req, keys, vals := utils.CreateContext(http.MethodGet, "/products", nil)

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.GetProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCategoriesController_AddProduct_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productName := "newcat"
	newID := "42"

	mockManager.EXPECT().AddProduct(gomock.Any(), productName).Return(newID, nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/products/create/:productName", map[string]string{
		"categoryName": productName,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.AddProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, newID, strings.TrimSpace(rec.Body.String()))
}

func TestCategoriesController_AddProduct_Conflict(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productName := "dupProd"

	mockManager.EXPECT().AddProduct(gomock.Any(), productName).Return("", models.ErrUnique)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/products/create/:productName", map[string]string{
		"productName": productName,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.AddProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategoriesController_DeleteProduct_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"

	mockManager.EXPECT().DeleteProduct(gomock.Any(), productID).Return(nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodDelete, "/product/:productId", map[string]string{
		"productId": productID,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.DeleteProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_DeleteProduct_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"

	mockManager.EXPECT().DeleteProduct(gomock.Any(), productID).Return(models.ErrNotFound)

	rec, req, keys, vals := utils.CreateContext(http.MethodDelete, "/product/:productId", map[string]string{
		"productId": productID,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.DeleteProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCategoriesController_SetProduct_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"
	productName := "updated"

	mockManager.EXPECT().SetProduct(gomock.Any(), productID, productName).Return(nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/update/:productName/:productId", map[string]string{
		"productId":   productID,
		"productName": productName,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.SetProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_SetProduct_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockproducts.NewMockProductManager(ctrl)
	logger := utils.NewTestLogger()
	handler := products.NewProductHandler(mockManager, logger)

	productID := "42"
	productName := "updated"

	mockManager.EXPECT().SetProduct(gomock.Any(), productID, productName).Return(models.ErrNotFound)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/update/:productName/:productId", map[string]string{
		"productId":   productID,
		"productName": productName,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.SetProduct(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
