package main

import (
	"flag"
	"fmt"

	"github.com/Prash766/go-redis/server"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	host := flag.String("host", "0.0.0.0", "Host to listen on")
	fmt.Println("port", *port, *host)
	flag.Parse()
	server.StartAsyncTCPServer()
	fmt.Println("TCP server started!")
}
