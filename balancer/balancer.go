package balancer

import "load-balancer/server"

type LoadBalancer interface {
	GetNextServer() *server.Server
	SetServerPool(*server.Pool)
}

type Algorithm string

const (
	RoundRubin      Algorithm = "round_rubin"
	LeastConnection Algorithm = "least_connection"
	Weight          Algorithm = "weight"
)

func NewLoadBalancer(algorithm Algorithm, pool *server.Pool) LoadBalancer {
	switch algorithm {
	case RoundRubin:
		return &RoundRubinBalancer{pool: pool}
	case LeastConnection:
		return &LeastConnectionBalancer{pool: pool}
	case Weight:
	default:
		return &RoundRubinBalancer{pool: pool}
	}
}
