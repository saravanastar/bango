package server_test

import (
	"fmt"
	"testing"

	"github.com/saravanastar/bango/internal/protocol"
	"github.com/saravanastar/bango/internal/server"
)

func TestAddRouteSuccessScenario(t *testing.T) {

	// prepare
	routingGuide := &server.RoutingGuide{Url: "/echo/{echoString}", Method: protocol.GET}

	router := server.NewRoute()

	// Execute
	doesAdded, err := router.AddRoute(*routingGuide)

	// Assert
	if err != nil || doesAdded != true {
		t.Errorf("Should add the route %v", routingGuide.Url)
	}

}

func TestAddRoutAndTestWithGetRoutingGuide(t *testing.T) {
	// prepare
	routingGuide := &server.RoutingGuide{Url: "/echo/{echoString}", Method: protocol.GET}

	router := server.NewRoute()

	// Execute
	doesAdded, err := router.AddRoute(*routingGuide)

	// Assert
	if err != nil || doesAdded != true {
		t.Errorf("Should add the route %v", routingGuide.Url)
	}
	http := protocol.Http{EndPoint: "/echo/test"}
	httpRequest := protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, _ := router.GetRoutingGuide(httpRequest)

	if responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}

}

func TestAddRoutAndTestWithGetRoutingGuideAndMoreRoutes(t *testing.T) {
	// prepare
	routingGuide1 := &server.RoutingGuide{Url: "/echo/{echoString}", Method: protocol.GET}
	routingGuide2 := &server.RoutingGuide{Url: "/", Method: protocol.GET}
	routingGuide3 := &server.RoutingGuide{Url: "/contents/app/index.html", Method: protocol.GET}

	router := server.NewRoute()

	// Execute
	doesAdded1, err1 := router.AddRoute(*routingGuide1)
	doesAdded2, err2 := router.AddRoute(*routingGuide2)
	doesAdded3, err3 := router.AddRoute(*routingGuide3)

	// Assert
	if err1 != nil || doesAdded1 != true {
		t.Errorf("Should add the route %v", routingGuide1.Url)
	}

	// Assert
	if err2 != nil || doesAdded2 != true {
		t.Errorf("Should add the route %v", routingGuide2.Url)
	}

	// Assert
	if err3 != nil || doesAdded3 != true {
		t.Errorf("Should add the route %v", routingGuide3.Url)
	}

	http := protocol.Http{EndPoint: "/echo/test"}
	httpRequest := protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, _ := router.GetRoutingGuide(httpRequest)

	if responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}
	fmt.Println("Printing the path params", httpRequest.PathParams)
	if httpRequest.PathParams["echoString"] != "test" {
		t.Errorf("Should copy the pathParams for the route %v", http.EndPoint)
	}

	http = protocol.Http{EndPoint: routingGuide2.Url}
	httpRequest = protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, _ = router.GetRoutingGuide(httpRequest)

	if responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}

	http = protocol.Http{EndPoint: routingGuide3.Url}
	httpRequest = protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, _ = router.GetRoutingGuide(httpRequest)

	fmt.Println("Response Routing guid3", responseRoutingGuide)
	if responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}

}
