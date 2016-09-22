package networking

import (
	"fmt"
	"net/http"
	"reflect"
	// "github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	RouteMap()[r.URL.Path](w, r)
}

func generateEndPoint(fn http.HandlerFunc) http.HandlerFunc {
	return func(respWrtr http.ResponseWriter, req *http.Request) {
		respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
		fn(respWrtr, req)
	}
}

func ServerEndPoints() {
	fmt.Println("STARTING END POINTS...")
	for _, endPnt := range reflect.ValueOf(RouteMap()).MapKeys() {
		fmt.Println("GENERATING END POINT: ", endPnt)
		http.HandleFunc(endPnt.String(), generateEndPoint(handler))
	}
	fmt.Println("READY TO SERVE\n\n")
}
