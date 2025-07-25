package utils

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
)

func NewTestLogger() *slog.Logger {
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

func CreateContext(method, path string, params map[string]string) (*httptest.ResponseRecorder, *http.Request, []string, []string) {
	req := httptest.NewRequest(method, path, strings.NewReader(""))
	rec := httptest.NewRecorder()

	keys := make([]string, 0, len(params))
	vals := make([]string, 0, len(params))

	for k, v := range params {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	return rec, req, keys, vals
}
