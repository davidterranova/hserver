// Created by davidterranova on 16/10/2017.

package hserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Server struct {
	address      string
	port         int
	server       *http.Server
	mainRouter   *mux.Router
	negroni      *negroni.Negroni
	serveMux     *http.ServeMux
	creationTime time.Time
}

type ServerOption func(s *Server) error

func Port(port int) ServerOption {
	return func(s *Server) error {
		s.port = port
		return nil
	}
}

func Address(address string) ServerOption {
	return func(s *Server) error {
		s.address = address
		return nil
	}
}

func New(options ...ServerOption) *Server {
	s := Server{
		address:      "127.0.0.1",
		port:         8001,
		mainRouter:   mux.NewRouter().StrictSlash(false),
		negroni:      negroni.New(negroni.NewRecovery(), httpLoggerMiddleware()),
		serveMux:     http.NewServeMux(),
		creationTime: time.Now(),
	}

	for _, opt := range options {
		opt(&s)
	}
	s.server = &http.Server{
		Addr: s.listeningAddress(),
		//Handler: HttpLogger(s.negroni),
		Handler: s.negroni,
		//Handler:        HttpLogger(s.mainRouter),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.GetRouter("/").HandleFunc(
		"/health_check",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			status := struct {
				Status     string    `json:"status"`
				ServerTime time.Time `json:"server_time"`
				UpTime     string    `json:"up_time"`
			}{
				Status:     "OK",
				ServerTime: time.Now(),
				UpTime:     fmt.Sprintf("%ds", int64(time.Now().Sub(s.creationTime).Seconds())),
			}
			bytes, err := json.Marshal(status)
			if err != nil {
				logrus.Error(err)
			}
			w.Write(bytes)
		},
	)

	s.serveMux.Handle("/", s.mainRouter)
	return &s
}

func (s *Server) Listen() {
	logrus.WithFields(logrus.Fields{"address": s.address, "port": s.port}).Info("start server")
	s.negroni.UseHandler(s.serveMux)
	go s.server.ListenAndServe()
}

func (s *Server) GetRouter(pathPrefix string) *mux.Router {
	if pathPrefix == "" || pathPrefix == "/" {
		return s.mainRouter
	} else {
		return s.mainRouter.PathPrefix(pathPrefix).Subrouter()
	}
}

func (s *Server) AddMiddleware(pathPrefix string, middleware negroni.HandlerFunc) {
	s.serveMux.Handle(
		pathPrefix,
		negroni.New(
			middleware,
			negroni.Wrap(s.mainRouter),
		),
	)
}

func (s *Server) Shutdown(ctx context.Context) {
	logrus.Println("Shutting down server...")
	s.server.Shutdown(ctx)
}

func (s *Server) listeningAddress() string {
	return fmt.Sprintf("%s:%d", s.address, s.port)
}
