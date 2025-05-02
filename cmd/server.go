package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shinobistack/gokakashi/internal/db"

	configv1 "github.com/shinobistack/gokakashi/internal/config/v1"
	restapiv1 "github.com/shinobistack/gokakashi/internal/restapi/v1"
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

func initDatabase(cfg *configv1.Config) {
	// Initialize a separate connection for configuration tasks
	configDB := restapiv1.InitDB(cfg.Database)
	defer configDB.Close()
	db.RunMigrations(configDB)
	db.PopulateDatabase(configDB, cfg)
}
