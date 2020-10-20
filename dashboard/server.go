package dashboard

import (
	"fmt"
	"net/http"

	"github.com/c-bata/goptuna"
)

func NewServer(storage goptuna.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

type Server struct {
	storage goptuna.Storage
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
