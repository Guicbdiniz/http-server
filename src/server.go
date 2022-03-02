package httpserver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// HTTP Server TODO DOC
type HttpServer struct {
	routes []Route
}

// Listen to any incoming requests in the passed address.
func (server HttpServer) Listen(address string) {
	log.Printf("Listenin to address %s...\n", address)

	listener, listeningError := net.Listen("tcp", address)
	if listeningError != nil {
		log.Printf("Error captured while listening in address %s\n", address)
		log.Fatalln(listeningError)
	}

	for {
		conn, connError := listener.Accept()
		if connError != nil {
			log.Printf("Error captured while accepting a new connection")
			log.Fatalln(connError)
		}
		go server.handleConnection(conn)
	}
}

// Set up a HTTP route to render a specific HTML file.
func (server *HttpServer) SetUpRoute(route Route) {
	server.routes = append(server.routes, route)
}

func (server HttpServer) handleConnection(conn net.Conn) {
	log.Println("New connection being handled...")

	scanner := bufio.NewScanner(conn)

	hasInitialLine := scanner.Scan()
	if !hasInitialLine {
		server.sendBadRequestResponse(conn)
		return
	}

	initialLine := scanner.Text()

	log.Printf("Connection initial line: %s\n", initialLine)

	initialLineWords := strings.Split(initialLine, " ")

	if len(initialLineWords) <= 2 {
		server.sendBadRequestResponse(conn)
		return
	}

	requestMethod := initialLineWords[0]
	requestedPath := initialLineWords[1]

	log.Printf("Request method: %s, request path: %s\n", requestMethod, requestedPath)

	if !HttpRequestMethodIsValid(requestMethod) {
		server.sendBadRequestResponse(conn)
		return
	}

	server.handleHttpRequest(conn, requestMethod, requestedPath)

}

func (server HttpServer) getRouteFromPath(path string) (Route, error) {
	for _, route := range server.routes {
		if route.Path == path {
			return route, nil
		}
	}

	return Route{}, errors.New("invalid path")
}

func (server HttpServer) handleHttpRequest(conn net.Conn, method string, path string) {
	selectedRoute, error := server.getRouteFromPath(path)

	if error != nil {
		server.sendNotFoundResponse(conn)
		return
	}

	htmlContent, error := selectedRoute.GetResponseContent(method)

	if error != nil {
		server.sendMethodNotAllowedResponse(conn)
		return
	}

	server.sendValidResponse(conn, htmlContent)
}

func (server HttpServer) sendNotFoundResponse(conn net.Conn) {
	log.Println("Not found response being send to connection")
	defer conn.Close()
	io.WriteString(conn, GetHttpResponse("404", "Not Found"))
}

func (server HttpServer) sendBadRequestResponse(conn net.Conn) {
	log.Println("Bad request response being send to connection")
	defer conn.Close()
	io.WriteString(conn, GetHttpResponse("400", "Bad Request"))
}

func (server HttpServer) sendMethodNotAllowedResponse(conn net.Conn) {
	log.Println("Method not allowed response being send to connection")
	defer conn.Close()
	io.WriteString(conn, GetHttpResponse("405", "Method Not Allowed"))
}

func (server HttpServer) sendValidResponse(conn net.Conn, content string) {
	log.Println("Valid response being send to connection")
	defer conn.Close()
	io.WriteString(conn, GetHttpResponse("200", "OK"))
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(content))
	fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, content)
}
