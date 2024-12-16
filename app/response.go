package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
)

const (
	CONTENT_TEXT       string = "text/plain"
	CONTENT_FILE       string = "application/octet-stream"
	NOT_FOUND          string = "404 Not Found"
	METHOD_NOT_ALLOWED string = "405 Method Not Allowed"
	OK                 string = "200 OK"
	CREATED            string = "201 Created"
)

type ResponseHeader struct {
	httpVersion     string
	status          string
	contentType     string
	contentLength   string
	contentEncoding string
}
type Response struct {
	header *ResponseHeader
	body   string
	// file   []byte
}

const CRLF string = "\r\n"

func newResponseHeader(httpVersion, status, contentType string) *ResponseHeader {
	return &ResponseHeader{
		httpVersion: httpVersion,
		status:      status,
		contentType: contentType,
	}
}

func newResponse(header *ResponseHeader, body string) *Response {
	return &Response{
		header: header,
		body:   body,
	}
}

func (r *Response) generateResponse() []byte {
	var bodyBytes []byte
	if r.header.contentEncoding == "gzip" {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		gz.Write([]byte(r.body))
		gz.Close()
		bodyBytes = b.Bytes()
	} else {
		bodyBytes = []byte(r.body)
	}

	r.header.contentLength = strconv.Itoa(len(bodyBytes))

	var response bytes.Buffer
	response.WriteString(fmt.Sprintf("%s %s\r\n", r.header.httpVersion, r.header.status))
	response.WriteString(fmt.Sprintf("Content-Type: %s\r\n", r.header.contentType))
	response.WriteString(fmt.Sprintf("Content-Length: %s\r\n", r.header.contentLength))

	if r.header.contentEncoding != "" {
		response.WriteString(fmt.Sprintf("Content-Encoding: %s\r\n", r.header.contentEncoding))
	}

	response.WriteString("\r\n")
	response.Write(bodyBytes)

	return response.Bytes()
}
