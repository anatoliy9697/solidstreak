package web

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
)

type Server struct {
	Addr string
	Res  common.Resources
	s    *http.Server
}

func (s Server) Run(mainCtx context.Context, doneCh chan struct{}) {
	defer func() { doneCh <- struct{}{} }()

	s.Res.Logger.Info("web server initialization...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		staticPath := filepath.Join("static", r.URL.Path)
		if info, err := os.Stat(staticPath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, staticPath)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	s.s = &http.Server{
		Addr:    s.Addr,
		Handler: mux,
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
