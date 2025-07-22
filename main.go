package main

import (
	"log"
	"redis-clone/internal/server"
)

func main() {
	err := server.Start(":6380")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
