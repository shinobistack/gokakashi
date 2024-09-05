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


.env: accept config file. 
config file can have 
- list images 
- creds
- website
  - port
  - name
  - file_path/

configurations
LTS config file
Latest hasura config file 

goKakashi <config_file> 

