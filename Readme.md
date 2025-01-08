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

gokakashi is a security platform to help ship secure software.

## Motivation 🔥

- Be vendor-agnostic and open(-sourced).
- Centralized: You need one place to understand your security posture.
- Help teams adopt industry standards like [SLSA](https://slsa.dev/).
- Educate: Security is not an afterthought.
- Any team, any size.

## Features 🎁

### Container Image Scanning

Find, analyze, and remediate vulnerabilities present in your container images.

- Multiple registries support - scan images from various container image registries — all in one place!
- Vulnerability scanner of your choice.
- Custom notifications - Customize notifications to suit your needs, including where to get notified and control over priority, assignment, due dates etc.
- Scheduled and on-demand scans - Automate your scans with in-built cron jobs or trigger them from your CI.

#### Image Registries

| Regisry | Status |
|--------------|:-----------------:|
| Docker Hub | ✅ [Enhancement in progress ⏳](https://github.com/shinobistack/gokakashi/issues/81) |
| Google Artifact Registry | ✅ [Enhancement in progress ⏳](https://github.com/shinobistack/gokakashi/issues/82) |
| GitHub Container Registry | [In progress ⏳](https://github.com/shinobistack/gokakashi/issues/83) |
| Amazon Elastic Container Registry | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/84)  |
| Azure Container Registry | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/85) |

#### Image Scanners

| Scanner |                                       Status                                       |
|---------|:----------------------------------------------------------------------------------:|
| Trivy | ✅ [Enhancement in progress ⏳](https://github.com/shinobistack/gokakashi/issues/86) |
| Snyk  |    [Open for contribution](https://github.com/shinobistack/gokakashi/issues/87)    |
| Clair |    [Open for contribution](https://github.com/shinobistack/gokakashi/issues/88)    |


#### Notification Systems

| Platform |                                    Status                                     |
|----------|:-----------------------------------------------------------------------------:|
| Linear   |                                  ✅ Complete                                   |
| Jira     | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/105) |
| Slack    | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/106) |

#### Database Integration
Integrated with PostgreSQL using the modern, type-safe [ent ORM](https://entgo.io/).

| Database     |                                    Status                                     |
|--------------|:-----------------------------------------------------------------------------:|
| PostgresSQL  |                                  ✅ Complete                                   |
| Other DBs | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/107) |
Currently integrated with PostgreSQL using the modern, type-safe [ent ORM](https://entgo.io/).
Future support for databases like MySQL, MariaDB, CockroachDB, SQLite, and more is planned.
With ent ORM, goKakashi is designed to be database-agnostic, allowing seamless integration with existing infrastructure, 
regardless of the database technology is used.

## Install 🛠️
Using Docker Compose

Here’s how you can set up gokakashi using Docker Compose for both the server and PostgreSQL database.
Add your configuration file, e.g., [`./config/latest_config.yaml`](config/latest_config.yaml)


```sh
docker-compose up --build -d
./gokakashi agent start --server=http://localhost:8000 --token=letsdoit --workspace=/tmp

```
Or
### DB
```sh 
docker network create gokakashi-network

docker run -d --rm --name postgresdb --network gokakashi-network -p 5432:5432 \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=postgres \
  postgres:latest

```
### Server

```sh
docker run -d --rm --name gokakashi-server --network gokakashi-network -p 8000:8000 \
  -v $(pwd)/lts.yaml:/app/lts.yaml \
  gokakashi server --config=lts.yaml

```

### Agent

```sh
docker run -it --rm --name gokakashi-agent --network gokakashi-network \
  -v /tmp:/tmp \
  gokakashi agent start --server=http://gokakashi-server:8000 --token=letsdoit --workspace=/tmp
  
```

## Transparency & Feedback ✨
We’re excited to share gokakashi early with the community to gather feedback and improve quickly.

Whether you're curious, have suggestions, or if your team is looking for a fast and efficient way to streamline vulnerability scanning (and get back to enjoying that extra ice cream or your favorite anime), we’d love to hear from you. Feel free to open an issue or submit a pull request or request any features that would help on GitHub. Let’s build something awesome together!
