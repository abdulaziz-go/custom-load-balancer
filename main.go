package main

import (
	"fmt"
	"load-balancer/balancer"
	"load-balancer/config"
	"load-balancer/healthcheck"
	"load-balancer/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	pool := server.NewPool()

	for _, backendServer := range cfg.BackendServers {
		serverURL := backendServer.Address
		if !strings.HasPrefix(serverURL, "http://") {
			serverURL = "http://" + serverURL
		}
		srv, err := server.NewServerAdd(serverURL, backendServer.Weight)
		if err != nil {
			log.Printf("Failed to create server %s: %v", backendServer.Address, err)
			continue
		}
		pool.AddServer(srv)
		log.Printf("Added server: %s", backendServer.Address)
	}

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

func handleRequest(w http.ResponseWriter, r *http.Request, lb balancer.LoadBalancer) {
	fmt.Println("request received")
	nxtServer := lb.GetNextServer()
	if nxtServer == nil {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}
	nxtServer.AddConnection()
	defer nxtServer.RemoveConnection()

	fmt.Printf(" Request sent to backend: %s\n", nxtServer.Url)
	nxtServer.ReverseProxy.ServeHTTP(w, r)
}
