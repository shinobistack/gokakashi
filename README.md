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
  <a href="https://shinobistack.zulipchat.com/#narrow/channel/486791-gokakashi"><img src="https://img.shields.io/badge/zulip-join_chat-brightgreen.svg" /></a>
<a href="https://github.com/shinobistack/gokakashi/releases"><img src="https://img.shields.io/github/downloads/shinobistack/gokakashi/total" /></a>
</p>

&nbsp;

gokakashi is a security platform to help ship secure software. It aims to

- Be centralized: You need one place to understand your security posture.
- Be vendor-agnostic and open(-sourced).
- Help teams adopt industry standards like [SLSA](https://slsa.dev/).
- For any team of any size.

## Features 🎁

### Container Image Scanning

Find, analyze, and remediate vulnerabilities present in your container images.

- Multiple registries support - scan images from various container image registries — all in one place!
- Vulnerability scanner of your choice.
- Custom notifications - Customize notifications to suit your needs, including where to get notified and control over priority, assignment, due dates etc.
- Scheduled and on-demand scans - Automate your scans with in-built cron jobs or trigger them from your CI.

#### Image Registries

| Regisry |                                    Status                                    |
|--------------|:----------------------------------------------------------------------------:|
| Docker Hub |       [✔️ Available](https://github.com/shinobistack/gokakashi/issues/81)        |
| Google Artifact Registry |       [✔️ Available](https://github.com/shinobistack/gokakashi/issues/82)        |
| GitHub Container Registry |     [In progress ⏳](https://github.com/shinobistack/gokakashi/issues/83)     |
| Amazon Elastic Container Registry | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/84) |
| Azure Container Registry | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/85) |

#### Image Scanners

| Scanner |                                    Status                                    |
|---------|:----------------------------------------------------------------------------:|
| Trivy |       [✔️ Available](https://github.com/shinobistack/gokakashi/issues/86)        |
| Snyk  | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/87) |
| Clair | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/88) |

#### Notification Systems

| Scanner |                                    Status                                    |
|---------|:----------------------------------------------------------------------------:|
| Linear  |                                 [✔️ Available]()                                 |
| Jira    | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/87) |
| Slack   | [Open for contribution](https://github.com/shinobistack/gokakashi/issues/88) |


## Install 🛠️

### Server

```sh
docker run -d ghcr.io/shinobistack/gokakashi server 
```

### Agent

```sh
docker run --rm -it ghcr.io/shinobistack/gokakashi agent
```

## Develop 🏗️

Thanks for your interest in contributing to the project.

You will need [docker](https://docs.docker.com/) and [docker-compose](https://docs.docker.com/compose/) for building gokakashi. You can follow the below workflow after having the software.

```sh
git clone git@github.com:shinobistack/gokakashi.git

cd gokakashi

# Make code changes

docker compose -f docker-compose/base.yaml -f docker-compose/dev.yaml up --build
```

## Transparency & Feedback ✨

We’re excited to share gokakashi early with the community to gather feedback and improve quickly.

Whether you're curious, have suggestions, or your team is looking for a fast and efficient way to streamline your security workflows (and get back to enjoying that extra ice cream 🍨 or your favorite anime 📺), we’d love to hear from you.

- Chat with us on [Zulip](https://shinobistack.zulipchat.com/#narrow/channel/486791-gokakashi) 🗯️
- Report bugs and feature requests on [GitHub](https://github.com/shinobistack/gokakashi/issues/new) :octocat:

Let’s build something awesome together!
