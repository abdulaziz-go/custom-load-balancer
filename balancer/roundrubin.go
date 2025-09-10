package balancer

import "load-balancer/server"

type RoundRubinBalancer struct {
	pool *server.Pool
}

func (rr *RoundRubinBalancer) GetNextServer() *server.Server {
	return rr.pool.GetNextRoundRobin()
}

func (rr *RoundRubinBalancer) SetServerPool(pool *server.Pool) {
	rr.pool = pool
}
