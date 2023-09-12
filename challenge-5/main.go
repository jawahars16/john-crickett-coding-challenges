package main

import "github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/server"

func main() {
	lb := server.New(
		server.Backend{Host: "localhost", Port: 4001},
		server.Backend{Host: "localhost", Port: 4002},
		server.Backend{Host: "localhost", Port: 4003})
	lb.Listen(4000)
}
