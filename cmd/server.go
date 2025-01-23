package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shinobistack/gokakashi/internal/assigner"
	"github.com/shinobistack/gokakashi/internal/db"
	"github.com/shinobistack/gokakashi/internal/notifier"

	configv1 "github.com/shinobistack/gokakashi/internal/config/v1"
	restapiv1 "github.com/shinobistack/gokakashi/internal/restapi/v1"
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

	var cfg *configv1.Config
	if *serverConfigFilePath != "" {
		customCfg, err := configv1.LoadAndValidateConfig(*serverConfigFilePath)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		cfg = customCfg
	}
	if cfg == nil {
		defaultCfg, err := configv1.DefaultConfig()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		cfg = defaultCfg
	}

	log.Println("==== Starting gokakashi ====")
	go startAPIServer(cfg)
	go startWebServer(cfg)
	logConfig(cfg)

	<-done
	log.Println("====  Exiting gokakashi. Bye! ==== ")
}

func startAPIServer(cfg *configv1.Config) {
	s := &restapiv1.Server{
		AuthToken: cfg.Site.APIToken,
		Websites:  cfg.Site.Host,
		Port:      cfg.Site.Port,
		DBConfig:  cfg.Database,
	}
	go s.Serve()

	dbConfig := cfg.Database

	// Initialize a separate connection for configuration tasks
	configDB := restapiv1.InitDB(dbConfig)
	defer configDB.Close()
	db.RunMigrations(configDB)
	// Populate the database
	db.PopulateDatabase(configDB, cfg)

	// ToDo: To be go routine who independently and routinely checks and assigns scans in agentTasks table
	go assigner.StartAssigner(cfg.Site.Host, cfg.Site.Port, cfg.Site.APIToken, 1*time.Minute)
	// Todo: To introduce API calls for scanNotify and remove client passing
	go notifier.StartScanNotifier(cfg.Site.Host, cfg.Site.Port, cfg.Site.APIToken, 1*time.Minute)
}

func startWebServer(cfg *configv1.Config) {
	webServer, err := webapp.New(fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port))
	if err != nil {
		log.Fatalln("Error creating web server", err)
	}
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatalln("Error starting web server", err)
	}
}

func logConfig(cfg *configv1.Config) {
	log.Println()
	log.Println("- - - - Configuration - - - - -")
	log.Printf("  API Server URL: %s\n", cfg.APIServerURL())
	if cfg.Site.LogAPITokenOnStartup {
		log.Printf("  API Token: %s\n", cfg.Site.APIToken)
	}
	log.Println("")
	log.Printf("  Web Server URL: %s\n", cfg.WebServerURL())
	log.Println("")
	log.Printf("  Database Host: %s\n", cfg.Database.Host)
	log.Printf("  Database Port: %d\n", cfg.Database.Port)
	log.Printf("  Database User: %s\n", cfg.Database.User)
	log.Printf("  Database Name: %s\n", cfg.Database.Name)
	log.Printf("  Database Password: %s\n", strings.Repeat("*", len(cfg.Database.Password)))
	log.Println("- - - - - - - - - - - - - - - -")
	log.Println()
}
