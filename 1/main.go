package main

import (
	"DistributedCache/1/HTTP"
	"DistributedCache/1/cache"
)

func main() {
	c := cache.New("inmemory")
	HTTP.New(c).Listen()
}
