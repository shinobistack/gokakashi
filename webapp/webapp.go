package webapp

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

type Server struct {
	*http.Server
}

//go:embed dist
var WebAssets embed.FS

func New(addr string) (*Server, error) {
	routes := http.NewServeMux()

	reactApp, err := fs.Sub(WebAssets, "dist")
	if err != nil {
		return nil, fmt.Errorf("error finding the dist folder: %w", err)
	}

	routes.Handle("/", http.FileServerFS(reactApp))

	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: routes,
		},
	}, nil
}
