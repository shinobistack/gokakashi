package registry

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type GCRIntegration struct {
	Authtype           string `json:"auth_type"`
	ServiceAccountPath string `json:"json_key_path"`
}

func (g *GCRIntegration) Authenticate() error {
	if _, err := os.Stat(g.ServiceAccountPath); os.IsNotExist(err) {
		return fmt.Errorf("service account file does not exist at %s", g.ServiceAccountPath)
	}

	if g.Authtype == "serviceAccount" && g.ServiceAccountPath != "" {
		log.Printf("Setting GOOGLE_APPLICATION_CREDENTIALS to %s", g.ServiceAccountPath)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", g.ServiceAccountPath)

		log.Println("Authenticating with GCR using service account...")
		cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", g.ServiceAccountPath)

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

func (g *GCRIntegration) PullImage(image string) error {
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
