package categories_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tradeservice/internal/models"
	"tradeservice/internal/server/handler/categories"

	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"log/slog"
	mockcategories "tradeservice/internal/server/handler/categories/mockCategories"
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

func TestCategoriesController_GetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoriesList := []models.CategoryDto{
		{ID: "1", Name: "cat1"},
		{ID: "2", Name: "cat2"},
	}

	mockManager.EXPECT().GetCategory(gomock.Any()).Return(categoriesList, nil)

	c, rec := createContext(http.MethodGet, "/categories", nil)

	err := handler.GetCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_GetCategory_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	mockManager.EXPECT().GetCategory(gomock.Any()).Return(nil, models.ErrDB)

	c, rec := createContext(http.MethodGet, "/categories", nil)

	err := handler.GetCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCategoriesController_AddCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryName := "newcat"
	productID := "prod123"
	newID := "42"

	mockManager.EXPECT().AddCategory(gomock.Any(), categoryName, productID).Return(newID, nil)

	c, rec := createContext(http.MethodPost, "/categories/:categoryName/:productId", map[string]string{
		"categoryName": categoryName,
		"productId":    productID,
	})

	err := handler.AddCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, newID, strings.TrimSpace(rec.Body.String()))
}

func TestCategoriesController_AddCategory_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryName := "dupCat"
	productID := "prod123"

	mockManager.EXPECT().AddCategory(gomock.Any(), categoryName, productID).Return("", models.ErrUnique)

	c, rec := createContext(http.MethodPost, "/categories/:categoryName/:productId", map[string]string{
		"categoryName": categoryName,
		"productId":    productID,
	})

	err := handler.AddCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategoriesController_DeleteCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"

	mockManager.EXPECT().DeleteCategory(gomock.Any(), categoryID).Return(nil)

	c, rec := createContext(http.MethodDelete, "/categories/:categoryId", map[string]string{
		"categoryId": categoryID,
	})

	err := handler.DeleteCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_DeleteCategory_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"

	mockManager.EXPECT().DeleteCategory(gomock.Any(), categoryID).Return(models.ErrNotFound)

	c, rec := createContext(http.MethodDelete, "/categories/:categoryId", map[string]string{
		"categoryId": categoryID,
	})

	err := handler.DeleteCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCategoriesController_SetCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"
	categoryName := "updated"

	mockManager.EXPECT().SetCategory(gomock.Any(), categoryID, categoryName).Return(nil)

	c, rec := createContext(http.MethodPost, "/categories/:categoryId/:categoryName", map[string]string{
		"categoryId":   categoryID,
		"categoryName": categoryName,
	})

	err := handler.SetCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_SetCategory_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryID := "42"
	categoryName := "updated"

	mockManager.EXPECT().SetCategory(gomock.Any(), categoryID, categoryName).Return(models.ErrNotFound)

	c, rec := createContext(http.MethodPost, "/categories/:categoryId/:categoryName", map[string]string{
		"categoryId":   categoryID,
		"categoryName": categoryName,
	})

	err := handler.SetCategory(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
