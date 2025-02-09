package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/saravanastar/bango/pkg/protocol"
	"github.com/saravanastar/bango/pkg/server"
)

var (
	port      = flag.String("port", ":4221", "Port to listen")
	directory = flag.String("directory", ".", "Directory to pull the content")
)

func main() {
	flag.Parse()
	server := server.NewServer(buildRouter())
	server.Start(port)
	fmt.Println("Server started at port", *port)
}

func buildRouter() *server.Router {
	router := server.NewRoute()
	router.AddRoute(server.RoutingGuide{Url: "/", Method: protocol.GET, Handler: homePage})
	router.AddRoute(server.RoutingGuide{Url: "/echo/:echoString", Method: protocol.GET, Handler: echo})
	router.AddRoute(server.RoutingGuide{Url: "/user-agent", Method: protocol.GET, Handler: userAgent})
	router.AddRoute(server.RoutingGuide{Url: "/files/**", Method: protocol.GET, Handler: serveContent, RouteType: server.CONTENT})
	router.AddRoute(server.RoutingGuide{Url: "/files/**", Method: protocol.POST, Handler: createContent})
	return router
}

func createContent(context *protocol.Context) {
	request := context.Request
	url := request.Http.EndPoint
	exp := regexp.MustCompile(`(/files)(/.*)`)
	rewrittenUrl := exp.ReplaceAllString(url, `$2`)
	fileContent := request.Body

	//read file
	err := os.WriteFile(*directory+rewrittenUrl, []byte(fileContent), 0644)
	if err != nil {
		context.WriteBytes(protocol.INTERNAL_SERVER_ERROR.Code, []byte(err.Error()))
	}
	context.WriteBytes(protocol.CREATED.Code, []byte("File Created"))

}

func serveContent(context *protocol.Context) {
	url := context.Request.Http.EndPoint
	exp := regexp.MustCompile(`(/files)(/.*)`)
	rewrittenUrl := exp.ReplaceAllString(url, `$2`)
	sufftix := strings.Split(rewrittenUrl, ".")

	// create the response
	var contentType string
	if len(sufftix) <= 0 {
		contentType = "application/octet-stream"
	} else {
		contentType = fmt.Sprintf("text/%v", sufftix[len(sufftix)-1])
	}
	fullPath := *directory + rewrittenUrl
	if _, err := os.Stat(fullPath); err != nil {
		context.WriteBytes(protocol.NOT_FOUND.Code, []byte(err.Error()))
		return
	}
	//read file
	file, err := os.ReadFile(fullPath)
	if err != nil {
		context.WriteBytes(protocol.INTERNAL_SERVER_ERROR.Code, []byte(err.Error()))
		return
	}

	context.Data(protocol.OK.Code, contentType, file)
}

func userAgent(context *protocol.Context) {
	content := []byte(strings.Join(context.Request.Http.Headers["User-Agent"], ""))
	context.Data(protocol.OK.Code, "text/plain", content)
}

func echo(context *protocol.Context) {

	content := []byte(context.Request.PathParams["echoString"])
	context.Data(protocol.OK.Code, "text/plain", content)
}

func homePage(context *protocol.Context) {
	content := []byte("All OK!")
	context.Data(protocol.OK.Code, "text/plain", content)
}
