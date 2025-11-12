package main

import (
	"comment-tree-service/intenal/config"
	delete2 "comment-tree-service/intenal/http-server/handlers/comment/delete"
	"comment-tree-service/intenal/http-server/handlers/comment/find"
	"comment-tree-service/intenal/http-server/handlers/comment/save"
	"comment-tree-service/intenal/service"
	"comment-tree-service/intenal/storage/mongodb"
	"comment-tree-service/lib/logger/logr"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := logr.InitLogrusLog(cfg.Env, os.Stdout)

	log.Info("INIT CONFIG")
	log.Debug("config ", *cfg)

	storage, err := mongodb.NewMongoStorage(cfg)
	if err != nil {
		log.Error(err)
		return
	}

	ser := service.NewCommentsService(storage)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/comments", find.NewGetHandler(log, ser))
	router.Get("/comments-search", find.NewSearchHandler(log, ser))
	router.Delete("/comments/{id}", delete2.New(log, ser))
	router.Post("/comments", save.New(log, ser))

	err = serve(log, cfg, router, storage)
	log.Info("stop", err)

}

func serve(log *logrus.Logger, cfg *config.Config, h http.Handler, storage *mongodb.MongoStorage) error {
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      h,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		log.Info("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Info("starting server wb-examples-l0",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.HTTPServer.Address),
	)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	storage.Close(context.Background())

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Info("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
