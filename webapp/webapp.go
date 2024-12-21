package webapp

import (
	"fmt"
	"net/http"
)

type Server struct {
	*http.Server
}

func New(addr string) *Server {
	routes := http.NewServeMux()

	routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "TODO")
	})

	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: routes,
		},
	}
}
