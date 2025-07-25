package categories_test

import (
	"net/http"
	"strings"
	"testing"
	"tradeservice/internal/models"
	"tradeservice/internal/server/handler/categories"
	"tradeservice/internal/server/utils"

	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	mockcategories "tradeservice/internal/server/handler/categories/mockCategories"
)

func TestCategoriesController_GetCategory(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoriesList := []models.CategoryDto{
		{ID: "1", Name: "cat1"},
		{ID: "2", Name: "cat2"},
	}

	mockManager.EXPECT().GetCategory(gomock.Any()).Return(categoriesList, nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodGet, "/categories", nil)

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.GetCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_GetCategory_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	mockManager.EXPECT().GetCategory(gomock.Any()).Return(nil, models.ErrDB)

	rec, req, keys, vals := utils.CreateContext(http.MethodGet, "/categories", nil)

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.GetCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCategoriesController_AddCategory_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryName := "newcat"
	productID := "prod123"
	newID := "42"

	mockManager.EXPECT().AddCategory(gomock.Any(), categoryName, productID).Return(newID, nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/categories/:categoryName/:productId", map[string]string{
		"categoryName": categoryName,
		"productID":    productID,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.AddCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, newID, strings.TrimSpace(rec.Body.String()))
}

func TestCategoriesController_AddCategory_Conflict(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryName := "dupCat"
	productID := "prod123"

	mockManager.EXPECT().AddCategory(gomock.Any(), categoryName, productID).Return("", models.ErrUnique)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/categories/:categoryName/:productId", map[string]string{
		"categoryName": categoryName,
		"productID":    productID,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.AddCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategoriesController_DeleteCategory_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"

	mockManager.EXPECT().DeleteCategory(gomock.Any(), categoryID).Return(nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodDelete, "/categories/:categoryId", map[string]string{
		"categoryId": categoryID,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.DeleteCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_DeleteCategory_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"

	mockManager.EXPECT().DeleteCategory(gomock.Any(), categoryID).Return(models.ErrNotFound)

	rec, req, keys, vals := utils.CreateContext(http.MethodDelete, "/categories/:categoryId", map[string]string{
		"categoryId": categoryID,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.DeleteCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCategoriesController_SetCategory_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"
	categoryName := "updated"

	mockManager.EXPECT().SetCategory(gomock.Any(), categoryID, categoryName).Return(nil)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/categories/:categoryId/:categoryName", map[string]string{
		"categoryId":   categoryID,
		"categoryName": categoryName,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.SetCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_SetCategory_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := utils.NewTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"
	categoryName := "updated"

	mockManager.EXPECT().SetCategory(gomock.Any(), categoryID, categoryName).Return(models.ErrNotFound)

	rec, req, keys, vals := utils.CreateContext(http.MethodPost, "/categories/:categoryId/:categoryName", map[string]string{
		"categoryId":   categoryID,
		"categoryName": categoryName,
	})

	e := echo.New()
	echoCtx := e.NewContext(req, rec)
	echoCtx.SetParamNames(keys...)
	echoCtx.SetParamValues(vals...)

	err := handler.SetCategory(echoCtx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
