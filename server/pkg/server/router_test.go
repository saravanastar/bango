package server_test

import (
	"testing"

	"github.com/saravanastar/bango/pkg/protocol"
	"github.com/saravanastar/bango/pkg/server"
)

func TestAddRouteSuccessScenario(t *testing.T) {
	// prepare
	routingGuide := &server.RoutingGuide{Url: "/echo/{echoString}", Method: protocol.GET}
	router := server.NewRoute()

	// Execute
	doesAdded, err := router.AddRoute(*routingGuide)

	// Assert
	if err != nil || !doesAdded {
		t.Errorf("Should add the route %v", routingGuide.Url)
	}
}

func TestAddRouteAndTestWithGetRoutingGuide(t *testing.T) {
	// prepare
	routingGuide := &server.RoutingGuide{Url: "/echo/:echoString", Method: protocol.GET}
	router := server.NewRoute()

	// Execute
	doesAdded, err := router.AddRoute(*routingGuide)

	// Assert
	if err != nil || !doesAdded {
		t.Errorf("Should add the route %v", routingGuide.Url)
	}

	http := protocol.Http{EndPoint: "/echo/test", Method: protocol.GET}
	httpRequest := protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, err := router.GetRoutingGuide(httpRequest)

	if err != nil || responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}
}

func TestAddRouteAndTestWithGetRoutingGuideAndMoreRoutes(t *testing.T) {
	// prepare
	routingGuide1 := &server.RoutingGuide{Url: "/echo/:echoString/user", Method: protocol.GET}
	routingGuide2 := &server.RoutingGuide{Url: "/", Method: protocol.GET}
	routingGuide3 := &server.RoutingGuide{Url: "/contents/app/index.html", Method: protocol.GET}

	router := server.NewRoute()

	// Execute
	doesAdded1, err1 := router.AddRoute(*routingGuide1)
	doesAdded2, err2 := router.AddRoute(*routingGuide2)
	doesAdded3, err3 := router.AddRoute(*routingGuide3)

	// Assert
	if err1 != nil || !doesAdded1 {
		t.Errorf("Should add the route %v", routingGuide1.Url)
	}
	if err2 != nil || !doesAdded2 {
		t.Errorf("Should add the route %v", routingGuide2.Url)
	}
	if err3 != nil || !doesAdded3 {
		t.Errorf("Should add the route %v", routingGuide3.Url)
	}

	http := protocol.Http{EndPoint: "/echo/test/user", Method: protocol.GET}
	httpRequest := protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, err := router.GetRoutingGuide(httpRequest)

	if err != nil || responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}
	if httpRequest.PathParams["echoString"] != "test" {
		t.Errorf("Should copy the pathParams for the route %v", http.EndPoint)
	}

	http = protocol.Http{EndPoint: routingGuide2.Url, Method: protocol.GET}
	httpRequest = protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, err = router.GetRoutingGuide(httpRequest)

	if err != nil || responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}

	http = protocol.Http{EndPoint: routingGuide3.Url, Method: protocol.GET}
	httpRequest = protocol.HttpRequest{Http: http, PathParams: make(map[string]string)}
	responseRoutingGuide, err = router.GetRoutingGuide(httpRequest)

	if err != nil || responseRoutingGuide == nil {
		t.Errorf("Should find the route %v", http.EndPoint)
	}
}
