package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
	go initDatabase(cfg)
	go startAPIServer(cfg)
	go startWebServer(cfg)
	go assigner.Start(cfg.Site.Host, cfg.Site.Port, cfg.Site.APIToken, 1*time.Minute)
	go notifier.Start(cfg.Site.Host, cfg.Site.Port, cfg.Site.APIToken, 1*time.Minute)
	log.Println(cfg)

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
}

func startWebServer(cfg *configv1.Config) {
	webServer, err := webapp.New(cfg.WebServerServingAddress(), cfg.WebServer.APIHostURL)
	if err != nil {
		log.Fatalln("Error creating web server", err)
	}
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatalln("Error starting web server", err)
	}
}

func initDatabase(cfg *configv1.Config) {
	// Initialize a separate connection for configuration tasks
	configDB := restapiv1.InitDB(cfg.Database)
	defer configDB.Close()
	db.RunMigrations(configDB)
	db.PopulateDatabase(configDB, cfg)
}
