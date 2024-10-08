package web

//import (
//	"html/template"
//	"log"
//	"net/http"
//	"os"
//	"path/filepath"
//	"strconv"
//)
//
//// StartPrivateServer starts a private web server to list and serve reports.
//func StartPrivateServer(reportRootPath string, port int) {
//	mux := http.NewServeMux()
//
//	// Handle requests to list directories inside the /reports directory
//	mux.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
//		// Ensure the root reports directory exists
//		err := os.MkdirAll(reportRootPath, os.ModePerm)
//		if err != nil {
//			http.Error(w, "Failed to create or access the reports directory", http.StatusInternalServerError)
//			return
//		}
//
//		// List all subdirectories inside the rootPath
//		dirs, err := listDirectories(reportRootPath)
//		if err != nil {
//			http.Error(w, "Failed to load directories", http.StatusInternalServerError)
//			return
//		}
//
//		// Display directories using a template
//		tmpl := template.Must(template.New("directories").Parse(directoryTemplate))
//		data := struct {
//			Directories []string
//		}{
//			Directories: dirs,
//		}
//
//		tmpl.Execute(w, data)
//	})
//
//	// Handle requests to list report files inside a selected directory
//	mux.HandleFunc("/reports/view", func(w http.ResponseWriter, r *http.Request) {
//		dir := r.URL.Query().Get("dir")
//		if dir == "" {
//			http.Error(w, "No directory specified", http.StatusBadRequest)
//			return
//		}
//
//		// List all report files in the specified directory
//		files, err := filepath.Glob(filepath.Join(dir, "*_report.json"))
//		if err != nil {
//			http.Error(w, "Failed to load reports", http.StatusInternalServerError)
//			return
//		}
//
//		// Display the reports using a template
//		tmpl := template.Must(template.New("reports").Parse(reportTemplate))
//		data := struct {
//			Reports []string
//			Dir     string
//		}{
//			Reports: files,
//			Dir:     dir,
//		}
//
//		tmpl.Execute(w, data)
//	})
//
//	// Handle requests to view individual report files
//	mux.HandleFunc("/reports/view/file", func(w http.ResponseWriter, r *http.Request) {
//		reportFile := r.URL.Query().Get("file")
//		if reportFile == "" {
//			http.Error(w, "No file specified", http.StatusBadRequest)
//			return
//		}
//
//		// Read the report file and serve it
//		reportData, err := os.ReadFile(reportFile)
//		if err != nil {
//			http.Error(w, "Failed to read report", http.StatusInternalServerError)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.Write(reportData)
//	})
//
//	log.Printf("Starting private server on port %d...", port)
//	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
//		log.Fatalf("Private server failed: %v", err)
//	}
//}
//
//// listDirectories lists all subdirectories inside the reports root directory
//func listDirectories(rootPath string) ([]string, error) {
//	var dirs []string
//
//	// Read all items in the reports root directory
//	files, err := os.ReadDir(rootPath)
//	if err != nil {
//		return nil, err
//	}
//
//	// Filter for directories
//	for _, file := range files {
//		if file.IsDir() {
//			dirs = append(dirs, filepath.Join(rootPath, file.Name()))
//		}
//	}
//
//	return dirs, nil
//}
//
//const directoryTemplate = `
//<!DOCTYPE html>
//<html>
//<head>
//    <title>Private Scan Reports - Directories</title>
//</head>
//<body>
//    <h1>Private Scan Report Directories</h1>
//    <ul>
//        {{range .Directories}}
//        <li><a href="/reports/view?dir={{.}}">{{.}}</a></li>
//        {{end}}
//    </ul>
//</body>
//</html>
//`
//
//const reportTemplate = `
//<!DOCTYPE html>
//<html>
//<head>
//    <title>Private Scan Reports - {{.Dir}}</title>
//</head>
//<body>
//    <h1>Reports in {{.Dir}}</h1>
//    <ul>
//        {{range .Reports}}
//        <li><a href="/reports/view/file?file={{.}}">{{.}}</a></li>
//        {{end}}
//    </ul>
//</body>
//</html>
//`
