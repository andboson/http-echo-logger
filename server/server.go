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
	address          string
	server           *http.Server
	history          History
	tpls             *templates.Templates
	echoEndpoints    []endpoint
	echoEndpointsMap map[string][]endpoint
}

type endpoint struct {
	Path            string   `json:"path"`
	ResponseHeaders []string `json:"headers"`
	MatchRequest    string   `json:"request"`
	MockResponse    string   `json:"mock"`
}

// NewServer returns instance of a service and sets up a server
func NewServer(addr string, tpls *templates.Templates, echoEndpointsCustom string) Server {
	mux := http.NewServeMux()
	endpoints, err := getEndpoints(echoEndpointsCustom)
	if err != nil {
		log.Fatalf("error getting enpoints:%+v", err)
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

	s.makeMockEndpoints(endpoints)

	mux.Handle(indexEndpoint, s.createHTTPHandler())
	mux.Handle(apiEndpoint, s.createAPIHandler())

	return s
}

func getEndpoints(endpointsArray string) ([]endpoint, error) {
	if endpointsArray == "" {
		return nil, nil
	}

	var result []endpoint
	return result, json.Unmarshal([]byte(endpointsArray), &result)
}

// Start starts a httpserver
func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Wrap(err, "can't create listener")
	}
	log.Printf("HTTP CLI LOGGER server started: %s", DefaultHTTPAddr)

	for _, echoEndpoint := range s.echoEndpoints {
		log.Printf("endpoint: %+v", echoEndpoint)
	}
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
		if r.RequestURI != "/" {
			if mocks, ok := s.echoEndpointsMap[r.RequestURI]; ok {
				s.createEchoHandler(mocks, w, r)
				return
			}

			s.processEchoMock(w, r, endpoint{})
			return
		}

		if s.tpls == nil {
			return
		}
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

func (s *server) makeMockEndpoints(endpoints []endpoint) {
	endpointsMap := map[string][]endpoint{}
	for _, e := range endpoints {
		endpointsMap[e.Path] = append(endpointsMap[e.Path], e)
	}

	s.echoEndpointsMap = endpointsMap
}

func (s *server) createEchoHandler(mocks []endpoint, w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body:%s", err)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	for _, mock := range mocks {
		if strings.Contains(string(reqBody), mock.MatchRequest) || mock.MatchRequest == "" {
			s.processEchoMock(w, r, mock)
			return
		}
	}

	fmt.Printf("\n == default enpoint ==")
	s.processEchoMock(w, r, endpoint{})
}

func (s *server) processEchoMock(w http.ResponseWriter, r *http.Request, mock endpoint) {
	item := HistoryItem{
		Request: r,
		Date:    time.Now(),
	}

	mockResponse := mock.MockResponse

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
}
