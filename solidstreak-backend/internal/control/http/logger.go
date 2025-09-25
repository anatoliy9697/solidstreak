package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

// responseLogger оборачивает http.ResponseWriter для логирования ответа
type responseLogger struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rl *responseLogger) WriteHeader(statusCode int) {
	rl.status = statusCode
	rl.ResponseWriter.WriteHeader(statusCode)
}

func (rl *responseLogger) Write(b []byte) (int, error) {
	rl.body.Write(b)
	return rl.ResponseWriter.Write(b)
}

func (s Server) Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := s.Res.Logger

			// Adding request ID to request context
			reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
			if reqID != "" {
				logger = logger.With("requestId", reqID)
			}

			// Request body reading
			var body string
			if r.Body != nil && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
				data, err := io.ReadAll(r.Body)
				if err == nil {
					body = string(data)
					r.Body = io.NopCloser(bytes.NewReader(data))
				}
			}

			// Request headers serialization
			reqHeaders := make(map[string]string)
			for k, v := range r.Header {
				if len(v) > 0 {
					reqHeaders[k] = v[0]
				}
			}

			logger.Info("request received",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Any("headers", reqHeaders),
				slog.String("body", body),
			)

			rl := &responseLogger{ResponseWriter: w, status: 200}
			next.ServeHTTP(rl, r)

			// Response headers serialization
			respHeaders := make(map[string]string)
			for k, v := range rl.Header() {
				if len(v) > 0 {
					respHeaders[k] = v[0]
				}
			}

			logger.Info("response sent",
				slog.Int("status", rl.status),
				slog.Any("headers", respHeaders),
				slog.String("body", rl.body.String()),
			)
		})
	}
}
