package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lroman242/redirector/config"
)

// ErrHTTPAddressIsNotSet describe error when which occur if host is not set in HTTP Server config.
var ErrHTTPAddressIsNotSet = errors.New("HTTP Address is not set")

// Server type describe http server instance.
type Server struct {
	config     *config.HTTPServerConf
	httpServer *http.Server
}

// NewServer create new instance of Server.
func NewServer(conf *config.HTTPServerConf, handler http.Handler) *Server {
	s := &Server{
		config: conf,
		httpServer: &http.Server{
			Addr:         conf.GetHTTPAddress(),
			ReadTimeout:  conf.GetHTTPReadTimeout(),
			WriteTimeout: conf.GetHTTPWriteTimeout(),
			Handler:      handler,
		},
	}

	return s
}

// Start function will start HTTP Server listener.
func (s *Server) Start() error {
	stop := s.subscribeForSignals()

	go func() {
		if err := s.startHTTPServer(); err != nil {
			panic(err)
		}
	}()

	<-stop
	slog.Info("Stop signal received..")

	s.stop()

	return nil
}

func (s *Server) subscribeForSignals() chan os.Signal {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, os.Kill)
	signal.Notify(stop, syscall.SIGTERM)

	return stop
}

func (s *Server) stop() {
	quit := make(chan struct{})

	go func() {
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()

			_ = s.httpServer.Shutdown(context.Background())
		}()

		//// TODO: close other services
		// wg.Add(1)
		// go func() {
		//	defer wg.Done()
		//
		//	// TODO:
		// }()

		wg.Wait()

		close(quit)
	}()

	select {
	case <-quit:
	case <-time.After(s.config.GetShutdownTimeout()):
		slog.Info("Close service by timeout")
	}
}

func (s *Server) startHTTPServer() error {
	if len(s.config.GetHTTPAddress()) == 0 {
		return ErrHTTPAddressIsNotSet
	}

	slog.Info("Start listening", "host", s.config.GetHTTPAddress())

	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server listener error: %w", err)
		}

		return fmt.Errorf("http listener error: %w", err)
	}

	return nil
}
