package http

import (
	"context"
	"net/http"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/errors"
	"github.com/google/uuid"
)

func (s Server) Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

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
					writeError(w, apperrors.ErrInternal(r.(string)))
				}
			}()

			// Calling the next handler
			next.ServeHTTP(w, r)
		})
	}
}
