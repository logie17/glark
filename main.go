import (
	"server"
)

func main() {
	s := Sever()
	err := s.Serve()
	if err != nil {
		println(err)
	}

	// TODO, add signal handlers to stop server
}
