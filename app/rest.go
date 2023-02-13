package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alesr/urltinyizer/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RESTApp is an app that implements the App interface.
type RESTApp struct {
	logger  *zap.Logger
	server  *http.Server
	service service.Service
}

// NewREST creates a new REST app.
func NewREST(logger *zap.Logger, router *chi.Mux, service service.Service) *RESTApp {
	return &RESTApp{
		logger: logger,
		server: &http.Server{
			ReadTimeout:       time.Duration(5) * time.Second,
			ReadHeaderTimeout: time.Duration(5) * time.Second,
			WriteTimeout:      time.Duration(10) * time.Second,
			Addr:              ":8080",
			Handler:           router,
		},
		service: service,
	}
}

func (app *RESTApp) RegisterRoutes() {
	app.server.Handler.(*chi.Mux).Post("/shorten", app.createShortURL())
	app.server.Handler.(*chi.Mux).Get("/{shortURL}", app.redirectToLongURL())
	app.server.Handler.(*chi.Mux).Get("/{shortURL}/stats", app.getStats())
}

// Run starts the REST API server and listens for cancellation signals.
func (app *RESTApp) Run(ctx context.Context) error {
	app.logger.Info("starting server on port 8080")

	go func() {
		<-ctx.Done()
		if err := app.server.Shutdown(ctx); err != nil {
			app.logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	if err := app.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("could not start server: %w", err)
	}
	return nil
}

func (app *RESTApp) Terminate(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

// CreateShortURL creates a new short URL.
func (app *RESTApp) createShortURL() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var reqPayload CreateShortURLRequest
		if err := json.NewDecoder(req.Body).Decode(&reqPayload); err != nil {
			app.logger.Error("could not decode request body", zap.Error(err))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := reqPayload.Validate(); err != nil {
			app.logger.Error("invalid request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		short, err := app.service.CreateShortURL(req.Context(), reqPayload.LongURL)
		if err != nil {
			app.logger.Error("could not create short URL", zap.Error(err))
			http.Error(w, "could not create short URL", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(CreateShortURLResponse{ShortURL: short}); err != nil {
			app.logger.Error("could not encode response", zap.Error(err))
			http.Error(w, "could not encode response", http.StatusInternalServerError)
			return
		}
	}
}

// RedirectToLongURL redirects to the long URL.
func (app *RESTApp) redirectToLongURL() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		escapedShortURL, err := url.PathUnescape(chi.URLParam(req, "shortURL"))
		if err != nil {
			app.logger.Error("could not unescape short URL", zap.Error(err))
			http.Error(w, "could not unescape short URL", http.StatusInternalServerError)
			return
		}

		shortURL := RedirectToLongURLRequest(escapedShortURL)

		if err := shortURL.Validate(); err != nil {
			app.logger.Error("invalid request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		longURL, err := app.service.RedirectToLongURL(req.Context(), string(shortURL))
		if err != nil {
			app.logger.Error("could not redirect to long URL", zap.Error(err))
			http.Error(w, "could not redirect to long URL", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, longURL, http.StatusFound)
	}
}

// GetStats returns the stats of a short URL.
func (app *RESTApp) getStats() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		escapedShortURL, err := url.PathUnescape(chi.URLParam(req, "shortURL"))
		if err != nil {
			app.logger.Error("could not unescape short URL", zap.Error(err))
			http.Error(w, "could not unescape short URL", http.StatusInternalServerError)
			return
		}

		shortURL := GetStatsRequest(escapedShortURL)

		app.logger.Info("getting stats", zap.String("shortURL", string(shortURL)))
		if err := shortURL.Validate(); err != nil {
			app.logger.Error("invalid request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stats, err := app.service.GetStats(req.Context(), string(shortURL))
		if err != nil {
			app.logger.Error("could not get stats", zap.Error(err))
			http.Error(w, "could not get stats", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		statsResp := GetStatsResponse{
			ShortURL: string(shortURL),
			Hits:     stats,
		}

		if err := json.NewEncoder(w).Encode(statsResp); err != nil {
			app.logger.Error("could not encode response", zap.Error(err))
			http.Error(w, "could not encode response", http.StatusInternalServerError)
			return
		}
	}
}
