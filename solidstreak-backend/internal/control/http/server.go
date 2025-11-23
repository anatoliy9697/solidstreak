package http

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/resources"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Addr string
	Res  resources.Resources
	s    *http.Server
}

func (s Server) Run(mainCtx context.Context, doneCh chan struct{}) {
	defer func() { doneCh <- struct{}{} }()

	s.Res.Logger.Info("web server initialization...")

	router := chi.NewRouter()

	api := chi.NewRouter()

	api.Use(s.Recovery())
	api.Use(s.Logger())
	api.Use(s.ValidateTelegramInitData())

	api.Post("/user-info/upsert", s.postUserInfo)
	api.Get("/users/{userId}", s.getUser)

	api.Post("/users/{userID}/habits", s.postHabit)
	api.Put("/users/{userID}/habits/{habitID}", s.putHabit)
	api.Get("/users/{userID}/habits/{habitID}", s.getHabit)
	api.Delete("/users/{userID}/habits/{habitID}", s.deleteHabit)
	api.Get("/users/{userID}/habits", s.getHabits)
	api.Post("/users/{userID}/habits/{habitID}/checks", s.postUserHabitCheck)
	api.Get("/users/{userID}/habits/{habitID}/checks", s.getUserHabitCompletedChecks)

	router.Mount("/api/v1", api)

	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		staticPath := filepath.Join("static", r.URL.Path)
		if info, err := os.Stat(staticPath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, staticPath)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	s.s = &http.Server{
		Addr:    s.Addr,
		Handler: router,
	}

	s.Res.Logger.Info("web server started on " + s.Addr)

	errCh := make(chan error, 1)
	go func() {
		err := s.s.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
		errCh <- err
	}()

	select {
	case <-mainCtx.Done():
		s.Res.Logger.Info("web server shutting down initiated")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.s.Shutdown(ctx); err != nil {
			s.Res.Logger.Error("web server shutdown error", "error", err)
		}
	case err := <-errCh:
		if err != nil {
			s.Res.Logger.Error("web server error", "error", err)
		}
	}

	s.Res.Logger.Info("web server stopped")
}

func getInt64FromURLParams(r *http.Request, key string, required bool) (int64, error) {
	strValue := chi.URLParam(r, key)
	if strValue == "" {
		if required {
			return 0, apperrors.ErrBadRequest("missing \"" + key + "\" in URL params")
		}
		return 0, nil
	}

	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return 0, apperrors.ErrBadRequest("invalid \"" + key + "\" in URL params")
	}

	return value, nil
}

// func getInt64FromURLQuery(r *http.Request, key string, required bool) (int64, error) {
// 	strValue := r.URL.Query().Get(key)

// 	if strValue == "" {
// 		if required {
// 			return 0, apperrors.ErrBadRequest("missing \"" + key + "\" in URL query")
// 		}
// 		return 0, nil
// 	}

// 	value, err := strconv.ParseInt(strValue, 10, 64)
// 	if err != nil {
// 		return 0, apperrors.ErrBadRequest("invalid \"" + key + "\" in URL query")
// 	}

// 	return value, nil
// }

func getDateFromURLQuery(r *http.Request, key string, required bool) (*date.Date, error) {
	dateStr := r.URL.Query().Get(key)

	if dateStr == "" {
		if required {
			return nil, apperrors.ErrBadRequest("missing \"" + key + "\" date in URL query")
		}
		return nil, nil
	}

	d, err := date.Parse(dateStr)
	if err != nil {
		return nil, apperrors.ErrBadRequest("invalid \"" + key + "\" date in URL query")
	}

	return &d, nil
}

func getBoolFromURLQuery(r *http.Request, key string, required bool) (bool, error) {
	strValue := r.URL.Query().Get(key)
	if strValue == "" {
		if required {
			return false, apperrors.ErrBadRequest("missing \"" + key + "\" in URL query")
		}
		return false, nil
	}

	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, apperrors.ErrBadRequest("invalid \"" + key + "\" in URL query")
	}

	return value, nil
}
