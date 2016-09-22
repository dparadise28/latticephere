package main

import (
	//"fmt"
	//"html"
	"log"
	"net/http"
	"networking"
	"tools"
	//"golang.org/x/net/http2"
)

func main() {
	var srv http.Server
	//http2.VerboseLogs = true
	srv.Addr = ":8000"

	// This enables http2 support
	//http2.ConfigureServer(&srv, nil)

	networking.ServerEndPoints()
	http.HandleFunc("/ggg", tools.RemodelJ)

	// Listen as https ssl server
	// NOTE: WITHOUT SSL IT WONT WORK!!
	// To self generate a test ssl cert/key you could go to
	// http://www.selfsignedcertificate.com/
	// or read the openssl manual
	log.Fatal(srv.ListenAndServe())
	//log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
