package registry

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ashwiniag/goKakashi/internal/config/v0"
)

type DockerHub struct{}

func NewDockerHub() *DockerHub {
	return &DockerHub{}
}

// Login authenticates to Docker Hub using --password-stdin
func (d *DockerHub) Login(target config.ScanTarget) error {
	if target.Auth.Type == "basicAuth" && target.Auth.Username == "" && target.Auth.Password == "" {
		log.Println("Skipping DockerHub login as no credentials are provided")
		return nil
	}

	log.Println("Attempting to log in to DockerHub using --password-stdin...")

	cmd := exec.Command("docker", "login", "-u", target.Auth.Username, "--password-stdin")
	cmd.Stdin = strings.NewReader(target.Auth.Password)

	output, err := cmd.CombinedOutput()
	log.Printf("Docker login output: %s", string(output))

	if err != nil {
		return fmt.Errorf("docker login failed: %v, %s", err, string(output))
	}

	log.Println("Successfully logged in to DockerHub.")
	return nil
}

// PullImage pulls the specified Docker image from Docker Hub
func (d *DockerHub) PullImage(image string) error {
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
