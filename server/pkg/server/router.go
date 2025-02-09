package server

import (
	"fmt"
	"strings"

	"github.com/saravanastar/bango/pkg/protocol"
)

type UrlTrie struct {
	urlString     string
	isComplete    bool
	routingGuides []*RoutingGuide
	urlMap        map[string]*UrlTrie
}

type RoutingGuide struct {
	Method  protocol.HttpMethod
	Url     string
	Handler func(request *protocol.Context)
	RouteType
}

type RouteType string

const (
	INTERNAL RouteType = "INTERNAL"
	CONTENT  RouteType = "CONTENT"
	PROXY    RouteType = "PROXY"
)

type IRouter interface {
	AddRoute(routingGuide RoutingGuide) (bool, error)
	GetRoutingGuide(request protocol.HttpRequest) (*RoutingGuide, error)
}

type Router struct {
	routingGuides []RoutingGuide
	urlTrie       *UrlTrie
}

func NewRoute() *Router {
	return &Router{urlTrie: &UrlTrie{urlMap: map[string]*UrlTrie{}}}
}

func (router *Router) AddRoute(routingGuide RoutingGuide) (bool, error) {
	router.routingGuides = append(router.routingGuides, routingGuide)
	return router.addRouteTrie(&routingGuide)
}

func (router *Router) GET(url string, handler func(request *protocol.Context)) (bool, error) {
	return router.AddRoute(RoutingGuide{Url: url, Method: protocol.GET, Handler: handler})
}

func (router *Router) POST(url string, handler func(*protocol.Context)) (bool, error) {
	return router.AddRoute(RoutingGuide{Url: url, Method: protocol.POST, Handler: handler})
}

func (router *Router) PUT(url string, handler func(*protocol.Context)) (bool, error) {
	return router.AddRoute(RoutingGuide{Url: url, Method: protocol.PUT, Handler: handler})
}

func (router *Router) DELETE(url string, handler func(*protocol.Context)) (bool, error) {
	return router.AddRoute(RoutingGuide{Url: url, Method: protocol.DELETE, Handler: handler})
}

// addRouteTrie adds the routing guide to the url trie
func (router *Router) addRouteTrie(routingGuide *RoutingGuide) (bool, error) {
	url := routingGuide.Url
	if url == "" {
		return false, fmt.Errorf("end point can't be empty")
	}

	urlArray := strings.Split(url, "/")
	currentTrie := router.urlTrie

	for index := 0; index < len(urlArray); index++ {
		currentUrlPath := urlArray[index]
		if (index + 1) == len(urlArray) {
			currentTrie.urlString = currentUrlPath
			currentTrie.isComplete = true
			currentTrie.routingGuides = append(currentTrie.routingGuides, routingGuide)
			break
		}
		/**
			1.  UrlPath exist
			2. UrlPath doesn't exist
		**/
		// url path exist
		if currentTrie.urlString == currentUrlPath {
			nextUrlPath := urlArray[index+1]
			if urlTrie, ok := currentTrie.urlMap[nextUrlPath]; ok {
				currentTrie = urlTrie
			} else {
				// url path doesn't exist
				newUrlTrie := &UrlTrie{urlString: nextUrlPath, urlMap: map[string]*UrlTrie{}}
				currentTrie.urlMap[nextUrlPath] = newUrlTrie
				currentTrie = newUrlTrie
			}
		}

	}
	return true, nil
}

// GetRoutingGuide returns the routing guide for the given request
func (router *Router) GetRoutingGuide(request protocol.HttpRequest) (*RoutingGuide, error) {
	url := request.Http.EndPoint
	if url == "" {
		return nil, fmt.Errorf("end point can't be empty")
	}

	urlArray := strings.Split(url, "/")
	currentTrie := router.urlTrie
	routingGuide := router.getRoutingGuideByRecursion(request, 0, urlArray, currentTrie)
	if routingGuide == nil {
		return nil, fmt.Errorf("routing guide not found for the endpoint %v", request.Http.EndPoint)
	}
	return routingGuide, nil
}

// getRoutingGuideByRecursion returns the routing guide for the given request
func (router *Router) getRoutingGuideByRecursion(request protocol.HttpRequest, urlIndex int, urlArray []string, currentTrie *UrlTrie) *RoutingGuide {
	if urlIndex+1 == len(urlArray) && (currentTrie == nil || !currentTrie.isComplete || currentTrie.urlString != urlArray[urlIndex]) {
		return nil
	} else if urlIndex == len(urlArray)-1 && currentTrie.isComplete && (currentTrie.urlString == urlArray[urlIndex] || currentTrie.urlString == "*" || currentTrie.urlString == "**") {
		return router.findMatchingRoutingGuideWithMethod(request, currentTrie)
	} else if currentTrie.urlString == urlArray[urlIndex] {
		// if url index end of the incoming url but current url trie is not complete return nil
		if urlIndex == len(urlArray)-1 && !currentTrie.isComplete {
			return nil
		}
		if urlIndex == len(urlArray)-1 {
			return nil
		}
		nextUrlTrie := currentTrie
		for nextUrlTrie != nil {
			if urlIndex == len(urlArray)-1 {
				return nil
			}
			// if url index end of the incoming url but current url trie is not complete return nil
			upComingTrie, ok := nextUrlTrie.urlMap[urlArray[urlIndex+1]]
			if !ok {
				urlIndex++
				for url, nextTrieVal := range nextUrlTrie.urlMap {

					if url == "**" || url == "*" {
						return router.findMatchingRoutingGuideWithMethod(request, nextTrieVal)
					}
					if len(url) > 1 && url[0] == ':' {
						request.PathParams[url[1:]] = urlArray[urlIndex]
						if urlIndex+1 >= len(urlArray) && nextTrieVal.isComplete {
							return router.findMatchingRoutingGuideWithMethod(request, nextTrieVal)
						}
						if urlIndex+1 >= len(urlArray) && !nextTrieVal.isComplete {
							return nil
						}
						nextUrlTrie = nextTrieVal
						break
					}
				}
			} else {
				return router.getRoutingGuideByRecursion(request, urlIndex+1, urlArray, upComingTrie)
			}

		}
	}
	return nil
}

func (router *Router) findMatchingRoutingGuideWithMethod(request protocol.HttpRequest, urlTrie *UrlTrie) *RoutingGuide {

	for _, currentRoutingGuide := range urlTrie.routingGuides {
		if strings.EqualFold(string(currentRoutingGuide.Method), string(request.Http.Method)) {
			return currentRoutingGuide
		}
	}
	return nil
}
