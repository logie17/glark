package main

import (
	"github.com/logie17/glark/server"
)

func main() {
	s := server.NewServer()
	err := s.Serve()
	if err != nil {
		println(err)
	}

	// TODO, add signal handlers to stop server
}
