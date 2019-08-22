package api

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/spacemeshos/collector/api/pb"
	"github.com/spacemeshos/go-spacemesh/log"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type JSONHTTPServer struct {
	Port   uint
	server *http.Server
	ctx    context.Context
	stop   chan bool
}

// NewJSONHTTPServer creates a new json http server.
func NewJSONHTTPServer(port int) *JSONHTTPServer {
	return &JSONHTTPServer{Port: uint(port), stop: make(chan bool)}
}

// StopService stops the server.
func (s JSONHTTPServer) StopService() {
	log.Debug("Stopping json-http service...")
	s.stop <- true
}

// Listens on gracefully stopping the server in the same routine.
func (s JSONHTTPServer) listenStop() {
	<-s.stop
	log.Debug("Shutting down json API server...")
	if err := s.server.Shutdown(s.ctx); err != nil {
		log.Error("Error during shutdown json API server : %v", err)
	}
}

// StartService starts the json api server and listens for status (started, stopped).
func (s JSONHTTPServer) StartService(status chan bool) {
	go s.startInternal(status)
}

func (s JSONHTTPServer) startInternal(status chan bool) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.ctx = ctx

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// register the http server on the local grpc server
	portStr := strconv.Itoa(int(s.Port))

	const endpoint = "api_endpoint"
	var echoEndpoint string

	fl := flag.Lookup(endpoint)
	if fl != nil {
		echoEndpoint = fl.Value.String()
	} else {
		echoEndpoint = *flag.String(endpoint, "localhost:"+portStr, "endpoint of api grpc service")
	}
	if err := gw.RegisterCollectorServiceHandlerFromEndpoint(ctx, mux, echoEndpoint, opts); err != nil {
		log.Error("failed to register http endpoint with grpc", err)
	}

	addr := ":" + strconv.Itoa(int(s.Port))

	log.Debug("json API listening on port %d", s.Port)

	go func() { s.listenStop() }()

	if status != nil {
		status <- true
	}

	s.server = &http.Server{Addr: addr, Handler: mux}
	err := s.server.ListenAndServe()

	if err != nil {
		log.Debug("listen and serve stopped with status. %v", err)
	}

	if status != nil {
		status <- true
	}
}
