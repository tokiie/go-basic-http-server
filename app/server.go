package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

type Server struct {
	addr   string
	port   string
	routes map[string][]Route
	// handler Handler
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.StringVar(&directory, "directory", "tmp/", "File Directory Path")
	flag.Parse()

	server := newServer("0.0.0.0", "4221")

	server.AddRoute("/", "GET", handleGetBase)
	server.AddRoute("/echo", "GET", handleGetEcho)
	server.AddRoute("/user-agent", "GET", handleGetUserAgent)
	server.AddRoute("/files", "GET", handleGetFiles)
	server.AddRoute("/files", "POST", handlePostFiles)
	// 	body = string(data)

	server.listenAndServe()
}

func newServer(addr, port string) *Server {
	return &Server{
		addr:   addr,
		port:   port,
		routes: make(map[string][]Route),
	}
}

func (s *Server) getAddress() string {
	return s.addr + ":" + s.port
}

func (s *Server) AddRoute(path, method string, handler func(map[string]string, []byte) (string, string, string)) {
	if _, exists := s.routes[path]; !exists {
		s.routes[path] = []Route{}
	}
	s.routes[path] = append(s.routes[path], Route{
		Method:  method,
		Handler: handler,
	})
}

func (s *Server) listenAndServe() {
	l, err := net.Listen("tcp", s.getAddress())
	if err != nil {
		fmt.Println("Failed to bind to port 4221, ", err.Error())
		os.Exit(1)
	}

	for {

		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()

	// Read the initial request data
	requestData := make([]byte, 1024)
	n, err := c.Read(requestData)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// Split request into headers and body sections
	requestParts := strings.Split(string(requestData[:n]), "\r\n\r\n")
	if len(requestParts) != 2 {
		fmt.Println("Invalid request format")
		return
	}

	headerSection := requestParts[0]
	body := []byte(requestParts[1])

	// Parse the request line
	headerLines := strings.Split(headerSection, "\r\n")
	method, rawPath, proto, err := parseRequestLine([]byte(headerLines[0]))
	if err != nil {
		fmt.Println("Error parsing request line:", err)
		return
	}

	// Determine path
	var path string
	if strings.HasPrefix(rawPath, "/files") {
		path = "/files"
	} else if strings.HasPrefix(rawPath, "/echo") {
		path = "/echo"
	} else {
		path = rawPath
	}

	// Parse headers
	headers := make(map[string]string)
	headers["Original-Path"] = rawPath
	for i := 1; i < len(headerLines); i++ {
		line := strings.TrimSpace(headerLines[i])
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	responseHeader := newResponseHeader(proto, OK, TEXT_PLAIN)

	if headers["Accept-Encoding"] != "" {
		supportedEncoding := parseAcceptEncoding(headers["Accept-Encoding"])

		if len(supportedEncoding) > 0 {
			for _, encoding := range supportedEncoding {
				if encoding == "gzip" {
					responseHeader.contentEncoding = encoding
					break
				}
			}
		}
	}

	// Find matching route
	routes, pathExists := s.routes[path]
	if !pathExists {
		responseHeader.status = NOT_FOUND
		response := newResponse(
			responseHeader,
			"Path not found",
		)
		c.Write(response.generateResponse())
		return
	}

	// Find matching method handler
	var matchedRoute *Route
	for _, route := range routes {
		if route.Method == method {
			matchedRoute = &route
			break
		}
	}

	if matchedRoute == nil {
		responseHeader.status = METHOD_NOT_ALLOWED
		response := newResponse(
			responseHeader,
			"Method not allowed",
		)
		c.Write(response.generateResponse())
		return
	}

	// Execute handler with body
	responseBody, status, contentType := matchedRoute.Handler(headers, body)

	responseHeader.status = status
	responseHeader.contentType = contentType

	response := newResponse(
		responseHeader,
		responseBody,
	)
	c.Write(response.generateResponse())
}
