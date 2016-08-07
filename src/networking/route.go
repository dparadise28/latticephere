package networking

import (
	"net/http"
)

//create type for consistent function signatures in routing map
type fn func (http.ResponseWriter, *http.Request)

func RouteMap()(map[string]fn){
	return map[string] fn {
		"/"  : CheckPath, 
		"/h2": ShowRequestInfoHandler,
	}
}
