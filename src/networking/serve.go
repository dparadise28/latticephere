package networking

import (
	"fmt"
	"net/http"
	"reflect"
	// "github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	APIRouteMap()[r.URL.Path](w, r)
}

func serveFile(res http.ResponseWriter, req *http.Request, file string) {
	http.ServeFile(res, req, file)
}

func generateAPIEndPoint(fn http.HandlerFunc) http.HandlerFunc {
	return func(respWrtr http.ResponseWriter, req *http.Request) {
		respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
		fn(respWrtr, req)
	}
}

func generateUIEndPoint(file string) http.HandlerFunc {
	return func(respWrtr http.ResponseWriter, req *http.Request) {
		respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
		serveFile(respWrtr, req, file)
	}
}

func ServerEndPoints() {
	fmt.Println("STARTING API END POINTS...")
	for _, endPnt := range reflect.ValueOf(APIRouteMap()).MapKeys() {
		fmt.Println("GENERATING END POINT: ", endPnt)
		http.HandleFunc(endPnt.String(), generateAPIEndPoint(handler))
	}

	fmt.Println("\n\nSTARTING UI END POINTS...")
	for endPnt, file := range UIRouteMap() {
		fmt.Println("GENERATING UI END POINT: ", endPnt)
		http.HandleFunc(endPnt, generateUIEndPoint(file))
	}
	fmt.Println("READY TO SERVE\n\n")
}
