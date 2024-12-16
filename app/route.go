package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Route struct {
	Method string
	// Route   string
	Handler func(headers map[string]string, body []byte) (string, string, string) // returns (body, status, contentType)
}

type Handler interface {
	ServeTCP(net.Conn)
}

type HandlerFunc func()

// case strings.Split(path, "/")[1] == "":
// 	break
// case strings.Split(path, "/")[1] == "echo":
// 	body = strings.Split(path, "/")[2]
// case strings.Split(path, "/")[1] == "user-agent":
// 	body = headers["User-Agent"]
// case strings.Split(path, "/")[1] == "files":
// 	fmt.Println("directory: ", directory)
// 	header.contentType = "application/octet-stream"
// 	fileName := strings.TrimPrefix(path, "/files/")
// 	data, err := os.ReadFile(filepath.Join(directory, fileName))
// 	if err != nil {
// 		header.status = "404 Not Found"
// 		fmt.Println("file not found: ", err.Error())
// 		break
// 	}

const (
	TEXT_PLAIN = "text/plain"
	FILE       = "application/octet-stream"
)

// TODO: base handler
func handleGetBase(headers map[string]string, body []byte) (string, string, string) {
	return "Welcome!", "200 OK", TEXT_PLAIN
}

// TODO: echo handler
func handleGetEcho(headers map[string]string, body []byte) (string, string, string) {
	echo := strings.TrimPrefix(headers["Original-Path"], "/echo/")
	return echo, "200 OK", TEXT_PLAIN
}

// TODO: user-agent handler
func handleGetUserAgent(headers map[string]string, body []byte) (string, string, string) {
	return headers["User-Agent"], "200 OK", TEXT_PLAIN
}

// TODO: files handler
func handleGetFiles(headers map[string]string, body []byte) (string, string, string) {
	fileName := strings.TrimPrefix(headers["Original-Path"], "/files/")
	data, err := os.ReadFile(filepath.Join(directory, fileName))

	if err != nil {
		fmt.Println("Error while opening the file: ", err)
		return handle404()
	}
	return string(data), "200 OK", FILE
}

// TODO: post handler
func handlePostFiles(headers map[string]string, body []byte) (string, string, string) {
	fileName := strings.TrimPrefix(headers["Original-Path"], "/files/")
	fmt.Println(body)
	err := os.WriteFile(filepath.Join(directory, fileName), body, 0644)
	if err != nil {
		return "Failed to save file", "500 Internal Server Error", TEXT_PLAIN
	}

	return "File saved successfully", "201 Created", FILE
}

func handle404() (string, string, string) {
	return "", "404 Not Found", TEXT_PLAIN
}
