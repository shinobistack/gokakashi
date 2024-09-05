package registry

import (
	"fmt"
	"os/exec"

	"github.com/hasura/goKakashi/pkg/config"
)

// DockerHub struct
type DockerHub struct{}

// NewDockerHub creates a new DockerHub instance
func NewDockerHub() *DockerHub {
	return &DockerHub{}
}

// Login authenticates to Docker Hub
func (d *DockerHub) Login(cfg *config.Config) error {
	cmd := exec.Command("docker", "login", "-u", cfg.DockerUsername, "-p", cfg.DockerPassword)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker login failed: %v, %s", err, string(output))
	}
	return nil
}

// PullImage pulls the specified Docker image
func (d *DockerHub) PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker pull failed: %v, %s", err, string(output))
	}
	return nil
}
