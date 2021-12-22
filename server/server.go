package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	endpoints := strings.Split(strings.TrimSpace(echoEndpointsCustom), "\n")
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
	mux.Handle(apiEndpoint, s.createAPIHandler())

	if echoEndpointsCustom == "" {
		s.echoEndpoints = append([]string{}, echoEndpointDefault)
	}

	for i, endpoint := range s.echoEndpoints {
		s.echoEndpoints[i] = strings.Trim(strings.TrimSpace(endpoint), "\n")
		if s.echoEndpoints[i] == "" {
			s.echoEndpoints = append(s.echoEndpoints[:i], s.echoEndpoints[i+1:]...)
			continue
		}

		mockResponse := ""
		if strings.Contains(s.echoEndpoints[i], endpointMockDelimiter) {
			splitted := strings.Split(s.echoEndpoints[i], endpointMockDelimiter)
			s.echoEndpoints[i] = splitted[0]
			mockResponse = splitted[1]
		}
		mux.Handle(s.echoEndpoints[i], s.createEchoHandler(mockResponse))
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

func (s *server) createAPIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.history.GetHistory(true)); err != nil {
			fmt.Fprintf(w, "%+v", err)
		}
	})
}

func (s *server) createEchoHandler(mockResponse string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item := HistoryItem{
			Request: r,
			Date:    time.Now(),
		}

		if mockResponse != "" {
			item = HistoryItem{
				Request: &http.Request{
					Header:        r.Header,
					Body:          ioutil.NopCloser(bytes.NewBufferString(mockResponse)),
					ContentLength: int64(len(mockResponse)),
					Method:        r.Method,
					URL:           r.URL,
					RemoteAddr:    r.RemoteAddr,
					RequestURI:    r.RequestURI,
				},
				Date: time.Now(),
			}
			item.Header.Set("Content-Length", fmt.Sprintf("%d", len(mockResponse)))
			item.Header.Set("Content-Type", "application/json")
		}

		item.PrintConsole(w)
		s.history.AddItem(item)
		w.Header().Add("referrer", "http-echo-server")
	})
}
