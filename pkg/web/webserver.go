package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/shinobistack/gokakashi/internal/config/v0"
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

// StartWebServers starts web servers based on the configuration
func (ws *WebServer) StartWebServers(cfg *config.Config) error {
	for websiteID, website := range cfg.Websites {
		if err := ws.validateConfig(website); err != nil {
			log.Printf("Invalid config for website %s: %v", websiteID, err)
			continue
		}

		// Ensure report directories are created
		reportDir := filepath.Join(ReportsRootDir, website.ReportSubDir)
		if err := os.MkdirAll(reportDir, os.ModePerm); err != nil {
			log.Printf("Failed to create directory %s: %v", reportDir, err)
			continue
		}
		log.Printf("Created directory for %s: %s", websiteID, reportDir)

		// Create handlers and servers for each website
		handler := ws.createWebsiteHandler(websiteID, website)
		server := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", website.Hostname, website.Port),
			Handler: handler,
		}
		ws.Servers[websiteID] = server

		// Start server in a goroutine
		go func(id string, s *http.Server) {
			log.Printf("Starting web server for %s on %s", id, s.Addr)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Web server for %s stopped: %v", id, err)
			}
		}(websiteID, server)
	}

	return nil
}

// createWebsiteHandler sets up the routes and handlers for a website
func (ws *WebServer) createWebsiteHandler(websiteID string, websiteConfig config.Website) http.Handler {
	mux := http.NewServeMux()
	// Set up routes for listing directories and viewing reports
	mux.HandleFunc("/reports/", ws.handleListDirectoriesAndFiles(websiteID, websiteConfig.ReportSubDir))
	mux.HandleFunc("/reports/view", ws.handleListReports(websiteConfig.ReportSubDir))
	mux.HandleFunc("/reports/view/file", ws.handleViewReportFile(websiteConfig.ReportSubDir))

	return mux
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

// validateConfig validates the configuration for each website
func (ws *WebServer) validateConfig(website config.Website) error {
	if website.Hostname == "" || website.Port == 0 || website.ReportSubDir == "" {
		return fmt.Errorf("hostname, port, or report_sub_dir is missing")
	}
	return nil
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
