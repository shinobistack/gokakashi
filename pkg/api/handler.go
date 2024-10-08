package api

//
//import (
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//	"time"
//)
//
//// Define a struct for scan jobs
//type scanJob struct {
//	scanID   string
//	image    string
//	severity string
//	filePath string
//}
//
//// Job queue to process scan jobs
//var scanJobQueue = make(chan scanJob, 100) // Buffer size of 100 to handle burst loads
//
//// Number of workers for the pool
//var maxWorkers = 10
//
//// Worker function to process scan jobs
//func runWorker(id int, jobQueue <-chan scanJob) {
//	for job := range jobQueue {
//		log.Printf("Worker %d started processing scan: %s", id, job.scanID)
//		runScan(job.scanID, job.image, job.severity, job.filePath)
//		log.Printf("Worker %d finished processing scan: %s", id, job.scanID)
//	}
//}
//
//// Initialize worker pool
//func initWorkerPool() {
//	for i := 0; i < maxWorkers; i++ {
//		go runWorker(i, scanJobQueue)
//	}
//}
//
//// Retry logic for transient errors
//func retry(attempts int, sleep time.Duration, fn func() error) error {
//	for i := 0; i < attempts; i++ {
//		if err := fn(); err != nil {
//			log.Printf("Attempt %d failed: %v", i+1, err)
//			time.Sleep(sleep)
//		} else {
//			return nil
//		}
//	}
//	return fmt.Errorf("after %d attempts, operation failed", attempts)
//}
//
//// StartScan POST /api/v0/scan?config=/path/to/file
////func StartScan(w http.ResponseWriter, r *http.Request) {
////	configFile := r.URL.Query().Get("config")
////	if configFile == "" {
////		http.Error(w, "config file path missing", http.StatusBadRequest)
////		return
////	}
////
////	// Check if the file exists
////	if _, err := os.Stat(configFile); os.IsNotExist(err) {
////		http.Error(w, "config file not found", http.StatusNotFound)
////		log.Printf("Error: Config file not found at path: %s", configFile)
////		return
////	}
////
////	cfg, err := config.LoadAndValidateConfig(configFile)
////	if err != nil {
////		log.Fatalf("Error: %v", err)
////	}
////
////	// Unique scan ID
////	scanID := generateScanID()
////	updateScanStatus(scanID, StatusQueued)
////
////	// Start the scan asynchronously
////	go runScan(scanID, cfg)
////
////	response := ScanResponse{
////		ScanID: scanID,
////		Status: string(StatusQueued),
////	}
////	w.Header().Set("Content-Type", "application/json")
////	json.NewEncoder(w).Encode(response)
////}
//
////	func runScan(scanID string, cfg *config.Config) {
////		updateScanStatus(scanID, StatusInProgress)
////
////		// Simulating long-running process
////		time.Sleep(10 * time.Second)
////
////		// Process scan targets and images
////		for _, target := range cfg.ScanTargets {
////			// Iterate over the images and scan them
////			for _, image := range target.Images {
////				// ToDo: save the report if visibility=private|public and status=completed
////				err := utils.RunImageScan(target, image, cfg)
////				if err != nil {
////					updateScanStatus(scanID, StatusFailed)
////					return
////				}
////			}
////		}
////
////		// On scan completion, update status
////		updateScanStatus(scanID, StatusCompleted)
////
////		// Save the report (based on scan_id and config)
////		saveTemporaryScanResult(scanID, cfg.Website.FilesPath)
////	}
////func generateScanID() string {
////	return fmt.Sprintf("scan-%d", time.Now().UnixNano())
////}
//
////// How would this work?
////func updateScanStatus(scanID string, status ScanStatus) {
////	statusMutex.Lock()
////	defer statusMutex.Unlock()
////	scanStatusStore[scanID] = string(status)
////}
////func saveTemporaryScanResult(scanID, result string) error {
////	// Logic for saving report, e.g., to /reports/scanID_report.json
////	// if visibilit=private and status==completed
////	// By default it is saving the report with imagename
////
////	// Store the result in /tmp or /tmp/gokakashi/apiscan/ temporary directory
////	resultFilePath := fmt.Sprintf("/tmp/%s_result.json", scanID)
////	// make and store the JSON response scanID:<>, status:<completed|failed>,result:<scanner result>
////
////	err := os.WriteFile(resultFilePath, []byte(result), 0644)
////	if err != nil {
////		return fmt.Errorf("failed to save report: %v", err)
////	}
////
////	log.Printf("Temporary scan result saved to %s", resultFilePath)
////	return nil
////
////}
//
//// GetScanStatus returns the status of a scan
//func GetScanStatus(w http.ResponseWriter, r *http.Request) {
//	scanID := r.URL.Path[len("/api/v0/scan/"):]
//
//	status, exists := getScanStatus(scanID)
//	if !exists {
//		http.Error(w, "Scan ID not found", http.StatusNotFound)
//		return
//	}
//
//	response := ScanResponse{
//		ScanID: scanID,
//		Status: status,
//	}
//
//	if status == string(StatusCompleted) {
//		// Optionally, return the scan result (from Trivy output)
//		responseFilePath := fmt.Sprintf("/reports/%s_report.json", scanID)
//		w.Header().Set("Content-Type", "application/json")
//		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_report.json", scanID))
//		http.ServeFile(w, r, responseFilePath)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(response)
//}
//
//func getScanStatus(scanID string) (string, bool) {
//	statusMutex.Lock()
//	defer statusMutex.Unlock()
//	status, exists := scanStatusStore[scanID]
//	return status, exists
//}
//
////type ScanJob struct {
////	ScanID string
////	Image  string
////	//Visibility string
////	Status   ScanStatus
////	Response string
////}
////
////var (
////	scanJobs  = make(map[string]*ScanJob)
////	jobsMutex sync.RWMutex
////)
////
////// POST /api/v0/scan?config=/path/to/file&visibility=<none|public|private>
////func scanHandler(w http.ResponseWriter, r *http.Request) {
////
////	configPath := r.URL.Query().Get("config")
////	if configPath == "" {
////		http.Error(w, "Config file path is missing", http.StatusUnauthorized)
////		log.Println("Error: Config file path is missing")
////		return
////	}
////	// Tod0: to check if file exists
////	// Todo: Validate the config file.
////	// Load and validate config using the reusable function
////	cfg, err := utils.LoadAndValidateConfig(configPath)
////	if err != nil {
////		http.Error(w, fmt.Sprintf("Config error: %v", err), http.StatusInternalServerError)
////		log.Printf("Error: %v", err)
////		return
////	}
////
////	scanID := fmt.Sprintf("scan_%s", time.Now().Format("20060102150405"))
////
////	// Log the scan request whatever image distinguish the log message with that its an API call
////	log.Printf("Scan requested for image")
////
////	job := &ScanJob{
////		ScanID:   scanID,
////		Image:    image,
////		Status:   StatusQueued,
////		Response: "", // whatever was responsed like the output of of after scanning,
////	}
////
////	jobsMutex.Lock()
////	scanJobs[scanID] = job
////	jobsMutex.Unlock()
////
////	// ToDo: asynchronously run the scans
////	go processScan(scanID, target, image, cfg) // Should we send the pointer for config values? Which is best?
////	response := map[string]string{
////		"scan_id": scanID,
////		"status":  string(job.Status),
////		Response:  "", // whatever was responsed like the output of of after scanning,
////	}
////
////	w.Header().Set("Content-Type", "application/json")
////	json.NewEncoder(w).Encode(response)
////	w.Header().Set("Content-Type", "application/json")
////	json.NewEncoder(w).Encode(response)
////
////}
////
////func processScan(scanID, target, image, cfg) {
////	time.Sleep(10 * time.Second)
////	log.Printf("something")
////	// separate process starts that executes RunImageScan(target config.ScanTarget, image config.Image, cfg *config.Config it should take with what gokakshi is running)
////	for _, target := range cfg.ScanTargets {
////		// Iterate over the images and scan them
////		for _, image := range target.Images {
////			utils.RunImageScan(target, image, cfg)
////		}
////	}
////	// Add logic to save the report and handle visibility
////}
//
////// GET /api/v0/scan/{scan_id}/status
////func statusHandler(w http.ResponseWriter, r *http.Request) {
////	// go fetches the scan_id
////	// gets the status of if scan is still in queue, in-progress, if done give the output
////	// response to requests
////	// same validation scan_id exists?
////
////	// Parse scan_id from the URL
////	parts := strings.Split(r.URL.Path, "/")
////	if len(parts) < 5 {
////		http.Error(w, "Invalid scan ID", http.StatusBadRequest)
////		return
////	}
////	scanID := parts[4]
////
////	// Simulate fetching scan status (in a real-world scenario, you'd check a database or state store)
////	status := "completed"   // Example status; you'd retrieve this dynamically
////	visibility := "private" // Example visibility; retrieve dynamically as well
////
////	response := map[string]string{
////		"scan_id":    scanID,
////		"status":     status,
////		"visibility": visibility,
////		"report_url": "/reports/" + scanID,
////	}
////	w.Header().Set("Content-Type", "application/json")
////	json.NewEncoder(w).Encode(response)
////}
////
////// This handler checks the status of the scan based on the scan_id and returns the current status (e.g., queued, in-progress, completed).
////
////func serveReport(w http.ResponseWriter, r *http.Request) {
////	scanID := r.URL.Query().Get("scan_id")
////
////	// Fetch scan metadata (including visibility)
////	metadata, err := getReportMetadata(scanID) // Assume this fetches the metadata
////	if err != nil {
////		http.Error(w, "Report not found", http.StatusNotFound)
////		return
////	}
////
////	// Enforce visibility rules (e.g., private reports only served on the private server)
////	if metadata.Visibility == "private" && r.Host != "localhost:9090" {
////		http.Error(w, "This report is only accessible via the private server", http.StatusForbidden)
////		return
////	}
////
////	// Serve the report if all checks pass
////	serveReportFile(scanID, w) // Assume you have a function to serve the file
////}
