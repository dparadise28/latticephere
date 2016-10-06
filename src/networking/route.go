package networking

import (
	"net/http"
	"tools"
)

//container for api methods
type API struct{}

//create type for consistent function signatures in routing map
type fn func(http.ResponseWriter, *http.Request)

func APIRouteMap() map[string]fn {
	return map[string]fn{
		//"/":          CheckPath,
		"/h2":        ShowRequestInfoHandler,
		"/transform": tools.RemodelJ,
	}
}

func UIRouteMap() map[string]string {
	return map[string]string{
		"/": "views/startbootstrap-sb-admin/index.html",
	}
}
