package main

import "strings"

// type RequestHeader struct {
// 	host      string
// 	userAgent string
// 	accept    string
// }
// type Request struct {
// 	header *RequestHeader
// }

// func newHeader(host, userAgent, accept string) *RequestHeader {
// 	return &RequestHeader{
// 		host:      host,
// 		userAgent: userAgent,
// 		accept:    accept,
// 	}
// }

func parseAcceptEncoding(header string) []string {
	encodings := strings.Split(header, ",")
	var supported []string

	for _, encoding := range encodings {
		encoding = strings.TrimSpace(encoding)
		// Check for supported encodings
		if encoding == "gzip" || encoding == "br" || encoding == "deflate" {
			supported = append(supported, encoding)
		}
	}
	return supported
}
