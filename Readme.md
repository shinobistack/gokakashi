<p align="center">
   <img src="https://github.com/user-attachments/assets/d5a52847-eeac-4cbc-a047-7991a003a523">
  <br><br>
  <span><b>gokakashi</b></span>
  <br><br>
  <i>The Centralized Security Platform 🔍 🚀</i>
  <br><br>
  <span>🚧 Heavy work in progress 🚧</span>
  <br><br>
  <a href="https://github.com/shinobistack/gokakashi/actions/workflows/build.yml"><image src="https://github.com/shinobistack/gokakashi/actions/workflows/build.yml/badge.svg" /></a>
</p>

&nbsp;

Make vulnerability management effortless with **gokakashi**! 
This tool simplifies the process of pulling, scanning, reporting, and notifying across all your container images. Gone are the days of manually juggling multiple tools and managing disparate processes—**gokakashi** brings everything under one roof.

## Key Features
1. **Multi-Platform Image Aggregation**
   Pull images from Dockerhub, ECR, GCR, ACR or private hosted repositories—all in one place!  
   _Current Support:_ Dockerhub integration. \
   **Continuously developing** to support more platforms.

2. **Comprehensive Image Scanning**
   Use gokakashi’s multi-scanner support to detect vulnerabilities in your images.
   You have the flexibility to scan based on severity levels like CRITICAL or HIGH or both CRITICAL AND HIGH. By default, gokakashi scans all severities.\
   _Current Support:_ Trivy scanner for detailed vulnerability scans.\
   **Continuously developing** to support more scanners.

4. **Scheduled Scans with Cron Jobs**
   Automate your scans with cron jobs. Schedule scans as needed or run them on-demand, eliminating manual work and setting up schedules. 

5. **Custom Notifications & Ticketing**
   Customize notifications to suit your needs, including where to get notified and control over priority, assignment, due dates etc.\
   Automatically create and assign issues based on the severity of detected vulnerabilities. gokakashi ensures that new issues are only created when relevant, helping you avoid unnecessary noise.\
   Meaningful tracking is maintained by creating new issues when key details change, such as Vulnerability (CVE), Severity, Installed Version, or Fixed Version\
   Here's an example of the information you'll receive in a notification:
   ```
       Image: ashwiniag/xxx:v2.36.0
       
       Library: libnghttp2-14
       Vulnerability: CVE-2023-44487
       Severity: HIGH
       Status: fixed
       Installed Version: 1.43.0-1build3
       Fixed Version: 1.43.0-1ubuntu0.1
       Title: HTTP/2: Multiple HTTP/2 enabled web servers are vulnerable to a DDoS attack (Rapid Reset Attack)
       More details: https://avd.aquasec.com/nvd/cve-2023-44487
   ```
   _Current Support:_ Linear.\
   Continuously developing to support more platforms like Jira, slack.
