package web

import (
	"html/template"
	"log"
	"net/http"
)

// StartPrivateServer starts the private web server on the specified port
func StartPrivateServer(report string, port string) {
	http.HandleFunc("/private/report", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("private").Parse(privateTemplate))
		data := struct {
			Report string
		}{
			Report: report,
		}
		tmpl.Execute(w, data)
	})

	log.Printf("Starting private server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Private server failed: %v", err)
	}
}

const privateTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Private Scan Report</title>
</head>
<body>
    <h1>Private Scan Report</h1>
    <pre>{{.Report}}</pre>
</body>
</html>
`
