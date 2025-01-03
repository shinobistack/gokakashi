package registry

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type DockerHubIntegration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d *DockerHubIntegration) Authenticate() error {
	if d.Username == "" || d.Password == "" {
		log.Println("Skipping DockerHub login as no credentials are provided")
		return nil
	}
	log.Println("Attempting to log in to DockerHub using --password-stdin...")

	cmd := exec.Command("docker", "login", "-u", d.Username, "--password-stdin")
	cmd.Stdin = strings.NewReader(d.Password)
	output, err := cmd.CombinedOutput()
	log.Printf("Docker login output: %s", string(output))

	if err != nil {
		return fmt.Errorf("docker login failed: %v, %s", err, string(output))
	}

	log.Println("Successfully logged in to DockerHub.")
	return nil
}

func (d *DockerHubIntegration) PullImage(image string) error {
	log.Printf("Pulling Docker image: %s...", image)

	cmd := exec.Command("docker", "pull", image)
	output, err := cmd.CombinedOutput()
	log.Printf("Docker pull output: %s", string(output))

	if err != nil {
		return fmt.Errorf("docker pull failed: %v, %s", err, string(output))
	}

	log.Println("Docker image pulled successfully.")
	return nil
}
