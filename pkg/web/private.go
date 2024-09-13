package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func StartPrivateServer(reportPath string, port int) {
	// Create a new ServeMux for the private server
	mux := http.NewServeMux()

	mux.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		files, err := filepath.Glob(reportPath + "/*_report.json")
		if err != nil {
			http.Error(w, "Failed to load reports", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.New("private").Parse(privateTemplate))
		data := struct {
			Reports []string
		}{
			Reports: files,
		}
		tmpl.Execute(w, data)
	})

	mux.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Starting private server on port %d...", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		log.Fatalf("Private server failed: %v", err)
	}
}

const privateTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Private Scan Reports</title>
</head>
<body>
    <h1>Private Scan Reports</h1>
    <ul>
        {{range .Reports}}
        <li><a href="/view?file={{.}}">{{.}}</a></li>
        {{end}}
    </ul>
</body>
</html>
`
