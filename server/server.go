package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"andboson/http-echo-logger/templates"

	"github.com/pkg/errors"
)

//Server holds a HTTP server methods
type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type server struct {
	address       string
	server        *http.Server
	history       History
	tpls          *templates.Templates
	echoEndpoints []string
}

// NewServer returns instance of a service and sets up a server
func NewServer(addr string, tpls *templates.Templates, echoEndpointsCustom string) Server {
	mux := http.NewServeMux()

	endpoints := strings.Split(echoEndpointsCustom, "\n")
	if !strings.Contains(echoEndpointsCustom, "\n") {
		endpoints = strings.Split(echoEndpointsCustom, " ")
	}

	s := &server{
		tpls:          tpls,
		history:       History{},
		address:       addr,
		echoEndpoints: endpoints,
		server: &http.Server{
			Handler: mux,
		},
	}

	mux.Handle(indexEndpoint, s.createHTTPHandler())

	if echoEndpointsCustom == "" {
		s.echoEndpoints = append(s.echoEndpoints, echoEndpointDedfault)
	}

	for i, endpoint := range s.echoEndpoints {
		s.echoEndpoints[i] = strings.Trim(strings.TrimSpace(endpoint), "\n")
		if s.echoEndpoints[i] == "" {
			s.echoEndpoints = append(s.echoEndpoints[:i], s.echoEndpoints[i+1:]...)
			continue
		}
		log.Printf("=>%s<", s.echoEndpoints[i])
		mux.Handle(s.echoEndpoints[i], s.createEchoHandler())
	}

	return s
}

// Start starts a httpserver
func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Wrap(err, "can't create listener")
	}
	log.Printf("HTTP CLI LOGGER server started: %s, echo endpoints: %s",
		DefaultHTTPAddr,
		strings.Join(s.echoEndpoints, ","))
	if err := s.server.Serve(ln); err != nil {
		return errors.Wrap(err, "can't start server")
	}

	return nil
}

func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) createHTTPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.tpls.Tpls.Execute(w, s.history.GetHistory(true)); err != nil {
			fmt.Fprintf(w, "%+v", err)
		}
	})
}

func (s *server) createEchoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item := HistoryItem{
			Request: r,
			Date:    time.Now(),
		}
		item.PrintConsole(w)
		s.history.AddItem(item)
		w.Header().Add("referrer", "http-echo-server")
	})
}
