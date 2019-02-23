package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type result struct {
	v []byte
	e error
}

//使用类型为无缓冲的管道实现goroutine间的同步
//resultCh为有缓冲的管道，响应消息顺序放入此管道，之后的处理
//异步进行，最后返回客户端的结果也是顺序执行的
func (s *Server) process(conn net.Conn) {
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, 5000)
	defer close(resultCh)
	go reply(conn, resultCh)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error: ", e)
			}
			return
		}
		if op == 'S' {
			s.set(resultCh, r)
		} else if op == 'G' {
			s.get(resultCh, r)
		} else if op == 'D' {
			s.del(resultCh, r)
		} else {
			log.Println("close connetion due to invalid operation: ", op)
			return
		}
	}
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
	}
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
}
