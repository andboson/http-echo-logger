package main

import (
	"log"
	"os"

	"andboson/http-echo-logger/server"
	"andboson/http-echo-logger/templates"
)

const customEchoEndpointEnv = "CUSTOM_ENDPOINT"

func main() {
	// check templates
	tpls, err := templates.NewTemplates()
	if err != nil {
		log.Fatalf("template load error: %+v", err)
	}

	// start the server
	serv := server.NewServer(server.DefaultHTTPAddr, tpls, os.Getenv(customEchoEndpointEnv))

	if err := serv.Start(); err != nil {
		log.Printf("unable to start the server: %+v", err)
	}
}
