package server

import (
	"fmt"

	"net"

	"github.com/saravanastar/bango/internal/io"
	"github.com/saravanastar/bango/internal/protocol"
)

type Server struct {
	listener net.Listener
	qchannel chan struct{}
	router   *Router
}

func NewServer(router *Router) Server {
	server := Server{qchannel: make(chan struct{}), router: router}
	return server
}

func (server *Server) Start(port *string) {
	listen, err := net.Listen("tcp", *port)

	if err != nil {
		fmt.Println("error listening the port", err)
	}
	defer listen.Close()
	server.listener = listen

	go server.acceptLoop()
	<-server.qchannel
}

func (server *Server) acceptLoop() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Print("Error Accepting the connection ", err)
		}
		go server.readLoop(conn)

	}

}

func (server *Server) readLoop(conn net.Conn) {
	ioHandler := io.NewIOHandler(conn)
	httpRequest, err := ioHandler.Read()

	defer conn.Close()

	if err != nil {
		fmt.Println("Error reading a request", err)
	}

	routingGuide, err := server.router.GetRoutingGuide(*httpRequest)
	if err != nil {
		fmt.Printf("Error getting routing guide for %v", httpRequest.Http.EndPoint)
	}
	httpResponse := protocol.HttpResponse{Http: protocol.Http{ProtocolVersion: "HTTP/1.1"}, ResponseCode: protocol.NOT_FOUND}
	if routingGuide != nil {
		httpResponseFromController := routingGuide.Handler(httpRequest)
		if httpResponseFromController != nil {
			httpResponse = *httpResponseFromController
		}
	}
	ioHandler.Write(httpResponse)
}
