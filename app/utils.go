package main

import (
	"fmt"
	"strings"
)

func parseRequestLine(request []byte) (method, path string, proto string, err error) {
	firstLine := strings.SplitN(string(request), "\r\n", 2)[0]

	parts := strings.SplitN(firstLine, " ", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid request line")
	}

	method = parts[0]
	if !isValidMethod(method) {
		return "", "", "", fmt.Errorf("invalid HTTP method")
	}

	path = parts[1]
	proto = parts[2]

	return method, path, proto, nil
}

func isValidMethod(method string) bool {
	validMethods := map[string]bool{
		"GET":  true,
		"POST": true,
		// "PUT":     true,
		// "DELETE":  true,
		// "HEAD":    true,
		// "OPTIONS": true,
		// "PATCH":   true,
	}
	return validMethods[method]
}
