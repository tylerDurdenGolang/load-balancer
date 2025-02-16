package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/tylerDurdenGolang/load-balancer/internal/balancer"
)

type Server struct {
	addr     string
	balancer balancer.IBalancer
	proxy    *httputil.ReverseProxy
}

func NewServer(addr string, b balancer.IBalancer) *Server {
	// Заглушка — нам понадобится подменять адрес бэкенда внутри Director
	director := func(req *http.Request) {}
	proxy := &httputil.ReverseProxy{Director: director}

	return &Server{
		addr:     addr,
		balancer: b,
		proxy:    proxy,
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/", s.handleRequest)
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	backend, err := s.balancer.GetBackend()
	if err != nil {
		http.Error(w, "No healthy backends", http.StatusServiceUnavailable)
		return
	}

	backendURL, err := url.Parse("http://" + backend)
	if err != nil {
		log.Printf("error parsing backend URL: %v", err)
		http.Error(w, "Invalid backend", http.StatusInternalServerError)
		return
	}

	// Переопределяем Director, чтобы он указывал на выбранный бэкенд
	s.proxy.Director = func(req *http.Request) {
		// Сохраняем метод, заголовки и путь:
		req.URL.Scheme = backendURL.Scheme
		req.URL.Host = backendURL.Host
		// В requestURI хранится полный путь (с query), так что можно использовать:
		// либо req.URL.Path, req.URL.RawQuery — если нужно подменять что-то точечно
		// либо напрямую req.URL.Opaque и т.п. — зависит от структуры.
		// Но обычно достаточно:
		req.URL.Path = r.URL.Path
		req.URL.RawQuery = r.URL.RawQuery
		if _, ok := req.Header["X-Forwarded-Host"]; !ok {
			req.Header.Set("X-Forwarded-Host", req.Host)
		}
		req.Host = backendURL.Host
	}

	s.proxy.ServeHTTP(w, r)
}
