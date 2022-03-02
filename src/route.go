package httpserver

import (
	"errors"
	"log"
	"os"
)

type Route struct {
	Path            string
	MethodToFileURI map[string]string
}

// Add an allowed method to this route, connecting it to a file URI.
func (route Route) AddMethod(method string, fileURI string) error {
	if HttpRequestMethodIsValid(method) {
		route.MethodToFileURI[method] = fileURI
		return nil
	} else {
		return errors.New("invalid http method")
	}
}

// Get HTML file content from requested method.
func (route Route) GetResponseContent(requestedMethod string) (string, error) {
	for method, fileURI := range route.MethodToFileURI {
		if method == requestedMethod {
			content, error := os.ReadFile(fileURI)

			if error != nil {
				log.Printf("Error captured while getting content from file:\n")
				log.Println(error)
				break
			}

			return string(content), nil
		}
	}
	return "", errors.New("not allowed http method")
}
