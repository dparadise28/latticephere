package networking

import (
	"net/http"
	"tools"
)

//container for api methods
type API struct{}

//create type for consistent function signatures in routing map
type fn func(http.ResponseWriter, *http.Request)

func RouteMap() map[string]fn {
	return map[string]fn{
		"/":          CheckPath,
		"/h2":        ShowRequestInfoHandler,
		"/transform": tools.RemodelJ,
	}
}
