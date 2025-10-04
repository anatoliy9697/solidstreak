package http

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/resources"
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

	api.Post("/habits", s.postHabit)
	api.Get("/habits/{id}", s.getHabit)
	api.Put("/habits/{id}", s.putHabit)
	api.Get("/habits", s.getHabits)
	api.Post("/users/{userId}/habit/{habitId}/check", s.postUserHabitCheck)

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

func getUserIDFromURL(r *http.Request) (int64, error) {
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		return 0, apperrors.ErrBadRequest("missing user id")
	}

	var (
		userID int64
		err    error
	)
	if userID, err = strconv.ParseInt(userIDStr, 10, 64); err != nil {
		return 0, apperrors.ErrBadRequest("invalid user id")
	}

	return userID, nil
}

func getHabitIDFromURL(r *http.Request) (int64, error) {
	habitIDStr := chi.URLParam(r, "habitId")
	if habitIDStr == "" {
		return 0, apperrors.ErrBadRequest("missing habit id")
	}

	var (
		habitID int64
		err     error
	)
	if habitID, err = strconv.ParseInt(habitIDStr, 10, 64); err != nil {
		return 0, apperrors.ErrBadRequest("invalid habit id")
	}

	return habitID, nil
}
