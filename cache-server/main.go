package main

import (
	"DistributedCache/cache-server/cache"
	"DistributedCache/cache-server/http"
	"DistributedCache/cache-server/tcp"
)

func main() {
	c := cache.New("inmemory")
	go tcp.New(c).Listen()
	http.New(c).Listen()
}
