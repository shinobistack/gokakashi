package registry

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/shinobistack/gokakashi/internal/config/v0"
)

type GCR struct{}

func NewGCR() *GCR {
	return &GCR{}
}

// Login authenticates to GCR using service account JSON key
func (g *GCR) Login(target config.ScanTarget) error {
	if target.Auth.Type == "serviceAccount" && target.Auth.JSONKeyPath != "" {
		log.Printf("Setting GOOGLE_APPLICATION_CREDENTIALS to %s", target.Auth.JSONKeyPath)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", target.Auth.JSONKeyPath)

		log.Println("Authenticating with GCR using service account...")

		cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", target.Auth.JSONKeyPath)
		output, err := cmd.CombinedOutput()
		log.Printf("gcloud auth output: %s", string(output))

		if err != nil {
			return fmt.Errorf("service account activation failed: %v, %s", err, string(output))
		}

		log.Println("Configuring Docker to use gcloud credentials for GCR...")
		cmd = exec.Command("gcloud", "auth", "configure-docker", "gcr.io")
		output, err = cmd.CombinedOutput()
		log.Printf("gcloud configure-docker output: %s", string(output))

		if err != nil {
			return fmt.Errorf("Docker configuration failed: %v, %s", err, string(output))
		}

		log.Println("Successfully authenticated with GCR.")
		return nil
	}
	return fmt.Errorf("Invalid or missing authentication method for GCR")
}

// PullImage pulls a Docker image from GCR
func (g *GCR) PullImage(image string) error {
	log.Printf("Pulling GCR image: %s...", image)

	cmd := exec.Command("docker", "pull", image)
	output, err := cmd.CombinedOutput()
	log.Printf("Docker pull output: %s", string(output))

	if err != nil {
		return fmt.Errorf("docker pull failed: %v, %s", err, string(output))
	}

	log.Println("GCR image pulled successfully.")
	return nil
}
