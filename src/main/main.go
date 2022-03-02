package main

import (
	"Guicdbiniz/httpserver"
	"log"
)

func main() {
	server := httpserver.HttpServer{}

	indexRoute := httpserver.Route{
		Path:            "/",
		MethodToFileURI: make(map[string]string),
	}

	err := indexRoute.AddMethod("GET", "/home/guicbdiniz/Repositories/http-server/src/main/html/index.html")

	if err != nil {
		log.Println("Error captured while adding method to route:")
		log.Println(err)
	}

	server.SetUpRoute(indexRoute)

	server.Listen(":8080")
}
