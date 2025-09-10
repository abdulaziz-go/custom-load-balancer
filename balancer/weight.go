package balancer

import (
	"load-balancer/server"
	"sync"
)

type WeightBalancer struct {
	pool         *server.Pool
	currentIndex int
	mux          sync.Mutex
}

func (wg *WeightBalancer) GetNextServer() *server.Server {
	var (
		weightedServers []*server.Server
	)

	wg.mux.Lock()
	defer wg.mux.Unlock()

	servers := wg.pool.GetAliveServers()
	if len(servers) == 0 {
		return nil
	}

	for _, srv := range servers {
		for i := 0; i < srv.Weight; i++ {
			weightedServers = append(weightedServers, srv)
		}
	}

	if len(weightedServers) == 0 {
		return nil
	}
	wg.currentIndex++
	return weightedServers[wg.currentIndex%len(weightedServers)]
}
