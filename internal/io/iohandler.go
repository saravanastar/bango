package io

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/saravanastar/bango/internal/protocol"
)

type IOHandler struct {
	reader *bufio.Reader
	writer io.Writer
}

const (
	SPACE = " "
)

func NewIOHandler(con net.Conn) *IOHandler {
	return &IOHandler{reader: bufio.NewReader(con), writer: con}
}
func (ioHandler *IOHandler) Write(response protocol.HttpResponse) {
	body := []byte(response.Http.Body)
	if response.Http.Headers != nil {
		response.Http.Headers["Content-Length"] = []string{strconv.Itoa(binary.Size(body))}
	}
	var responseByte []byte
	responseByte = append(responseByte, response.Http.ProtocolVersion...)
	responseByte = append(responseByte, []byte(SPACE)...)
	responseByte = append(responseByte, []byte(strconv.Itoa(response.ResponseCode.Code))...)
	responseByte = append(responseByte, []byte(SPACE)...)
	responseByte = append(responseByte, []byte(response.ResponseCode.ResponseString)...)
	responseByte = append(responseByte, []byte("\r\n")...)
	for headerName, headerValue := range response.Http.Headers {
		responseByte = append(responseByte, []byte(headerName)...)
		responseByte = append(responseByte, []byte(": ")...)
		responseByte = append(responseByte, []byte(strings.Join(headerValue, ";"))...)
		responseByte = append(responseByte, []byte("\r\n")...)
	}
	responseByte = append(responseByte, []byte("\r\n")...)
	// responseByte = append(responseByte, []byte("\r\n")...)
	responseByte = append(responseByte, []byte(body)...)
	ioHandler.writer.Write(responseByte)
}

func (ioHandler *IOHandler) Read() (*protocol.HttpRequest, error) {
	http := protocol.Http{}
	line, err := ioHandler.readLine()

	if err != nil {
		return nil, err
	}
	http.Method, http.EndPoint, http.ProtocolVersion, err = formatProtocol(line)

	if err != nil {
		return nil, err
	}

	headerPointer, err := ioHandler.readHeader()

	if err != nil {
		return nil, err
	}

	http.Headers = *headerPointer
	contentLengthString, ok := http.Headers["Content-Length"]
	if ok {
		contentLength, err := strconv.ParseInt(contentLengthString[0], 10, 32)
		if err != nil {
			return nil, err
		}
		body, err := ioHandler.readByByteLength(int(contentLength))
		if err != nil {
			return nil, err
		}
		http.Body = string(body)
	}

	return &protocol.HttpRequest{Http: http, PathParams: make(map[string]string), QueryParams: make(map[string]string)}, nil
}

func (ioHandler *IOHandler) readHeader() (*protocol.Header, error) {
	line, err := ioHandler.readLine()
	if err != nil {
		return nil, err
	}

	headers := protocol.Header{}
	for len(strings.TrimSpace(string(line))) > 0 {
		key, value, err := formatHeader(line)
		if err != nil {
			return nil, err
		}
		headers[key] = value
		line, err = ioHandler.readLine()
		if err != nil {
			return nil, err
		}
	}
	return &headers, nil
}

func formatHeader(line []byte) (string, []string, error) {
	lineString := string(line)
	index := strings.Index(lineString, ":")

	return lineString[0:index], strings.Split(lineString[index+2:], ","), nil
}

func formatProtocol(line []byte) (protocol.HttpMethod, string, string, error) {
	lineString := string(line)

	headerLine := strings.Split(lineString, " ")

	if len(headerLine) < 3 {
		return "", "", "", errors.New("head line protocol length must be 3")
	}
	return protocol.HttpMethod(headerLine[0]), headerLine[1], headerLine[2], nil
}

func (ioHandler *IOHandler) readByByteLength(byteLen int) ([]byte, error) {
	line := make([]byte, byteLen)
	_, err := ioHandler.reader.Read(line)

	if err != nil {
		return nil, err
	}
	return line, nil
}

func (ioHandler *IOHandler) readLine() ([]byte, error) {
	line := make([]byte, 0)
	for {
		b, err := ioHandler.reader.ReadByte()

		if err != nil {
			return nil, err
		}

		if b == '\r' {
			_, err = ioHandler.reader.ReadByte()

			if err != nil {
				return nil, err
			}
			break
		}

		line = append(line, b)
	}
	return line, nil
}
