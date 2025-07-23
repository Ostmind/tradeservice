package categories_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tradeservice/internal/server/handler/categories"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"tradeservice/internal/models"

	"log/slog"
	mockcategories "tradeservice/internal/server/handler/categories/mockCategories"
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

func createContext(t *testing.T, method, path string, params map[string]string, body string) (echo.Context, *httptest.ResponseRecorder) {
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

func TestCategoriesController_GetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoriesList := []models.CategoryDto{
		{Id: "1", Name: "cat1"},
		{Id: "2", Name: "cat2"},
	}

	mockManager.EXPECT().GetCategory(gomock.Any()).Return(categoriesList, nil)

	c, rec := createContext(t, http.MethodGet, "/categories", nil, "")

	err := handler.GetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_GetCategory_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	mockManager.EXPECT().GetCategory(gomock.Any()).Return(nil, errors.New("db error"))

	c, rec := createContext(t, http.MethodGet, "/categories", nil, "")

	err := handler.GetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCategoriesController_AddCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryName := "newcat"
	productId := "prod123"
	newID := "42"

	mockManager.EXPECT().AddCategory(gomock.Any(), categoryName, productId).Return(newID, nil)

	c, rec := createContext(t, http.MethodPost, "/categories/:categoryName/:productId", map[string]string{
		"categoryName": categoryName,
		"productId":    productId,
	}, "")

	err := handler.AddCategory(c)
	assert.NoError(t, err)
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
	productId := "prod123"

	mockManager.EXPECT().AddCategory(gomock.Any(), categoryName, productId).Return("", models.ErrUnique)

	c, rec := createContext(t, http.MethodPost, "/categories/:categoryName/:productId", map[string]string{
		"categoryName": categoryName,
		"productId":    productId,
	}, "")

	err := handler.AddCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategoriesController_DeleteCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryId := "42"

	mockManager.EXPECT().DeleteCategory(gomock.Any(), categoryId).Return(nil)

	c, rec := createContext(t, http.MethodDelete, "/categories/:categoryId", map[string]string{
		"categoryId": categoryId,
	}, "")

	err := handler.DeleteCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_DeleteCategory_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryId := "42"

	mockManager.EXPECT().DeleteCategory(gomock.Any(), categoryId).Return(models.ErrNotFound)

	c, rec := createContext(t, http.MethodDelete, "/categories/:categoryId", map[string]string{
		"categoryId": categoryId,
	}, "")

	err := handler.DeleteCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCategoriesController_SetCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryId := "42"
	categoryName := "updated"

	mockManager.EXPECT().SetCategory(gomock.Any(), categoryId, categoryName).Return(nil)

	c, rec := createContext(t, http.MethodPost, "/categories/:categoryId/:categoryName", map[string]string{
		"categoryId":   categoryId,
		"categoryName": categoryName,
	}, "")

	err := handler.SetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCategoriesController_SetCategory_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mockcategories.NewMockCategoryManager(ctrl)
	logger := newTestLogger()
	handler := categories.NewCategoriesHandler(mockManager, logger)

	categoryId := "42"
	categoryName := "updated"

	mockManager.EXPECT().SetCategory(gomock.Any(), categoryId, categoryName).Return(models.ErrNotFound)

	c, rec := createContext(t, http.MethodPost, "/categories/:categoryId/:categoryName", map[string]string{
		"categoryId":   categoryId,
		"categoryName": categoryName,
	}, "")

	err := handler.SetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
