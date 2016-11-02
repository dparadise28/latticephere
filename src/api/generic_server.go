package main

import (
	"{service}"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func main() {left_curly_bracket}
	var srv http.Server
	http2.VerboseLogs = true
	srv.Addr = ":{port}"

	// This enables http2 support
	http2.ConfigureServer(&srv, nil)

	http.HandleFunc("/{service}", {service}.{method})
	// Listen as https ssl server
	// NOTE: WITHOUT SSL IT WONT WORK!!
	// To self generate a test ssl cert/key you could go to
	// http://www.selfsignedcertificate.com/
	// or read the openssl manual
	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
{right_curly_bracket}
