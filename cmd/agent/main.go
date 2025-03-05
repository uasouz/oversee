package main

import (
	"fmt"
	"oversee/agent"
)

func main() {
	server := agent.NewIngestionAPI("localhost:4093")

	fmt.Println(server.Serve())
}
