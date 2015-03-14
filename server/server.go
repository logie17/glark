package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"bufio"
)

type Server struct {
	ClientPort uint16
	ServerPort uint16

	DefaultLanguage string

	Hostname string

	Logger *log.Logger

	clientListener net.Listener
	serverListener net.Listener
}

func NewServer() *Server {
	s := new(Server)
	s.ClientPort = 9000
	s.ServerPort = 9001
	s.DefaultLanguage = "en"
	s.Hostname = "localhost"

	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		names, _ := net.LookupAddr(addr.String())
		if len(names) > 0 {
			s.Hostname = names[0]
			break
		}
	}

	return s
}

func (s *Server) Serve() (err error) {
	s.serverListener, err = net.Listen("tcp", fmt.Sprint("127.0.0.1:", s.ServerPort))
	if err != nil {
		s.log(err)
		return err
	}
	s.clientListener, err = net.Listen("tcp", fmt.Sprint("127.0.0.1:", s.ClientPort))
	if err != nil {
		s.log(err)
		return err
	}
	var wg sync.WaitGroup
	go s.handleConn(s.serverListener)
	wg.Add(1)
	go s.handleConn(s.clientListener)
	wg.Add(1)
	wg.Wait()
	return nil
}

func (s *Server) handleConn(l net.Listener) {
	ch := make(chan string)
	done := make(chan string)
	
CONNECTION:
	for {
		conn, err := l.Accept()
		if err != nil {
			s.log(err)
			return
		}

		go s.handleStream(conn, ch, done)
		println ("connection started")
		for {
			select {
			case line := <-ch:
				println(line)
			case end := <-done:
				println(end)
				break CONNECTION;
			}
		}
	}
}

func (s *Server) handleStream(conn net.Conn, ch chan string, done chan string) {
	defer close(ch)
	bufc := bufio.NewReader(conn)

	for {
		line, _, err := bufc.ReadLine()
		if err != nil {
			break
		}
		if string(line) == "neat" {
			done<-"Stream Closed"
		}
		ch<-string(line)
	}
}

func (s *Server) log(v interface{}) {
	if s.Logger != nil && v != nil {
		s.Logger.Println(v)
	}
}
