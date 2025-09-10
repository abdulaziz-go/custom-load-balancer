package server

import (
	"math"
	"sync"
)

type Pool struct {
	servers []*Server
	current int
	mux     sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		servers: make([]*Server, 0),
		current: 0,
	}
}

func (pool *Pool) AddServer(server *Server) {
	pool.mux.Lock()
	defer pool.mux.Unlock()
	pool.servers = append(pool.servers, server)
}

func (pool *Pool) GetServers() []*Server {
	pool.mux.RLock()
	defer pool.mux.RUnlock()
	return pool.servers
}

func (pool *Pool) GetAliveServers() []*Server {
	pool.mux.RLock()
	defer pool.mux.RUnlock()

	var aliveServers []*Server
	for _, server := range pool.servers {
		if server.IsAlive() {
			aliveServers = append(aliveServers, server)
		}
	}
	return aliveServers
}

func (pool *Pool) GetNextRoundRobin() *Server {
	pool.mux.Lock()
	defer pool.mux.Unlock()

	aliveServers := make([]*Server, 0)
	for _, server := range pool.servers {
		if server.IsAlive() {
			aliveServers = append(aliveServers, server)
		}
	}

	if len(aliveServers) == 0 {
		return nil
	}

	server := aliveServers[pool.current%len(aliveServers)]
	pool.current++
	return server
}

func (pool *Pool) GetLeastConnections() *Server {
	pool.mux.Lock()
	defer pool.mux.Unlock()

	var leastConnServer *Server
	minConnections := math.MaxInt

	for _, server := range pool.servers {
		if server.IsAlive() && server.GetConnections() < minConnections {
			minConnections = server.GetConnections()
			leastConnServer = server
		}
	}

	return leastConnServer
}
