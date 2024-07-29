package app

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type HttpServer struct {
	server *http.Server
}

func NewServer(cfg *config.Config, router chi.Router) *HttpServer {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	return &HttpServer{
		server: &http.Server{
			Addr:         ":" + cfg.HTTPServer.Port,
			Handler:      router,
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
			TLSConfig:    tlsConfig,
		},
	}
}

func (hs *HttpServer) Run(log *zap.Logger, cfg *config.Config, certFile string, keyFile string) error {
	log.Info("starting server: ", zap.String("port", cfg.HTTPServer.Port))

	if certFile == "" || keyFile == "" {
		log.Fatal("Certificate or key file path not set in environment variables: ", zap.String("certFile", certFile), zap.String("keyFile", keyFile))
		return errors.New("certificate or key file path not set in environment variables")
	}

	go func() {
		err := hs.server.ListenAndServeTLS(certFile, keyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("error occurred on server shutting down: ", zap.String("error", err.Error()))
			return
		}
	}()

	log.Info("server started")

	return nil
}

func (hs *HttpServer) Shutdown(ctx context.Context, log *zap.Logger, stopTimeout time.Duration) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalChan:
		log.Info("received shutdown signal")

		shutdownCtx, cancel := context.WithTimeout(ctx, stopTimeout)
		defer cancel()

		if err := hs.server.Shutdown(shutdownCtx); err != nil {
			log.Error("error occurred on server shutting down", zap.String("error", err.Error()))
		} else {
			log.Info("server stopped")
		}
	case <-ctx.Done():
		log.Info("context done, skipping shutdown")
	}
}
