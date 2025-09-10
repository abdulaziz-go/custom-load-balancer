package balancer

import (
	"load-balancer/server"
)

type LeastConnectionBalancer struct {
	pool *server.Pool
}

func (lc *LeastConnectionBalancer) GetNextServer() *server.Server {
	return lc.pool.GetLeastConnections()
}
