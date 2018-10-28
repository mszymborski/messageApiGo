package main

import (
	"Utils"
	"github.com/gorilla/mux"
	"log"
	"message"
	"net/http"
)

type Route struct {
	Name string
	Path string
	HandlerFunc http.HandlerFunc
	HttpMethod string
}

var routes = []Route {
	Route{"Message", "/api/message", message.CreateMessage, "POST"},
	Route{"Send", "/api/send", message.SendMessage, "POST"},
	Route{"GetByEmailAddress", "/api/messages/{emailValue}", message.GetEmailsByAddress, "GET"},
}

func handleRequest()  {

	port := Utils.GetConfig(Utils.WEBSERVICE_PORT_KEY)

	router := mux.NewRouter()
	for _, route := range routes {
		router.Methods(route.HttpMethod).Path(route.Path).Name(route.Name).Handler(route.HandlerFunc)
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}


