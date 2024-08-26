package app

import (
	"geoprovider/internal/controller/jsonrpc_v1"
	"geoprovider/internal/controller/rpc_v1"
	"geoprovider/internal/usecase"
	"github.com/wanomir/e"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server interface {
	ListenAndServe() error
}

type ServerConfig struct {
	serviceName string
	port        string
}

type RpcServer struct {
	config     ServerConfig
	controller *rpc_v1.GeoController
}

func NewRpcServer(geoUsecase usecase.GeoServicer, serverConfig ServerConfig) *RpcServer {
	return &RpcServer{
		config:     serverConfig,
		controller: rpc_v1.NewGeoController(geoUsecase),
	}
}

func (s *RpcServer) ListenAndServe() (err error) {
	if err = rpc.RegisterName(s.config.serviceName, s.controller); err != nil {
		return e.Wrap("error registering server:", err)
	}

	listener, err := net.Listen("tcp", ":"+s.config.port)
	if err != nil {
		return e.Wrap("failed to listen", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return e.Wrap("accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}

type JsonRpcServer struct {
	config     ServerConfig
	controller *jsonrpc_v1.GeoController
}

func NewJsonRpcServer(geoUsecase usecase.GeoServicer, serverConfig ServerConfig) *JsonRpcServer {
	return &JsonRpcServer{
		config:     serverConfig,
		controller: jsonrpc_v1.NewGeoController(geoUsecase),
	}
}

func (s *JsonRpcServer) ListenAndServe() (err error) {
	if err = rpc.RegisterName(s.config.serviceName, s.controller); err != nil {
		return e.Wrap("error registering server:", err)
	}

	listener, err := net.Listen("tcp", ":"+s.config.port)
	if err != nil {
		return e.Wrap("failed to listen", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return e.Wrap("accept error:", err)
		}

		go jsonrpc.ServeConn(conn)
	}
}
