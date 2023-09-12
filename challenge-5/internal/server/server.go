package server

import (
	"container/list"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/httpparser"
	"golang.org/x/exp/slog"
)

var semaphore = make(chan struct{}, 1)

type Backend struct {
	Host string
	Port int
}

type Server struct {
	Backends          *list.List
	UnhealthyBackends *list.List
}

func New(backends ...Backend) *Server {
	bknds := list.New()
	uBknds := list.New()
	for _, backend := range backends {
		bknds.PushBack(backend)
	}
	return &Server{
		Backends:          bknds,
		UnhealthyBackends: uBknds,
	}
}

func (server *Server) Listen(port int) {
	addr := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Println("Listening on", addr)
	defer listener.Close()

	go server.performHealthCheck()

	for {
		incomingConnection, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("Panic:", err)
				}
			}()
			defer incomingConnection.Close()

			semaphore <- struct{}{} // acquire lock
			backendNode := server.Backends.Front()
			backend := backendNode.Value.(Backend)
			fmt.Println("Using the backend", backend.Host, backend.Port)
			server.Backends.MoveToBack(backendNode)
			<-semaphore // release lock

			backendConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
			if err != nil {
				slog.Error(err.Error())
			}
			defer backendConnection.Close()

			if err = handle(incomingConnection, backendConnection); err != nil {
				slog.Error(err.Error())
			}
		}()
	}
}

func handle(incomingConnection io.ReadWriter, backendConnection io.ReadWriter) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic:", err)
		}
	}()

	fmt.Println("Request")
	incomingData, err := httpparser.ReadBytes(incomingConnection)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	_, err = backendConnection.Write(incomingData)
	if err != nil {
		slog.Error(err.Error())
	}

	if tcpConnection, ok := backendConnection.(*net.TCPConn); ok {
		tcpConnection.CloseWrite()
	}

	fmt.Println("Response")
	backendResponse, err := httpparser.ReadBytes(backendConnection)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	incomingConnection.Write(backendResponse)
	return nil
}

func (lb *Server) performHealthCheck() {
	go cycle(lb.Backends, func(e *list.Element) {
		backend := e.Value.(Backend)
		_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
		if err != nil {
			fmt.Println("Health check failed for", backend.Host, backend.Port)
			lb.Backends.Remove(e)
			lb.UnhealthyBackends.PushBack(backend)
		}
		time.Sleep(time.Millisecond * 100)
	})
	go cycle(lb.UnhealthyBackends, func(e *list.Element) {
		backend := e.Value.(Backend)
		_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backend.Host, backend.Port))
		if err == nil {
			fmt.Println("Health check passed for", backend.Host, backend.Port)
			lb.UnhealthyBackends.Remove(e)
			lb.Backends.PushBack(backend)
		}
		time.Sleep(time.Millisecond * 100)
	})
}

func cycle(list *list.List, action func(*list.Element)) {
	node := list.Front()
	for {
		if node == nil {
			node = list.Front()
			continue
		}
		action(node)
		if node = node.Next(); node != nil {
			continue
		}

		node = list.Front()
	}
}
