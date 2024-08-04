package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/saravanastar/bango/internal/protocol"
)

type RoutingGuide struct {
	Method  protocol.HttpMethod
	Url     string
	Handler func(request *protocol.HttpRequest) *protocol.HttpResponse
	RouteType
}

type RouteType string

const (
	INTERNAL RouteType = "INTERNAL"
	CONTENT  RouteType = "CONTENT"
	PROXY    RouteType = "PROXY"
)

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

func (router *Router) addRouteTrie(routingGuide *RoutingGuide) (bool, error) {
	url := routingGuide.Url
	if url == "" {
		return false, errors.New("url can't be empty")
	}

	urlArray := strings.Split(url, "/")
	currentTrie := router.urlTrie

	for index := 0; index < len(urlArray); index++ {
		currentUrlPath := urlArray[index]
		// if (index+1) != len(urlArray) && currentUrlPath == "" {
		// 	continue
		// }
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

/*
*
/api/request/users/{userId}
*
*/
func (router *Router) GetRoutingGuide(request protocol.HttpRequest) (*RoutingGuide, error) {
	url := request.Http.EndPoint
	if url == "" {
		return nil, errors.New("url can't be empty")
	}

	urlArray := strings.Split(url, "/")
	currentTrie := router.urlTrie

	return router.getRoutingGuideByRecursion(request, 0, urlArray, currentTrie), nil

}

func (router *Router) getRoutingGuideByRecursion(request protocol.HttpRequest, urlIndex int, urlArray []string, currentTrie *UrlTrie) *RoutingGuide {

	if urlIndex+1 == len(urlArray) && (currentTrie == nil || !currentTrie.isComplete || currentTrie.urlString != urlArray[urlIndex]) {
		return nil
	} else if urlIndex+1 == len(urlArray) && currentTrie.isComplete && currentTrie.urlString == urlArray[urlIndex] {
		return router.findMatchingRoutingGuideWithMethod(request, currentTrie)
	} else if currentTrie.urlString == urlArray[urlIndex] {
		// if url index end of the incoming url but current url trie is not complete return nil
		if urlIndex+1 == len(urlArray) && !currentTrie.isComplete {
			return nil
		}
		// if next url part match not found in the map then return nil
		nextUrlTrie, ok := currentTrie.urlMap[urlArray[urlIndex+1]]
		if !ok {
			// if match not found check whether any path param is present
			for url, nextTrieVal := range currentTrie.urlMap {
				if url == "**" {
					return router.findMatchingRoutingGuideWithMethod(request, nextTrieVal)
				}
				if len(url) > 1 && url[0] == '{' {
					fmt.Printf("Formating %v", len(urlArray))
					// Check the whether its end of input url slice and trie also complete, if so return it
					if urlIndex+1 == len(urlArray)-1 && nextTrieVal.isComplete {
						// copy the path params and value
						request.PathParams[nextTrieVal.urlString[1:len(nextTrieVal.urlString)-1]] = urlArray[urlIndex+1]
						return router.findMatchingRoutingGuideWithMethod(request, nextTrieVal)
					}
					// if not end of the array, then skip the pathparam matching and move cursor to further more
					if urlIndex+2 < len(urlArray) {
						if furhterMore, ok := nextTrieVal.urlMap[urlArray[urlIndex+2]]; ok {
							pathParamRouteExist := router.getRoutingGuideByRecursion(request, urlIndex+2, urlArray, furhterMore)
							if pathParamRouteExist != nil {
								request.PathParams[url[0:len(url)-1]] = urlArray[urlIndex]
								return pathParamRouteExist
							}
						}
					}

				}
			}
		}
		// if match found in the next url part, call recursive function to move further
		return router.getRoutingGuideByRecursion(request, urlIndex+1, urlArray, nextUrlTrie)
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

type UrlTrie struct {
	urlString     string
	isComplete    bool
	routingGuides []*RoutingGuide
	urlMap        map[string]*UrlTrie
}
