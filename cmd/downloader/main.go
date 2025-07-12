package main

import (
	"downloader/internal/server"
	"fmt"
)

var Version = "0.1"

func main() {
	fmt.Printf("Version: %s\n", Version)

	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
