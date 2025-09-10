package healthcheck

import (
	"fmt"
	"load-balancer/server"
	"net"
	"time"
)

type HealthChecker struct {
	pool     *server.Pool
	interval time.Duration
	timeout  time.Duration
	done     chan bool
}

func NewHealthChecker(pool *server.Pool, interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		pool:     pool,
		interval: interval,
		timeout:  timeout,
		done:     make(chan bool),
	}
}

func (hc *HealthChecker) Start() {
	go hc.healthCheck()
}

func (hc *HealthChecker) healthCheck() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hc.checkAllServer()
		case <-hc.done:
			fmt.Print("Health Checker stopped")
			return
		}
	}
}

func (hc *HealthChecker) checkAllServer() {
	for _, srv := range hc.pool.GetServers() {
		go hc.checkServer(srv)
	}
}

func (hc *HealthChecker) checkServer(srv *server.Server) {
	address := fmt.Sprintf("%s:%s", srv.Url.Hostname(), srv.Url.Port())
	conn, err := net.DialTimeout("tcp", address, hc.timeout)
	if err != nil {
		if srv.IsAlive() {
			srv.SetAlive(false)
		}
		return
	}
	conn.Close()

	if !srv.IsAlive() {
		srv.SetAlive(true)
	}
}
