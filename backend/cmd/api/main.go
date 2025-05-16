package main

import (
	"github.com/hendraadwii/library/internal/server"
)

func main() {
	// This is just a proxy to the actual server implementation
	// in server.go to allow better code organization
	server.Run()
}
