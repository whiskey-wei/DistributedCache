package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var typ, server, operation string
var total, valueSize, threads, keyspacelen, pipelen int

//指令集
func init() {
	flag.StringVar(&typ, "type", "redis", "cache server type")
	flag.StringVar(&server, "h", "localhost", "cache server address")
	flag.IntVar(&total, "n", 1000, "total number of request")
	flag.IntVar(&valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&threads, "c", 1, "number of parallel connections")
	flag.StringVar(&operation, "t", "set", "test set, could be get/set/mixed")
	flag.IntVar(&keyspacelen, "r", 0, "keyspacelen, use random keys from 0 to keyspacelen-1")
	flag.IntVar(&pipelen, "P", 1, "pipeline length")
	flag.Parse()
	fmt.Println("type is:", typ)
	fmt.Println("serveris:", server)
	fmt.Println("total", total, "request")
	fmt.Println("data size is:", valueSize)
	fmt.Println("we have ", threads, "connections")
	fmt.Println("operation is:", keyspacelen)
	fmt.Println("pipeline length is:", pipelen)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	ch := make(chan *result, threads)
	res := &result{0, 0, 0, make([]statistic, 0)}
	start := time.Now()
	for i := 0; i < threads; i++ {
		go oprate(i, total, ch)
	}

	//<-chan操作每次只会读出一个元素值，当threads个goroutine没有完成任务时
	//这里会就读不到threads个result值，就会阻塞等待goroutine向chan中写入result，
	//利用channal实现了goroutine的同步
	for i := 0; i < threads; i++ {
		res.addResult(<-ch)
	}

	d := time.Now().Sub(start)
	totalCount := res.getCount + res.missCount + res.setCount
	fmt.Printf("%d records get\n", res.getCount)
	fmt.Printf("%d records miss\n", res.missCount)
	fmt.Printf("%d records set\n", res.setCount)
	fmt.Printf("%f seconds total\n", d.Seconds())
	statCountSum := 0
	statTimeSum := time.Duration(0)
	for b, s := range res.statBuckets {
		if s.count == 0 {
			continue
		}
		statCountSum += s.count
		statTimeSum += s.time
		fmt.Printf("%d%% requests < %d ms\n", statCountSum*100/totalCount, b+1)
	}
	fmt.Printf("%d usec average for each request \n", int64(statTimeSum/time.Microsecond)/int64(statCountSum))
	fmt.Printf("throughput is %f MB/s\n", float64((res.getCount+res.setCount)*valueSize)/1e6/d.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(d.Seconds()))
}
