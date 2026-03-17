package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	service := NewResolverService(cfg)
	handler := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/resolve", handler.Resolve)

	loggedMux := LoggingMiddleware(logger)(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: loggedMux,
	}

	logger.Info("starting server", "port", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
