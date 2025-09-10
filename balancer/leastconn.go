package balancer

import (
	"load-balancer/server"
)

type LeastConnectionBalancer struct {
	pool *server.Pool
}

func (lc *LeastConnectionBalancer) GetNextServer() *server.Server {
	return lc.pool.GetLeastConnectionServerAlgorithm()
}

func (lc *LeastConnectionBalancer) SetServerPool(pool *server.Pool) {
	lc.pool = pool
}
