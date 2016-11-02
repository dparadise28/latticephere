package networking

import (
	//"github.com/gorilla/websocket"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"reflect"
)

func handler(w http.ResponseWriter, r *http.Request) {
	APIRouteMap()[r.URL.Path](w, r)
}

func Log(message string) {
	var origin = "http://72.69.174.66:8000"
	var url = "ws://72.69.174.66:8000/api/log"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	mesg := []byte(message) //"hello, world!")
	_, err = ws.Write(mesg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", mesg)

	var msg = make([]byte, 512)
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg)
}

/*
func Log1(message string) {
	c, err := upgrader.Upgrade(APILogConnection, APILogRequest, nil)
	if err != nil {
		log.Print("upgrade:", err)
	} else {
		for {
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}
}
*/
func generateAPIEndPoint(fn http.HandlerFunc) http.HandlerFunc {
	return func(respWrtr http.ResponseWriter, req *http.Request) {
		log.Println("!!!!!!!!!!!     SETTING    {Access-Control-Allow-Origin: *}     !!!!!!!")
		respWrtr.Header().Set("Access-Control-Allow-Origin", "*")
		respWrtr.Header().Set("Sec-Websocket-Version", "13")
		fn(respWrtr, req)
	}
}

func ServeEndPoints() {
	for _, endPnt := range reflect.ValueOf(APIRouteMap()).MapKeys() {
		log.Println("GENERATING END POINT: ", endPnt)
		http.HandleFunc(endPnt.String(), generateAPIEndPoint(handler))
	}
}
