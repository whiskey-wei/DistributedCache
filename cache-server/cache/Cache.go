package cache

import (
	"log"
	"sync"
)

/*
 key-value的数目，
 key的总长度，
 value的总长度
*/
type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(k string, v []byte) {
	s.Count++
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count--
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

/*
 Cache接口
*/
type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}

/*
 接口具体实现
*/
type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, exist := c.c[k]
	if exist {
		c.del(k, tmp)
	}
	c.c[k] = v
	c.add(k, v)
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.c[k], nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v)
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}

func New(typ string) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache()
	}
	if c == nil {
		panic("unknown cache type" + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
