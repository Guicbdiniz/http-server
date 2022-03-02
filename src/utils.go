package httpserver

import "fmt"

// Create HTTP response initial line using specific status code and message.
func GetHttpResponse(code string, message string) string {
	return fmt.Sprintf("HTTP/1.1 %s %s", code, message)
}

// Check if HTTP request method is among the valid ones.
func HttpRequestMethodIsValid(method string) bool {
	for _, validMethod := range getValidHttpMethods() {
		if method == validMethod {
			return true
		}
	}
	return false
}

func getValidHttpMethods() []string {
	return []string{
		"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
	}
}
