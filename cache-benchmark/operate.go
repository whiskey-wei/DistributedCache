package main

import (
	"DistributedCache/cache-benchmark/cacheClient"
	"fmt"
	"math/rand"
	"strings"
)

func oprate(id, count int, ch chan *result) {
	client := cacheClient.New(typ, server)
	cmds := make([]*cacheClient.Cmd, 0)
	valuePrefix := strings.Repeat("a", valueSize) //构建用于测试的value
	r := &result{0, 0, 0, make([]statistic, 0)}
	for i := 0; i < count; i++ {
		var tmp int

		//构建测试的key值
		if keyspacelen > 0 {
			tmp = rand.Intn(keyspacelen)
		} else {
			tmp = id*count + i
		}

		key := fmt.Sprintf("%d", tmp)
		value := fmt.Sprintf("%s%d", valuePrefix, tmp)
		name := operation
		if operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}
		c := &cacheClient.Cmd{Name: name, Key: key, Value: value, Error: nil}

		if pipelen > 1 { //如果管线数>1先缓存这条命令
			cmds = append(cmds, c)
			if len(cmds) == pipelen { //管线数跟缓存的命令数目一样时，发送命令
				pipeline(client, cmds, r)
				cmds = make([]*cacheClient.Cmd, 0)
			}
		} else {
			run(client, c, r)
		}
	}
	ch <- r
}
