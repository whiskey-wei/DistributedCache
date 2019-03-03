package cache

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
 接口具体实现
*/
type value struct {
	k       string
	v       []byte
	created time.Time
}
type inMemoryCache struct {
	c     map[string]*list.Element
	mem   *list.List
	mutex sync.RWMutex
	Stat
	ttl     time.Duration
	MemSize int64
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Set:", k, " ", v)
	tmp, exist := c.c[k]
	log.Println(exist)
	if exist {
		c.del(k, tmp.Value.(*value).v)
		c.mem.Remove(tmp)
	}
	c.c[k] = c.mem.PushFront(&value{k, v, time.Now()})
	//c.c[k] = &value{v, time.Now()}
	c.add(k, v)
	fmt.Println("LUR：", c.ValueSize, " ", c.MemSize)
	if c.ValueSize < c.MemSize {
		return nil
	}
	for c.ValueSize > c.MemSize {
		v := c.mem.Remove(c.mem.Back())
		delete(c.c, v.(*value).k)
		c.del(v.(*value).k, v.(*value).v)
	}
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Get:", k)
	tmp, exist := c.c[k]
	if !exist {
		return nil, nil
	}
	c.mem.MoveToFront(tmp)
	return tmp.Value.(*value).v, nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Println("Del:", k)
	tmp, exist := c.c[k]
	if exist {
		c.mem.Remove(tmp)
		delete(c.c, k)
		c.del(k, tmp.Value.(*value).v)
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func (c *inMemoryCache) NewScanner() Scanner {
	pairCh := make(chan *pair)
	closeCh := make(chan struct{})
	go func() {
		defer close(pairCh)
		c.mutex.RLock()
		for k, e := range c.c {
			c.mutex.RUnlock()
			select {
			case <-closeCh:
				return
			case pairCh <- &pair{k, e.Value.(value).v}:
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}()
	return &inMemoryScanner{pair{}, pairCh, closeCh}
}

func newInMemoryCache(ttl int, memsize int64) *inMemoryCache {
	c := &inMemoryCache{make(map[string]*list.Element), list.New(), sync.RWMutex{}, Stat{}, time.Duration(ttl) * time.Second, memsize}
	if ttl > 0 {
		go c.expirer()
	}
	return c
}

func (c *inMemoryCache) expirer() {
	for {
		time.Sleep(c.ttl)
		c.mutex.RLock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			if v.Value.(*value).created.Add(c.ttl).Before(time.Now()) {
				c.Del(k)
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}
}

func New(typ string, ttl int, memsize int64) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache(ttl, memsize)
	}
	if c == nil {
		panic("unknown cache type" + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
