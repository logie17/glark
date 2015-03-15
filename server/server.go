package server

import (
	"fmt"
	"log"
	"net"
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
	s.handleConn(s.serverListener)
	s.handleConn(s.clientListener)

	defer func() {
		s.serverListener.Close()
		s.clientListener.Close()
	}()

	return nil
}

func (s *Server) handleConn(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			s.log(err)
			return
		}
		s.log("Connection started")
		go s.handleStream(conn)
	}

}

func (s *Server) handleStream(conn net.Conn) {
	bufc := bufio.NewReader(conn)

	for {
		line, _, err := bufc.ReadLine()
		if err != nil {
			break
		}
		if string(line) == "close" {
			s.log("Closing connection")
			conn.Close()
		}
		s.log(string(line))
	}
}

func (s *Server) log(v interface{}) {
	if s.Logger != nil && v != nil {
		s.Logger.Println(v)
	}
}
