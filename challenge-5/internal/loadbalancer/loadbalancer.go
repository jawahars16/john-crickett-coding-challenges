package loadbalancer

import (
	"container/list"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/backend"
	"github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/conf"
	"golang.org/x/exp/slog"
)

var sem = make(chan struct{}, 1)

type LoadBalancer struct {
	conf              conf.Conf
	server            *http.Server
	backends          list.List
	unHealthyBackends list.List
}

func New(conf *conf.Conf) *LoadBalancer {
	lb := &LoadBalancer{
		conf: *conf,
	}
	if conf.Port <= 0 {
		panic("Loadbalancer configuration is invalid.")
	}
	server := newServer(conf.Port, lb.handle)
	if conf.Backends != nil {
		for _, bknd := range conf.Backends {
			lb.AddBackend(bknd)
		}
	}
	lb.server = server
	return lb
}

func (lb *LoadBalancer) Listen(ctx context.Context) error {
	if lb.server == nil {
		return fmt.Errorf("No server to listen")
	}
	slog.Info("Listening...", "addr", lb.server.Addr)
	go lb.performHealthCheck()
	return lb.server.ListenAndServe()
}

func (lb *LoadBalancer) performHealthCheck() {
	var wg sync.WaitGroup
	for {
		for b := lb.backends.Front(); b != nil; b = b.Next() {
			wg.Add(1)
			go func(b *list.Element) {
				bcknd := b.Value.(*backend.Backend)
				healthy := bcknd.CheckHealth(lb.conf.Health.URL)
				if !healthy {
					lb.backends.Remove(b)
					lb.unHealthyBackends.PushBack(bcknd)
				}
				wg.Done()
			}(b)
		}
		for b := lb.unHealthyBackends.Front(); b != nil; b = b.Next() {
			wg.Add(1)
			go func(b *list.Element) {
				bcknd := b.Value.(*backend.Backend)
				healthy := bcknd.CheckHealth(lb.conf.Health.URL)
				if healthy {
					lb.backends.PushBack(bcknd)
					lb.unHealthyBackends.Remove(b)
				}
				wg.Done()
			}(b)
		}
		wg.Wait()
		time.Sleep(time.Second * time.Duration(lb.conf.Health.Interval))
	}
}

func (lb *LoadBalancer) AddBackend(addr string) {
	b := backend.New(addr)
	lb.backends.PushBack(b)
	slog.Info("New backend added.", "addr", addr)
}

func (lb *LoadBalancer) Stop(ctx context.Context) {
	if lb.server == nil {
		panic("Cannot stop before starting the server")
	}

	lb.server.Shutdown(ctx)
}

func newServer(port int, handler func(http.ResponseWriter, *http.Request)) *http.Server {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serverMux,
	}
	return &server
}

func (lb *LoadBalancer) handle(res http.ResponseWriter, req *http.Request) {
	if lb.backends.Len() <= 0 {
		slog.Error("No backends available")
		res.WriteHeader(503)
		return
	}

	sem <- struct{}{}
	node := lb.backends.Front()
	backend := node.Value.(*backend.Backend)
	lb.backends.MoveToBack(node)
	<-sem

	response, err := backend.Handle(req)
	if err != nil {
		slog.Error(err.Error())
		res.WriteHeader(503)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		slog.Error(err.Error())
		res.WriteHeader(500)
		return
	}
	fmt.Println()
	fmt.Println("Recieved request from", req.RemoteAddr)
	fmt.Println(req.Method, "/", req.Proto)
	fmt.Println("Host:", req.Host)
	fmt.Println("User-Agent:", req.UserAgent())
	fmt.Println("Accept", req.Header.Get("Accept"))
	fmt.Println()
	fmt.Println("Response from server", response.Proto, response.StatusCode, response.Status)
	res.Write(data)
}

func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
