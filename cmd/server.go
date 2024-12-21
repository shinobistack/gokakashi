package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	config "github.com/shinobistack/gokakashi/internal/config/v0"
	configv1 "github.com/shinobistack/gokakashi/internal/config/v1"
	restapi "github.com/shinobistack/gokakashi/internal/restapi/v0"
	restapiv1 "github.com/shinobistack/gokakashi/internal/restapi/v1"
	"github.com/shinobistack/gokakashi/pkg/utils"
	"github.com/shinobistack/gokakashi/pkg/web"
	"github.com/shinobistack/gokakashi/webapp"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the GoKakshi Server",
	Run:   runServer,
}

var serverConfigFilePath *string

func runServer(cmd *cobra.Command, args []string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	go func() {
		if os.Getenv("WEB_ONLY") != "" || *serverConfigFilePath == "" {
			return
		}

		handleConfigV1()
	}()

	go func() {
		webServerAddr := ":5555" // TODO make this come from a config
		log.Println("Starting webapp server at", webServerAddr)
		webServer, err := webapp.New(webServerAddr)
		if err != nil {
			log.Fatalln("Error creating web app server", err)
		}
		if err := webServer.ListenAndServe(); err != nil {
			log.Fatalln("Error starting web server", err)
		}
	}()

	go func() {
		if os.Getenv("WEB_ONLY") != "" || *serverConfigFilePath != "" {
			return
		}

		// TODO: get rid of the old config at some point.
		handleConfigV0()
	}()

	<-done
	log.Println("Exiting gokakashi. Bye!")
}

func handleConfigV1() {
	log.Println("=== Starting goKakashi Tool ===")

	// Load and validate the configuration file
	cfg, err := configv1.LoadAndValidateConfig(*serverConfigFilePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Starting API server for scan functionality...")
	s := &restapiv1.Server{
		AuthToken: cfg.Site.APIToken,
		Websites:  cfg.Site.Host,
		Port:      cfg.Site.Port,
	}
	go s.Serve()

	log.Println("Shutting down goKakashi gracefully...")
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
