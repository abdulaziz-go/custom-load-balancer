package server

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	Url          *url.URL
	Alive        bool
	Weight       int
	Connections  int
	ReverseProxy *httputil.ReverseProxy
	mux          sync.RWMutex
}

func NewServerAdd(serverUrl string, weight int) (*Server, error) {
	urlstr, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	return &Server{
		Url:          urlstr,
		Alive:        true,
		Weight:       weight,
		Connections:  0,
		ReverseProxy: httputil.NewSingleHostReverseProxy(urlstr),
	}, nil
}

func (s *Server) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.Alive
}

func (s *Server) SetAlive(alive bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Alive = alive
}

func (s *Server) GetConnections() int {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.Connections
}

func (s *Server) AddConnection() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Connections++
}

func (s *Server) RemoveConnection() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.Connections > 0 {
		s.Connections--
	}
}
