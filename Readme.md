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

### Brain dump

This readme will be updated to better. Here for brain dumping. 

**Phase I**
yaml look
```
scan_targets:
  - registry: dockerhub
    auth:
      username: email
      password: xxxx
    images:
      - name: hasura/graphql-engine
        tags:
          - v2.36.0
          - v2.36.3
          - v2.11.13
      - name: hasura/graphql-engine # other repository
        tags:
          - v2.36.4
          - v2.11.8
    scanner:
      - tool: Trivy
  - registry: gcr
    auth:
      username: your-gcr-username
      password: your-gcr-password
    images:
      - name: hasura/graphql-engine
        tags:
          - v2.34.5
          - v2.35.3
          - v2.19.0
    scanner:
      - tool: Trivy

website:
  hostname: localhost
  files_path: /absolute_path/goKakashi/test_data  # desired absolute path
  public:
    port: 8080
  private:
    port: 9090

```

1. having config file. Example ./goKakashi LTS.yaml ./goKakashi stable_release.yaml
   Now incorporate yaml changes into the go code

2. File needs to be stored in the local server. Example ask file_path=?
   and store these reports that get generated by scanner. Example Report generated by trivy

3. Webserver. On the UI, it should list images and upon clicking on these images it should show the report.
   on http://localhost:8080/public/report
   <v2.34.5 image name> --> on clicking it will show me the report
   <2.35.3 image name> --> on clicking it will show me the report
4. It understands absolute path and not relative path for file_path. Change it to relative path
5. Maybe we can restructure the image lists. Like Repository --> drop down images of it ---> on clicking display report
   Some future thinking.
6. Having a Scan policy. Example scan_policy.go would does the some job that on the images that have detected
   HIGH|CRITICAL it should create a linear ticket - this ticket should have the details of HIGH|CRITICAL detected
   and tagging the teams.
   scan policy such that like in yaml we could say if this image detects critical then if Notify is set to linear
   then report it on liner for all CRITICAL details and if Team is mentioned then tag that person on the ticket
   or assign it to the person. Maybe like for critical we can have linear and slack, for High just linear or
   medium and low just linear and no tagging the team
lets, maybe in the yaml we might have something
images:
- name: hasura/graphql-engine
tags:
- v2.36.5
- v2.36.3
- v2.11.8
scan_policy:
- Priority: CRITICAL
- Notify: Linear
- Team: @\ashwini
- name: hasurace/graphql-engine
tags:
- v2.36.5
- v2.36.3
- v2.11.8

==> Tag multiple teams
==> send multiple notify - linear, slack, jira etc
==> On UI should we show only the detected under scan_policy.vulnerabilties reports? Or should we give option to report all
7. Limitation error on description Linear API Response: {"errors":[{"message":"Argument Validation Error","path":["issueCreate"],"locations":[{"line":3,"column":4}],"extensions":{"code":"INVALID_INPUT","type":"invalid input","userError":true,"userPresentableMessage":"description must be shorter than or equal to 250000 characters.","meta":{}}}],"data":null}

8. Authentication. Write now we are just taking basic creds. What if we have basic, token etc authentication
9. beautifying the JSON - like readable.

10. Implement ACR|GCR|ECR Support:

11. Hosting

12. Secure the Private Web Server:

13. Improve Report Presentation:
<.....> 

14. Storage of reports:
<...>

15. Support Additional Scanners:

16. Improve Logging 


