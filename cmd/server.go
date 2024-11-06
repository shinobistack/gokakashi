package cmd

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	restapi "github.com/ashwiniag/goKakashi/internal/restapi/server"
	"github.com/ashwiniag/goKakashi/pkg/config/v0"
	"github.com/ashwiniag/goKakashi/pkg/utils"
	"github.com/ashwiniag/goKakashi/pkg/web"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the GoKakshi Server",
	Run:   runServer,
}

var serverConfigFilePath *string

func runServer(cmd *cobra.Command, args []string) {
	if *serverConfigFilePath != "" {
		handleConfigV0()
		return
	}
}

func handleConfigV0() {
	log.Println("=== Starting goKakashi Tool ===")

	// Load and validate the configuration file
	cfg, err := config.LoadAndValidateConfig(*serverConfigFilePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Start web servers to serve reports
	log.Println("Starting web servers...")
	// Initialize web servers
	webServer := web.NewWebServer()
	err = webServer.StartWebServers(cfg)
	if err != nil {
		log.Fatalf("Failed to start web servers: %v", err)
	}

	log.Println("Starting API server for scan functionality at port 8000...")
	s := &restapi.Server{
		AuthToken: cfg.APIToken,
		Websites:  cfg.Websites,
		Port:      8000,
	}
	go s.Serve()

	// Initialize cron job for scheduling scans
	cronSchedule := cron.New()
	// Ensure cron is stopped when program exits
	defer cronSchedule.Stop()

	// Register cron jobs for each scan target
	// Process scan targets and images
	for _, target := range cfg.ScanTargets {
		// Iterate over the images and scan them
		for _, image := range target.Images {
			if image.ScanPolicy.CronSchedule != "" {
				schedule := image.ScanPolicy.CronSchedule
				_, err := cronSchedule.AddFunc(schedule, func() {
					start := time.Now()
					log.Printf("Scheduled scan started at %v for %s:%s", start, image.Name, strings.Join(image.Tags, ", "))
					err := utils.RunImageScan(target, image, cfg)
					if err != nil {
						log.Println("Error running image scan", err)
						return
					}
					log.Printf("Scheduled scan completed at %v for %s:%s", time.Now(), image.Name, strings.Join(image.Tags, ", "))
				})
				if err != nil {
					log.Printf("Invalid cron schedule for image %s: %v", image.Name, err)
				} else {
					log.Printf("Scheduled scan for image %s:%s with cron schedule %s", image.Name, strings.Join(image.Tags, ", "), schedule)
				}
			} else {
				log.Printf("No cron schedule for image %s. Running scan immediately.", image.Name)
				err := utils.RunImageScan(target, image, cfg)
				if err != nil {
					log.Println("Error running image scan", err)
					return
				}
			}
		}
	}
	// Start cron scheduler
	log.Println("Starting cron scheduler...")
	cronSchedule.Start()

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down goKakashi gracefully...")
}
