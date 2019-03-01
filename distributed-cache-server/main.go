package main

import (
	"DistributedCache/distributed-cache-server/cache"
	"DistributedCache/distributed-cache-server/cluster"
	"DistributedCache/distributed-cache-server/http"
	"DistributedCache/distributed-cache-server/tcp"
	"flag"
	"log"
)

func main() {
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Println("node is:", *node)
	log.Println("cluster is:", *clus)
	c := cache.New("inmemory")
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go tcp.New(c, n).Listen()
	log.Println("tcp listen")
	http.New(c, n).Listen()
	log.Println("http listen")
}
