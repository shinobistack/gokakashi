package web

import (
	"html/template"
	"log"
	"net/http"
)

// StartPublicServer starts the public web server on the specified port
func StartPublicServer(report string, port string) {
	http.HandleFunc("/public/report", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("public").Parse(publicTemplate))
		data := struct {
			Report string
		}{
			Report: report,
		}
		tmpl.Execute(w, data)
	})

	log.Printf("Starting public server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Public server failed: %v", err)
	}
}

const publicTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Public Scan Report</title>
</head>
<body>
    <h1>Public Scan Report</h1>
    <pre>{{.Report}}</pre>
</body>
</html>
`
