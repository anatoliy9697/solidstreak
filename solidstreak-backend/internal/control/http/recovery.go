package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s Server) Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Adding request ID to request context
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			r = r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID{}, reqID))

			// Updating logger with request ID
			logger := s.Res.Logger
			if reqID != "" {
				logger = logger.With("requestId", reqID)
			}

			// Recovering from panics
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic recovered", "panic", r)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			// Calling the next handler
			next.ServeHTTP(w, r)
		})
	}
}
