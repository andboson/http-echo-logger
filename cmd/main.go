package main

import (
	"log"

	"andboson/http-echo-logger/server"
	"andboson/http-echo-logger/templates"
)

func main() {
	// check templates
	tpls, err := templates.NewTemplates()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// start the server
	serv := server.NewServer(server.DefaultHTTPAddr, tpls)

	_ = serv.Start()
}
