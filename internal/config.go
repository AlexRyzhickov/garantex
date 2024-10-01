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
	pgUser      string
	pgPass      string
	pgDb        string
	pgHost      string
	pgPort      string
}

func (c *Config) Read() {
	flag.IntVar(&c.Port, "port", 8080, "grpc gateway port")
	flag.IntVar(&c.GRPCPort, "grps-port", 9090, "grpc port")
	flag.IntVar(&c.MetricsPort, "metrics-port", 9092, "metrics port")

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DATABASE")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	flag.StringVar(&c.pgUser, "pg-user", user, "")
	flag.StringVar(&c.pgPass, "pg-pass", pass, "")
	flag.StringVar(&c.pgDb, "pg-db", db, "")
	flag.StringVar(&c.pgHost, "pg-host", host, "")
	flag.StringVar(&c.pgPort, "pg-port", port, "")
	flag.Parse()

	c.Conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.pgUser, c.pgPass, c.pgHost, c.pgPort, c.pgDb)
}

func (c *Config) Print() {
	log.Println("========== GARANTEX PROXY CONFIG ==========")
	log.Println("GATEWAY_PORT............", c.Port)
	log.Println("GRPC_PORT...............", c.GRPCPort)
	log.Println("METRICS_PORT............", c.MetricsPort)
	log.Println("==================================")
}
