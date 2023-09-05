package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/conf"
	"github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/loadbalancer"
)

func main() {
	conf, err := conf.Load("conf.yml")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	lb := loadbalancer.New(conf)
	if err := lb.Listen(context.Background()); err != nil {
		slog.Error(err.Error())
	}
}
