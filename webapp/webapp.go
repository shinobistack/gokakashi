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

func ReactApp() (http.Handler, error) {
	reactApp, err := fs.Sub(WebAssets, "dist")
	if err != nil {
		return nil, fmt.Errorf("error finding the dist folder: %w", err)
	}

	return http.FileServerFS(reactApp), nil
}
