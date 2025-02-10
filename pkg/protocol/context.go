package protocol

import (
	"encoding/json"
	"strconv"
)

type Context struct {
	Request  *HttpRequest
	Response *HttpResponse
}

func NewContext(request *HttpRequest, response *HttpResponse) *Context {
	return &Context{Request: request, Response: response}
}

func (context *Context) JSON(code int, data interface{}) {
	context.Response.ResponseCode.Code = code
	jsonData, err := json.Marshal(data)
	if err != nil {
		context.Response.ResponseCode.Code = 500
		context.Response.Body = []byte(`{"error": "Internal Server Error"}`)
		return
	}
	context.Response.Body = jsonData
	context.Response.Http.Headers["Content-Type"] = []string{"application/json"}
	context.Response.Http.Headers["Content-Length"] = []string{strconv.Itoa(len(context.Response.Body))}
}

func (context *Context) String(code int, data string) {
	context.Response.ResponseCode.Code = code
	context.Response.Body = []byte(data)
	context.Response.Http.Headers["Content-Type"] = []string{"text/plain"}
	context.Response.Http.Headers["Content-Length"] = []string{strconv.Itoa(len(context.Response.Body))}
}

func (context *Context) HTML(code int, data string) {
	context.Response.ResponseCode.Code = code
	context.Response.Body = []byte(data)
	context.Response.Http.Headers["Content-Type"] = []string{"text/html"}
	context.Response.Http.Headers["Content-Length"] = []string{strconv.Itoa(len(context.Response.Body))}
}

func (context *Context) Status(code int) {
	context.Response.ResponseCode.Code = code
}

func (context *Context) Redirect(code int, url string) {
	context.Response.ResponseCode.Code = code
	context.Response.Http.Headers["Location"] = []string{url}
}

func (context *Context) SetHeader(key, value string) {
	context.Response.Http.Headers[key] = []string{value}
}

func (context *Context) GetHeader(key string) string {
	return context.Request.Http.Headers[key][0]
}

func (context *Context) GetQueryParam(key string) string {
	return string(context.Request.QueryParams[key][0])
}

func (context *Context) GetPathParam(key string) string {
	return string(context.Request.PathParams[key][0])
}

func (context *Context) GetRequestBody() []byte {
	return context.Request.Body
}

func (context *Context) GetRequestMethod() string {
	return string(context.Request.Http.Method)
}

func (context *Context) GetRequestPath() string {
	return context.Request.Http.EndPoint
}

func (context *Context) GetRequestProtocol() string {
	return context.Request.Http.ProtocolVersion
}

func (context *Context) GetRequestHeaders() map[string][]string {
	return context.Request.Http.Headers
}

func (context *Context) WriteBytes(code int, data []byte) {
	context.Response.Body = data
	context.Response.Http.Headers["Content-Length"] = []string{strconv.Itoa(len(data))}
	context.Response.ResponseCode.Code = code
}

func (context *Context) Data(responsCode int, contentType string, data []byte) {
	context.Response.Body = data
	context.Response.ResponseCode.Code = responsCode
	context.Response.Http.Headers["Content-Length"] = []string{strconv.Itoa(len(data))}
	context.Response.Http.Headers["Content-Type"] = []string{contentType}
}