5. **Reporting**
   You can define which severity levels to report on, ensuring you only receive the most relevant information. This makes tracking vulnerabilities streamlined and focused. By default, scans for all severity.\
   Enjoy the flexibility to host reports wherever you need, with full control over access—whether through Cloudflare tunnels, SSO login, or other methods.\
   Seamlessly share reports via hosted endpoints, enabling smooth collaboration and quick discussions with your team or clients.\
   Serves scan reports for both public and private access under a unified path `/reports`. Public and private servers run on different ports and can be accessed as follows:
   ```
      
      - Public reports: `http://localhost:Port/reports`
      - Private reports: `http://localhost:Port/reports`
      
      To view an individual report:
      - `/view?file=<filename>`
   ```

7. **API Integrations**
   Need to scan an image during development? Use our API endpoint to scan and get reports on the fly!\
   _Current Support:_ Under development.

## Why Use gokakashi?
- **Reduce Engineering Overhead:** By centralizing the scanning process, gokakashi removes the need for multiple tools and need for managing and collaborating at multiple places.
- **Streamline Release Management:** Automate the detection, reporting, and discussing resolution of vulnerabilities, reducing last-minute firefights. 
- **Increase Security Proactivity:** Catch vulnerabilities before your customers do and maintain their trust with proactive management.
- **Scalability:** Designed to support long-term solutions for managing large-scale image vulnerability detection, gokakashi streamlines everything into a single, centralized platform.
- **Unified Platform:** One tool to rule them all—be it for vulnerability scanning, reporting, or even access and communicating directly with your team!


## Getting Started
1. **Setup Credentials:** Provide your ECR, GCR, Dockerhub, or self-registry credentials. You have the flexibility on how you would like pass it to gokakashi.
   _Current Support:_ Dockerhub.
2. **Schedule Scans:** Set up a cron job to scan your images periodically.
3. **Choose Notification Integration:** Customize your notifications—integrate with Linear Jira or slack to get vulnerability alerts directly in your workflow.
   _Current Support:_ Linear.
4. **Check Reports:** Access both public and private reports via the endpoints and where to store generated reports, defined by gokakashi. Go crazy and customize how you share them internally or with your clients.

**Configuration Example:**
The gokakashi tool is highly configurable, giving you the flexibility to manage different scanning use cases. 
Below is an example of a typical config file:\
```
scan_targets:
  - registry: dockerhub # <current support dockerhub registry>
    auth:
      username: ${DOCKER_USERNAME}
      password: ${DOCKER_PASSWORD}
    images:
      - name: <registry>
        tags:
          - v2.08.0
          - v2.36.3
        scan_policy:
          severity:
            - CRITICAL
            - HIGH
          notify:
            Linear:
              api_key: ${LINEAR_API_KEY}
              project_id: UUID
              team_id: UUID
              issue_title: "Vulnerability Report"
              issue_priority: 2 # INT
              issue_assignee_id: UUID of Assignee
              issue_state_id: UUID of Backlog, Triage, In Progress, etc.
              issue_due_date: 2024-12-01  # YYYY-MM-DD
      - name: <registry>
        tags:
          - v2.36.4
          - v2.11.8
        scan_policy:
          severity:
            - CRITICAL
          notify:
            Linear:
              api_key: ${LINEAR_API_KEY}
              project_id: UUID
              team_id: UUID
              issue_title: "Vulnerability Report"
              issue_priority: 2 # INT
              issue_assignee_id: UUID of Assignee
              issue_state_id: UUID of Backlog, Triage, In Progress, etc.
              issue_due_date: 2024-12-01  # YYYY-MM-DD
    scanner:
      - tool: Trivy
website:
  hostname: localhost
  files_path: /app/website  # absolute
  public:
    port: 8080
  private:
    port: 9090

```
**Current Support:** Continuously developing.

## Execution
```
#binary
./gokakshi --config=/config/config.yaml
# Or docker run with mount
docker run -it -v /Users/ashwiniag/config:/app/config -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker -p 8080:8080 -p 9090:9090 gokakashi:latest --config=/app/config/config.yaml

``` 

## Roadmap
- Jira Integration
- Slack Notifications
- API Endpoints for CI/CD to scan during development phase
- GCR, ACR, and Self-hosted registry integration
<more to be dumped from notes>

## Current Phase
gokakashi is currently in active development. Right now, we support:

- Dockerhub Integration
- Trivy for Vulnerability Scanning
- Linear Notifications and Issue Management
- Cron Functionality
- avoids Deduplication of Issues

More features are on the way! 🚀 Stay tuned as we continue to build and improve. Your feedback and pain points are highly appreciated! 🌻 

## Transparency & Feedback ✨
We’re excited to share gokakashi early with the community to gather feedback and improve quickly.\
Whether you're curious, have suggestions, or if your team is looking for a fast and efficient way to streamline vulnerability scanning (and get back to enjoying that extra ice cream or your favorite anime), we’d love to hear from you. Feel free to open an issue or submit a pull request or request any features that would help on GitHub. Let’s build something awesome together!

## Reach Out 💭
If you have any questions, ideas, or just want to connect, feel free to reach me on X (formerly Twitter) at [@AshwiniGaddagi](https://x.com/AshwiniGaddagi). I'd love to hear from you!


