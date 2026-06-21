package server

import (
	"fmt"
	"syscall"

	"github.com/Prash766/go-redis/core"
)

func StartAsyncTCPServer() {
	fmt.Printf("Connected to Async Server")
	core.Init()
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
	err = syscall.SetsockoptInt(serverFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		fmt.Println("Error setting SO_REUSEADDR:", err)
		return
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
	epollFd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		fmt.Printf("error creating epoll: %v\n", err)
		return
	}
	syscall.Listen(serverFd, maxConn)
	syscall.SetNonblock(serverFd, true)
	syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, serverFd, &event)
	for {
		count, _ := syscall.EpollWait(epollFd, epollEvents, -1)
		for i := 0; i < count; i++ {
			if epollEvents[i].Fd == int32(serverFd) {
				clientFd, _, err := syscall.Accept(serverFd)
				syscall.SetNonblock(clientFd, true)
				if err != nil {
					syscall.Close(serverFd)
					fmt.Println("Error connecting to a new client")
					continue
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
					connectedClients--
					syscall.Close(command.Fd)
					fmt.Printf("Client Disconnected , remaining clients %d", connectedClients)
					continue
				}
				respond(cmd, command)
			}
		}

	}

}
