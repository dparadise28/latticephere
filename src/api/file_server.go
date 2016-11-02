package main

import (
	"file_server"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func main() {
	var srv http.Server
	http2.VerboseLogs = true
	srv.Addr = ":8000"

	// This enables http2 support
	http2.ConfigureServer(&srv, nil)

	http.HandleFunc("/", file_server.ManageFile)
	// Listen as https ssl server
	// NOTE: WITHOUT SSL IT WONT WORK!!
	// To self generate a test ssl cert/key you could go to
	// http://www.selfsignedcertificate.com/
	// or read the openssl manual
	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
