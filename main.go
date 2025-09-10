package main

import (
	"fmt"
	"load-balancer/balancer"
	"load-balancer/config"
	"load-balancer/healthcheck"
	"load-balancer/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	pool := server.NewPool()

	addServersToPool(pool, cfg.BackendServers)

	if len(pool.GetAliveServers()) == 0 {
		fmt.Printf("no server available")
		return
	}

	lb := balancer.NewLoadBalancer(balancer.Algorithm(cfg.LoadBalancer.Algorithm), pool)
	hc := healthcheck.NewHealthChecker(pool, time.Duration(cfg.HealthCheck.Interval)*time.Second, time.Duration(cfg.HealthCheck.Timeout)*time.Second)
	hc.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, lb)
	})

	addr := fmt.Sprintf(":%d", cfg.LoadBalancer.Port)
	fmt.Printf("Load balancer starting on %s with algorithm %s \n", addr, cfg.LoadBalancer.Algorithm)
	fmt.Printf("Backend server count %d", len(pool.GetServers()))

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Printf("failed %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	fmt.Println("server shutdown")
	hc.Stop()
}

func addServersToPool(pool *server.Pool, servers []config.BackendServer) {
	for _, backendServer := range servers {
		srv, err := server.NewServerAdd(fmt.Sprintf("%v", backendServer.Address), backendServer.Weight)
		if err != nil {
			fmt.Printf("not connected to server %v", err)
			continue
		}
		pool.AddServer(srv)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, lb balancer.LoadBalancer) {
	nxtServer := lb.GetNextServer()
	if nxtServer == nil {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	nxtServer.AddConnection()
	defer nxtServer.RemoveConnection()

	fmt.Println("request send to ", nxtServer.Url)
	nxtServer.ReverseProxy.ServeHTTP(w, r)
}
