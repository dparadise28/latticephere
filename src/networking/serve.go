package networking

import (
	"fmt"
	"reflect"
	"net/http"
	// "github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, r.URL.Path)
	RouteMap()[r.URL.Path](w, r)
}

func generateEndPoint(fn http.HandlerFunc) http.HandlerFunc{
	fmt.Println("1")
	return func(respWrtr http.ResponseWriter, req *http.Request){
		fmt.Println("2")
		fn(respWrtr, req)
	}
}

func ServerEndPoints(){
	for _, endPnt := range reflect.ValueOf(RouteMap()).MapKeys(){
		fmt.Println(endPnt)
		http.HandleFunc(endPnt.String(), generateEndPoint(handler))
	}
}
