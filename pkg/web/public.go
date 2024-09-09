package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func StartPublicServer(reportPath string, port int) {
	http.HandleFunc("/public/report", func(w http.ResponseWriter, r *http.Request) {
		files, err := filepath.Glob(reportPath + "/*_report.json")
		if err != nil {
			http.Error(w, "Failed to load reports", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.New("public").Parse(publicTemplate))
		data := struct {
			Reports []string
		}{
			Reports: files,
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/public/view", func(w http.ResponseWriter, r *http.Request) {
		reportFile := r.URL.Query().Get("file")
		if reportFile == "" {
			http.Error(w, "No file specified", http.StatusBadRequest)
			return
		}

		reportData, err := os.ReadFile(reportFile)
		if err != nil {
			http.Error(w, "Failed to read report", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(reportData)
	})

	log.Printf("Starting public server on port %d...", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil { // Convert port to string
		log.Fatalf("Public server failed: %v", err)
	}
}

const publicTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Public Scan Reports</title>
</head>
<body>
    <h1>public Scan Reports</h1>
    <ul>
        {{range .Reports}}
        <li><a href="/public/view?file={{.}}">{{.}}</a></li>
        {{end}}
    </ul>
</body>
</html>
`
