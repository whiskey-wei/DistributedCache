package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"stathat.com/c/consistent"

	"github.com/stuarthu/go-implement-your-cache-server/cache-benchmark/cacheClient"
)

func main() {
	server := flag.String("h", "localhost", "cache server address")
	op := flag.String("c", "get", "command, could be get/set/del")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	flag.Parse()

	client := cacheClient.New("tcp", *server)
	cmd := &cacheClient.Cmd{Name: *op, Key: *key, Value: *value, Error: nil}
	client.Run(cmd)
	if cmd.Error != nil {
		fmt.Println("error:", cmd.Error)
	} else {
		fmt.Println(cmd.Value)
	}
}

func getServer(server string) (string, error) {
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://"+server+":12345/cluster", strings.NewReader(""))
	if err != nil {
		return "", err
	}
	reps, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	var memberlist []string
	dec := json.NewDecoder(reps.Body)
	for {
		var member string
		if err = dec.Decode(&member); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		fmt.Println(member)
		memberlist = append(memberlist, member)
	}
	circle := consistent.New()
	circle.NumberOfReplicas = 256
	circle.Set(memberlist)
	return server, err
}
