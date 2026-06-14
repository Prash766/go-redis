package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/Prash766/go-redis/core"
)

func readCommands(c io.ReadWriter) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf)
	fmt.Println("Client sent value ", string(buf[:n]))
	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}
	fmt.Println("tokens", tokens)
	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c io.ReadWriter) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
}

func StartTCPServer(port string, host string) {
	var connected_clients int = 0
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		panic("Error starting TCP server: " + err.Error())

	}
	fmt.Printf("TCP server started at %s:%s\n", host, port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		connected_clients++
		fmt.Println("New client connected")
		// buf := make([]byte, 1024)
		for {
			cmd, err := readCommands(conn)
			if err != nil {
				connected_clients--
				fmt.Println("Error reading from connection:", err)
				conn.Close()
				fmt.Printf("Client disconnected. Total connected clients: %d\n", connected_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			// fmt.Println("recieved buffer", buf[:n])
			// write_n, err := conn.Write(buf[:n])
			// if err != nil {
			// 	fmt.Println("Error writing to connection:", err)
			// 	break
			// }
			// fmt.Printf("Received %d bytes, echoed back %d bytes\n", n, write_n)
			respond(cmd, conn)
		}
	}
}
