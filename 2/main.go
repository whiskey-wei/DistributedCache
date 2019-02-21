package main

import (
	"DistributedCache/1/HTTP"
	"DistributedCache/1/cache"
	"DistributedCache/2/TCP"
)

func main() {
	ca := cache.New("inmemory")
	go TCP.New(ca).Listen()
	HTTP.New(ca).Listen()
}
