package server

import (
	"fmt"
	"syscall"

	"github.com/Prash766/go-redis/core"
)

func StartAsyncTCPServer() {
	fmt.Printf("Connected to Async Server")
	var maxConn = 20000
	var connectedClients = 0
	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("error creating socket: %v\n", err)
		return
	}
	//  addr := "127.0.0.1"
	// parts := strings.Split(addr, ".")
	// var ip [4]byte
	// for i, p := range parts {
	//     n, err := strconv.Atoi(p)
	//     if err != nil {
	//         panic(err)
	//     }
	//     ip[i] = byte(n)
	// }

	sockAddr := &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{0, 0, 0, 0},
	}
	bindErr := syscall.Bind(serverFd, sockAddr)
	if bindErr != nil {
		fmt.Println("Unable to bind the socket", bindErr)
		return
	}
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFd),
	}
	connectedClients++
	fmt.Println("server", serverFd)
	epollEvents := make([]syscall.EpollEvent, maxConn)
	epollFd, err := syscall.EpollCreate1(serverFd)
	if err != nil {
		fmt.Printf("error creating epoll: %v\n", err)
		return
	}
	syscall.Listen(serverFd, maxConn)
	syscall.SetNonblock(serverFd, true)
	syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, serverFd, &event)
	for {
		count, _ := syscall.EpollWait(epollFd, epollEvents, -1)
		for i := range count {
			if epollEvents[i].Fd == int32(serverFd) {
				clientFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					fmt.Println("Error connecting to a new client")
				}
				clientEvent := &syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(clientFd),
				}
				syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, clientFd, clientEvent)
			} else {
				command := core.FDComm{
					Fd: int(epollEvents[i].Fd),
				}
				cmd, err := readCommands(command)
				if err != nil {
					defer syscall.Close(command.Fd)
					connectedClients--
					fmt.Printf("Client Disconnected , remaining clients %d", connectedClients)
				}
				respond(cmd, command)
			}
		}

	}

}
