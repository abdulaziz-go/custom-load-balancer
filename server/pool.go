package server

import "sync"

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

func (p *Pool) GetServerByNextRoundRubinAlgorithm() *Server {
	p.mux.Lock()
	defer p.mux.Unlock()
	servers := p.GetAliveServers()

	if len(servers) == 0 {
		return nil
	}

	server := servers[p.current%len(servers)]
	p.current++
	return server
}

func (pool *Pool) GetLeastConnectionServerAlgorithm() *Server {
	pool.mux.RLock()
	defer pool.mux.RUnlock()

	var leastConnServer *Server

	if len(pool.servers) > 0 {
		minConnectedServer := pool.servers[0].Connections
		for _, server := range pool.servers {
			if server.GetConnections() < minConnectedServer {
				minConnectedServer = server.GetConnections()
				leastConnServer = server
			}
		}
	}
	return leastConnServer
}
