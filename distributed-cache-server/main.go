package main

import (
	"DistributedCache/distributed-cache-server/cache"
	"DistributedCache/distributed-cache-server/cluster"
	"DistributedCache/distributed-cache-server/http"
	"DistributedCache/distributed-cache-server/tcp"
	"flag"
	"fmt"
	"log"
)

func main() {
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	ttl := flag.Int("ttl", 30, "cache time to live")
	Memory := flag.Int64("Memory", 50, "LUR memory")
	flag.Parse()
	log.Println("node is:", *node)
	log.Println("cluster is:", *clus)
	log.Println("ttl is: ", *ttl)
	fmt.Println("LUR memory is: ", *Memory)
	c := cache.New("inmemory", *ttl, *Memory)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go tcp.New(c, n).Listen()
	log.Println("tcp listen")
	http.New(c, n).Listen()
	log.Println("http listen")
}
