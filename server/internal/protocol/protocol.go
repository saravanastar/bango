package protocol

type Http struct {
	ProtocolVersion string
	EndPoint        string
	Headers         Header
	Method          HttpMethod
}

type HttpRequest struct {
	Http        Http
	PathParams  map[string]string
	QueryParams map[string]string
	Body        []byte
}

type HttpResponse struct {
	Http         Http
	ResponseCode ResponseCodes
	Body         []byte
}

func NewHttpResponse(request *HttpRequest) *HttpResponse {
	http := &Http{ProtocolVersion: "HTTP/1.1", Headers: make(Header)}

	return &HttpResponse{Http: *http}
}

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

type Header map[string][]string

type ResponseCodes struct {
	Code           int
	ResponseString string
}

var (
	OK                    ResponseCodes = ResponseCodes{Code: 200, ResponseString: "OK"}
	CREATED               ResponseCodes = ResponseCodes{Code: 201, ResponseString: "Created"}
	NOT_FOUND             ResponseCodes = ResponseCodes{Code: 404, ResponseString: "Not Found"}
	INTERNAL_SERVER_ERROR ResponseCodes = ResponseCodes{Code: 500, ResponseString: "Internal Server Error"}
)

type Context struct {
	Request  *HttpRequest
	Response *HttpResponse
}
