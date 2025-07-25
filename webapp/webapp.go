package webapp

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type Server struct {
	*http.Server
}

//go:embed dist
var WebAssets embed.FS

func serveIndexHTML(w http.ResponseWriter, reactApp fs.FS) {
	indexFile, indexErr := reactApp.Open("index.html")
	if indexErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := []byte("500 - Could not find index.html")
		if _, err := w.Write(errMsg); err != nil {
			log.Println("error writing response:", err)
		}
		return
	}
	defer indexFile.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, indexFile)
}

func ReactApp() (http.Handler, error) {
	reactApp, err := fs.Sub(WebAssets, "dist")
	if err != nil {
		return nil, fmt.Errorf("error finding the dist folder: %w", err)
	}

	fileServer := http.FileServerFS(reactApp)

	customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If root path, serve index.html
		requestedPath := r.URL.Path[1:] // strip leading '/'
		if requestedPath == "" {
			serveIndexHTML(w, reactApp)
			return
		}

		// Try to open the requested file
		file, err := reactApp.Open(requestedPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Serve index.html for SPA routing
				serveIndexHTML(w, reactApp)
				return
			}
			// Other errors
			w.WriteHeader(http.StatusInternalServerError)
			errMsg := []byte("500 - Internal Server Error")
			if _, err := w.Write(errMsg); err != nil {
				log.Println("error writing response:", err)
			}
			return
		}
		file.Close()

		// File exists, serve it
		fileServer.ServeHTTP(w, r)
	})

	return customHandler, nil
}
