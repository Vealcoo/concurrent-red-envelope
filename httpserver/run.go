package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test/config"
)

type server struct {
	httpServer *http.Server
}

func New() *server {
	return &server{}
}

func (s *server) Run(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Setting.GetInt("service.httpPort")),
		Handler: handler,
	}
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s *server) Close() {
	log.Println("HttpServer closing...")
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}
