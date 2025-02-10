package io

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

func TestNewIOHandler(t *testing.T) {
	mock := &mockConn{
		readBuffer:  bytes.NewBuffer([]byte{}),
		writeBuffer: bytes.NewBuffer([]byte{}),
	}
	ioHandler := NewIOHandler(mock)
	assert.NotNil(t, ioHandler)
	assert.NotNil(t, ioHandler.reader)
	assert.NotNil(t, ioHandler.writer)
}

func TestWrite(t *testing.T) {
	mock := &mockConn{
		readBuffer:  bytes.NewBuffer([]byte{}),
		writeBuffer: bytes.NewBuffer([]byte{}),
	}
	ioHandler := NewIOHandler(mock)
	response := protocol.HttpResponse{
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
	ioHandler.Write(response)
	expected := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, World!"
	assert.Equal(t, expected, mock.writeBuffer.String())
}

func TestRead(t *testing.T) {
	request := "GET /hello HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello"
	mock := &mockConn{
		readBuffer:  bytes.NewBufferString(request),
		writeBuffer: bytes.NewBuffer([]byte{}),
	}
	ioHandler := NewIOHandler(mock)
	httpRequest, err := ioHandler.Read()
	assert.NoError(t, err)
	assert.Equal(t, protocol.GET, httpRequest.Http.Method)
	assert.Equal(t, "/hello", httpRequest.Http.EndPoint)
	assert.Equal(t, "HTTP/1.1", httpRequest.Http.ProtocolVersion)
	assert.Equal(t, "Hello", string(httpRequest.Body))
}

func TestReadHeader(t *testing.T) {
	request := "Content-Type: text/plain\r\nContent-Length: 5\r\n\r\n"
	mock := &mockConn{
		readBuffer:  bytes.NewBufferString(request),
		writeBuffer: bytes.NewBuffer([]byte{}),
	}
	ioHandler := NewIOHandler(mock)
	headers, err := ioHandler.readHeader()
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", (*headers)["Content-Type"][0])
	assert.Equal(t, "5", (*headers)["Content-Length"][0])
}

func TestReadLine(t *testing.T) {
	request := "GET /hello HTTP/1.1\r\n"
	mock := &mockConn{
		readBuffer:  bytes.NewBufferString(request),
		writeBuffer: bytes.NewBuffer([]byte{}),
	}
	ioHandler := NewIOHandler(mock)
	line, err := ioHandler.readLine()
	assert.NoError(t, err)
	assert.Equal(t, "GET /hello HTTP/1.1", string(line))
}
