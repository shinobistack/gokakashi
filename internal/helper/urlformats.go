package helper

import (
	"log"
	"net/url"
	"strings"
)

func NormalizeServer(server string) string {
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server
	}
	return server
}

func ConstructURL(server string, path string) string {
	base := NormalizeServer(server)
	u, err := url.Parse(base)
	if err != nil {
		log.Fatalf("Invalid server URL: %s", base)
	}
	u.Path = path
	return u.String()
}
