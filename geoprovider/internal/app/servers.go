package app

import (
	"geoprovider/internal/controller/jsonrpc_v1"
	"geoprovider/internal/controller/rpc_v1"
	"geoprovider/internal/usecase"
	"github.com/wanomir/e"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server interface {
	ListenAndServe() error
	Shutdown()
}

type BaseServer struct {
	name       string
	port       string
	shutdownCh chan bool
	listener   net.Listener
}

func NewBaseServer(name, port string) *BaseServer {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("failed to listen", err)
		return nil
	}

	return &BaseServer{
		name:       name,
		port:       port,
		listener:   listener,
		shutdownCh: make(chan bool),
	}
}

func (s *BaseServer) ListenAndServe() error {
	return nil
}

func (s *BaseServer) Shutdown() {
	s.shutdownCh <- true
}

type RpcServer struct {
	*BaseServer
	controller *rpc_v1.GeoController
}

func NewRpcServer(service usecase.GeoServicer, name, port string) *RpcServer {
	return &RpcServer{
		BaseServer: NewBaseServer(name, port),
		controller: rpc_v1.NewGeoController(service),
	}
}

func (s *RpcServer) ListenAndServe() (err error) {
	if err = rpc.RegisterName(s.BaseServer.name, s.controller); err != nil {
		return e.Wrap("error registering server:", err)
	}

	for {
		select {
		case <-s.shutdownCh:
			log.Println("shutting down rpc server...")
			return nil

		default:
			conn, err := s.BaseServer.listener.Accept()
			if err != nil {
				return e.Wrap("accept error:", err)
			}

			go rpc.ServeConn(conn)
		}
	}
}

type JsonRpcServer struct {
	*BaseServer
	controller *jsonrpc_v1.GeoController
}

func NewJsonRpcServer(service usecase.GeoServicer, name, port string) *JsonRpcServer {
	return &JsonRpcServer{
		BaseServer: NewBaseServer(name, port),
		controller: jsonrpc_v1.NewGeoController(service),
	}
}

func (s *JsonRpcServer) ListenAndServe() (err error) {
	if err = rpc.RegisterName(s.BaseServer.name, s.controller); err != nil {
		return e.Wrap("error registering server:", err)
	}

	for {
		select {
		case <-s.shutdownCh:
			log.Println("shutting down json-rpc server...")
			return nil

		default:
			conn, err := s.BaseServer.listener.Accept()
			if err != nil {
				return e.Wrap("accept error:", err)
			}

			go jsonrpc.ServeConn(conn)
		}
	}
}
