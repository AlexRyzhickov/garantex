package internal

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Conn        string
	Port        int
	GRPCPort    int
	MetricsPort int
}

func (c *Config) Read() {
	flag.IntVar(&c.Port, "http", 8080, "grpc gateway port")
	flag.IntVar(&c.GRPCPort, "grps-port", 9090, "grpc port")
	flag.IntVar(&c.MetricsPort, "metrics-port", 9092, "metrics port")
	flag.Parse()

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DATABASE")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	c.Conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, db)
}

func (c *Config) Print() {
	log.Println("========== GARANTEX PROXY CONFIG ==========")
	log.Println("GATEWAY_PORT............", c.Port)
	log.Println("GRPC_PORT...............", c.GRPCPort)
	log.Println("METRICS_PORT............", c.MetricsPort)
	log.Println("==================================")
}
