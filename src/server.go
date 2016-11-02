package main

import (
	//"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net/http"
	"networking"
)

func main() {
	var srv http.Server
	//http2.VerboseLogs = true
	srv.Addr = ":8000"
	resp, err := http.Get("http://myexternalip.com/raw")
	if err == nil {
		extip, extipErr := ioutil.ReadAll(resp.Body)
		if extipErr == nil {
			log.Println("Setting Server Address", string(extip[:len(extip)-1])+srv.Addr)
		} else {
			log.Println("\n\nTrouble Parsing external ip\n\nSetting Server Address", srv.Addr)
		}
	} else {
		// shouldnt stop the server from starting
		log.Println(err.Error())
		log.Println("\n\nTrouble Retreiving external ip\n\nSetting Server Address", srv.Addr)
	}
	resp.Body.Close()

	//log.Println("Enabling http2 support")
	//http2.ConfigureServer(&srv, nil)

	log.Println("\n\n-----Starting Endpoints\n")
	networking.ServeEndPoints()
	// Listen as https ssl server
	// NOTE: WITHOUT SSL IT WONT WORK!!
	// To self generate a test ssl cert/key you could go to
	// http://www.selfsignedcertificate.com/
	// or read the openssl manual
	//log.Println("Starting TLS")
	log.Fatal(srv.ListenAndServe())
	//log.Fatal(srv.ListenAndServeTLS("certs/server/cert.pem", "certs/server/key.pem"))
}
