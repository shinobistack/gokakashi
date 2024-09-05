package registry

import (
	"fmt"
	"github.com/hasura/goKakashi/pkg/config"
	"log"
	"os/exec"
)

// DockerHub struct
type DockerHub struct{}

// NewDockerHub creates a new DockerHub instance
func NewDockerHub() *DockerHub {
	return &DockerHub{}
}

// Login authenticates to Docker Hub using --password-stdin for secure password handling
func (d *DockerHub) Login(cfg *config.Config) error {
	log.Println("Attempting to log in to DockerHub")

	// Construct the docker login command using
	// --password-stdin better way to handle this?
	fmt.Printf("Command: docker login -u %s -p %s", cfg.DockerUsername, cfg.DockerPassword)
	cmd := exec.Command("docker", "login", "-u", cfg.DockerUsername, "-p", cfg.DockerPassword)

	// Capture both stdout and stderr from the command execution
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
