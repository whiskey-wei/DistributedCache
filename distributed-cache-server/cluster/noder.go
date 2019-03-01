package cluster

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/memberlist"
	"stathat.com/c/consistent"
)

type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}

func New(addr string, cluster string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}

	circle := consistent.New()
	circle.NumberOfReplicas = 256
	go func() {
		for {
			//每隔一秒扫描集群中存在的节点，重新构建哈希环
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
			//log.Println(nodes)
		}
	}()
	log.Println("cluster join")
	return &node{circle, addr}, nil
}
