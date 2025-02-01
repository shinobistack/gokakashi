package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

const ReportsRootDir = "reports/"

// WebServer struct stores a reference to all running web servers
type WebServer struct {
	Servers map[string]*http.Server
}

// NewWebServer creates a new WebServer instance
func NewWebServer() *WebServer {
	return &WebServer{
		Servers: make(map[string]*http.Server),
	}
}

// handleListDirectoriesAndFiles lists both directories and files inside the report_sub_dir
func (ws *WebServer) handleListDirectoriesAndFiles(websiteID, rootPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// List both directories and files inside the report_sub_dir
		reportDir := filepath.Join(ReportsRootDir, rootPath)

		// List both directories and files
		entries, err := ws.listDirectoriesAndFiles(reportDir)
		if err != nil {
			http.Error(w, "Failed to list directories and files", http.StatusInternalServerError)
			return
		}

		log.Printf("WebsiteID: %s - Entries found: %v", websiteID, entries)

		// Use the directoryTemplate to display both directories and files
		tmpl := template.Must(template.New("directories").Parse(directoryTemplate))
		err = tmpl.Execute(w, struct{ Entries []string }{Entries: entries})
		if err != nil {
			http.Error(w, "Failed to template site", http.StatusInternalServerError)
			return
		}
	}
}

// handleListReports lists the reports inside a selected directory for the given websiteID
func (ws *WebServer) handleListReports(rootPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("dir")
		if dir == "" {
			http.Error(w, "No directory specified", http.StatusBadRequest)
			return
		}

		reportDir := filepath.Join(ReportsRootDir, rootPath, dir)

		// Check if the dir is actually a file
		if fileInfo, err := os.Stat(reportDir); err == nil && !fileInfo.IsDir() {
			// If it's a file, serve it directly
			http.Redirect(w, r, fmt.Sprintf("/reports/view/file?file=%s", reportDir), http.StatusFound)
			return
		}

		files, err := filepath.Glob(filepath.Join(reportDir, "*_report.json"))
		if err != nil {
			http.Error(w, "Failed to load reports", http.StatusInternalServerError)
			return
		}

		log.Printf("Listing reports in directory: %s", reportDir)

		tmpl := template.Must(template.New("reports").Parse(reportTemplate))
		err = tmpl.Execute(w, struct {
			Reports []string
			Dir     string
		}{Reports: files, Dir: dir})
		if err != nil {
			http.Error(w, "Failed to template site", http.StatusInternalServerError)
			return
		}
	}
}

// handleViewReportFile serves a specific report file for the given websiteID
func (ws *WebServer) handleViewReportFile(rootPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("file")
		if file == "" {
			http.Error(w, "No file specified", http.StatusBadRequest)
			return
		}

		// If the file path contains the full path, just use it directly
		filePath := filepath.Clean(file)
		// Explore reading file in chunks and responding back to avoid
		// potential memory exhaustion in case of large reports
		reportData, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Failed to read report", http.StatusInternalServerError)
			return
		}

		log.Printf("Serving report file: %s", filePath)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(reportData)
		if err != nil {
			log.Println("Error writing report data in response", err)
			return
		}
	}
}

// listDirectoriesAndFiles lists both directories and files inside the given path
func (ws *WebServer) listDirectoriesAndFiles(rootPath string) ([]string, error) {
	var entries []string
	files, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// Add both directories and files to the list
		entries = append(entries, file.Name())
	}

	return entries, nil
}

// HTML templates for directory and report listing
const directoryTemplate = `
<!DOCTYPE html>
<html>
<head>
  <title>Scan Reports - Entries</title>
</head>
<body>
  <h1>Scan Report Entries</h1>
  <ul>
      {{range .Entries}}
      <li><a href="/reports/view?dir={{.}}">{{.}}</a></li>
      {{end}}
  </ul>
</body>
</html>
`

const reportTemplate = `
<!DOCTYPE html>
<html>
<head>
  <title>Private Scan Reports - {{.Dir}}</title>
</head>
<body>
  <h1>Reports in {{.Dir}}</h1>
  <ul>
      {{range .Reports}}
      <li><a href="/reports/view/file?file={{.}}">{{.}}</a></li>
      {{end}}
  </ul>
</body>
</html>
`
