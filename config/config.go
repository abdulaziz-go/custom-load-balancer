package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	LoadBalancer   LoadBalancerConfig
	BackendServers []BackendServer
	HealthCheck    HealthCheckConfig
}

type LoadBalancerConfig struct {
	Port      int
	Algorithm string
}

type BackendServer struct {
	Address string
	Weight  int
}

type HealthCheckConfig struct {
	Interval int
	Timeout  int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	portStr := getOrDefault("LB_PORT", "8080")
	port, _ := strconv.Atoi(portStr)

	algorithm := getOrDefault("LB_ALGORITHM", "round_robin")

	intervalStr := getOrDefault("HC_INTERVAL", "10")
	interval, _ := strconv.Atoi(intervalStr)

	timeoutStr := getOrDefault("HC_TIMEOUT", "5")
	timeout, _ := strconv.Atoi(timeoutStr)

	var servers []BackendServer
	for i := 1; ; i++ {
		key := "SERVER" + strconv.Itoa(i) + "_ADDRESS"
		addr := os.Getenv(key)
		if addr == "" {
			break
		}
		servers = append(servers, BackendServer{
			Address: addr,
			Weight:  1,
		})
	}

	cfg := &Config{
		LoadBalancer: LoadBalancerConfig{
			Port:      port,
			Algorithm: algorithm,
		},
		BackendServers: servers,
		HealthCheck: HealthCheckConfig{
			Interval: interval,
			Timeout:  timeout,
		},
	}

	return cfg, nil
}

func getOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
