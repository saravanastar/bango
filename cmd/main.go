package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/saravanastar/bango/internal/protocol"
	"github.com/saravanastar/bango/internal/server"
)

var (
	port      = flag.String("port", ":4221", "Port to listen")
	directory = flag.String("directory", "/tmp", "Directory to pull the content")
)

func main() {
	flag.Parse()
	server := server.NewServer(buildRouter())
	server.Start(port)
}

func buildRouter() *server.Router {
	router := server.NewRoute()
	router.AddRoute(server.RoutingGuide{Url: "/", Method: protocol.GET, Handler: homePage})
	router.AddRoute(server.RoutingGuide{Url: "/echo/{echoString}", Method: protocol.GET, Handler: echo})
	router.AddRoute(server.RoutingGuide{Url: "/user-agent", Method: protocol.GET, Handler: userAgent})
	router.AddRoute(server.RoutingGuide{Url: "/files/**", Method: protocol.GET, Handler: servContent})
	router.AddRoute(server.RoutingGuide{Url: "/files/**", Method: protocol.POST, Handler: createContent})
	return router
}

func createContent(request *protocol.HttpRequest) *protocol.HttpResponse {
	url := request.Http.EndPoint
	exp := regexp.MustCompile(`(/files)(/.*)`)
	rewrittenUrl := exp.ReplaceAllString(url, `$2`)
	fileContent := request.Http.Body

	// create the response
	httpResponse := protocol.NewHttpResponse(request)
	// httpResponse.Http.Headers["Content-Type"] = []string{"application/octet-stream"}

	//read file
	err := os.WriteFile(*directory+rewrittenUrl, []byte(fileContent), 0644)
	if err != nil {
		httpResponse.Http.Headers["Content-Type"] = []string{"text/plain"}
		httpResponse.ResponseCode = protocol.INTERNAL_SERVER_ERROR
		httpResponse.Http.Body = err.Error()
		return nil
	}
	httpResponse.ResponseCode = protocol.CREATED

	return httpResponse
}

func servContent(request *protocol.HttpRequest) *protocol.HttpResponse {
	url := request.Http.EndPoint
	exp := regexp.MustCompile(`(/files)(/.*)`)
	rewrittenUrl := exp.ReplaceAllString(url, `$2`)
	sufftix := strings.Split(rewrittenUrl, ".")

	// create the response
	httpResponse := protocol.NewHttpResponse(request)
	if len(sufftix) <= 0 {
		httpResponse.Http.Headers["Content-Type"] = []string{"application/octet-stream"}
	} else {
		contentType := fmt.Sprintf("text/%v", sufftix[len(sufftix)-1])
		httpResponse.Http.Headers["Content-Type"] = []string{contentType}
	}

	//read file
	file, err := os.ReadFile(*directory + rewrittenUrl)
	if err != nil {
		// httpResponse.Http.Headers["Content-Type"] = []string{"text/plain"}
		// httpResponse.ResponseCode = protocol.INTERNAL_SERVER_ERROR
		// httpResponse.Http.Body = err.Error()
		return nil
	}
	httpResponse.ResponseCode = protocol.OK
	httpResponse.Http.Body = string(file)
	return httpResponse
}

func userAgent(request *protocol.HttpRequest) *protocol.HttpResponse {
	httpResponse := protocol.NewHttpResponse(request)
	httpResponse.Http.Body = strings.Join(request.Http.Headers["User-Agent"], "")
	httpResponse.Http.Headers["Content-Type"] = []string{"text/plain"}
	httpResponse.ResponseCode = protocol.OK
	return httpResponse
}

func echo(request *protocol.HttpRequest) *protocol.HttpResponse {

	httpResponse := protocol.NewHttpResponse(request)
	httpResponse.Http.Body = request.PathParams["echoString"]
	httpResponse.Http.Headers["Content-Type"] = []string{"text/plain"}
	httpResponse.ResponseCode = protocol.OK
	return httpResponse
}

func homePage(request *protocol.HttpRequest) *protocol.HttpResponse {
	httpResponse := protocol.NewHttpResponse(request)
	httpResponse.Http.Body = "All OK!"
	httpResponse.ResponseCode = protocol.OK
	return httpResponse
}
