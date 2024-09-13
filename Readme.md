**goKakashi** is a Go-based tool designed to:

1. Authenticate and Pull Docker Images from multiple container registries:
- Docker Hub 
- AWS Elastic Container Registry (ECR) (To be implemented)
- Google Container Registry (GCR) (To be implemented)
- Azure Container Registry (ACR) (To be implemented)
2. Scan Docker Images using Trivy for vulnerabilities.

3. Serve Scan Reports via:
- Public Web Server: Accessible without authentication. 
- Private Web Server: Secured access for internal use.

### Project Structure

```
goKakashi/
├── cmd/                     # Main entry points of the project
│   └── kakashi.go           # Main function
├── pkg/                     # Core functionality, reusable packages
│   ├── config/
│   │   └── config.go        # Configuration loader
│   ├── registry/
│   │   ├── acr.go           # Azure Container Registry implementation (To be implemented)
│   │   ├── dockerhub.go     # Docker Hub implementation
│   │   ├── ecr.go           # AWS ECR implementation (To be implemented)
│   │   ├── gcr.go           # Google GCR implementation (To be implemented)
│   │   └── interface.go     # Registry interface
│   ├── scanner/
│   │   ├── interface.go     # # Scanner interface (pluggable)
│   │   └── trivy.go         # Trivy scanner implementation
│   └── web/
│       ├── private.go       # Private web server
│       └── public.go        # Public web server
├── internal/
│   └── ...                  # Future internal packages
├── go.mod                   # Go module file
└── go.sum                   # Go module checksum file
```
### Build and Test
```
go build -o goKakashi ./cmd
./goKakashi --config=lts.yaml

```
YAML Configuration Example
goKakashi uses a YAML configuration file to define scan targets, authentication details, and image scan policies.

Example: config.yaml
```
scan_targets:
  - registry: dockerhub
    auth:
      username: ${DOCKER_USERNAME}
      password: ${DOCKER_PASSWORD}
    images:
      - name: xx/xx
        tags:
          - v2.36.0
          - v2.36.3
        scan_policy:
          vulnerabilities:
            - CRITICAL
            - HIGH
          notify:
            Linear:
              api_key: ${LINEAR_API_KEY}
              project_id: UUID
              team_id: UUID
              issue_title: TEST2 Vulnerability Report
              issue_priority: 2
              issue_assignee_id: UUID of Assignee
              issue_state_id: UUID of Backlog, Triage, In progres etc
              issue_due_date: 2024-12-01 #YYYY-MM-DD
website:
  hostname: localhost
  files_path: /app/website # absolute
  public:
    port: 8080
  private:
    port: 9090

```
