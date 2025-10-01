package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/errors"
)

type ctxKeyUserTgID struct{}

func (s Server) ValidateTelegramInitData() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error

			logger := s.Res.Logger

			// Adding request ID to request context
			reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
			if reqID != "" {
				logger = logger.With("requestId", reqID)
			}

			initData := r.Header.Get("X-Telegram-InitData")
			if initData == "" {
				processError(w, logger, apperrors.ErrUnauthorized("missing telegram initData"))
				return
			}

			var userTgID int64
			if userTgID, err = validateAndGetUserTgID(initData, s.Res.TgBotAPI.Token); err != nil {
				processError(w, logger, apperrors.ErrUnauthorized(err.Error()))
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyUserTgID{}, userTgID))

			next.ServeHTTP(w, r)
		})
	}
}

func sign(payload, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}

func validateAndGetUserTgID(initData, token string) (int64, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, err
	}

	var (
		hash     string
		userTgID int64
		pairs    = make([]string, 0, len(values))
	)

	for k, v := range values {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "user" {
			if userTgID, err = extractUserTgIDFromJSONData(v[0]); err != nil {
				return 0, err
			}
		}
		pairs = append(pairs, k+"="+v[0])
	}

	if hash == "" {
		return 0, errors.New("missing hash in telegram initData")
	}

	sort.Strings(pairs)

	if sign(strings.Join(pairs, "\n"), token) != hash {
		return 0, errors.New("request authentication failed")
	}

	return userTgID, nil
}

func extractUserTgIDFromJSONData(jsonStr string) (int64, error) {
	type userData struct {
		ID int64 `json:"id"`
	}

	var u userData
	if err := json.Unmarshal([]byte(jsonStr), &u); err != nil {
		return 0, err
	}

	return u.ID, nil
}
