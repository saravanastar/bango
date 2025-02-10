package server

import (
	"bytes"
	"net"
	"testing"

	"github.com/saravanastar/bango/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

type mockConn struct {
	net.Conn
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

type mockRouter struct{}

func (r *mockRouter) AddRoute(routingGuide RoutingGuide) (bool, error) {
	// Mock implementation
	return true, nil
}

func (r *mockRouter) GetRoutingGuide(request protocol.HttpRequest) (*RoutingGuide, error) {
	return &RoutingGuide{
		Handler: func(context *protocol.Context) {
			context.Response = &protocol.HttpResponse{
				Http: protocol.Http{
					ProtocolVersion: "HTTP/1.1",

					Headers: map[string][]string{
						"Content-Type": {"text/plain"},
					},
				},
				Body: []byte("Hello, World!"),
				ResponseCode: protocol.ResponseCodes{
					Code:           200,
					ResponseString: "OK",
				},
			}
		},
	}, nil
}

func TestReadLoop(t *testing.T) {
	request := "GET /hello HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello"
	mock := &mockConn{
		readBuffer:  bytes.NewBufferString(request),
		writeBuffer: bytes.NewBuffer([]byte{}),
	}
	router := &mockRouter{}
	server := NewServer(router)
	server.readLoop(mock)

	expectedResponse := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, World!"
	assert.Equal(t, expectedResponse, mock.writeBuffer.String())
}
