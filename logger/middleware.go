package logger

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const key = "logger"
const traceparent = "traceparent"

func Middleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newLogger := logger.With(zap.String(traceparent, r.Header.Get(traceparent)))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger", newLogger)))
		})
	}
}

func L(ctx context.Context) *zap.Logger {
	v := ctx.Value(key)
	if v == nil {
		return zap.NewExample()
	}

	l, ok := v.(*zap.Logger)
	if ok {
		return l
	}
	return zap.NewExample()
}
