---
sidebar_position: 2
---

gokakashi has the following software components:

- __Postgres Database__ : For all the data storage.
- __gokakashi Server__ : Serves the API and dashboard UI.
- __gokakashi CLI__ : Used for automations.

## Via Docker

Following software are pre-requisites in this installation method:

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

Download the docker-compose configuration.

```sh
# Download the docker compose file
curl -o docker-compose.yml https://raw.githubusercontent.com/shinobistack/gokakashi/refs/heads/main/docker/compose/releases/latest.yaml

# or
wget -O docker-compose.yml https://raw.githubusercontent.com/shinobistack/gokakashi/refs/heads/main/docker/compose/releases/latest.yaml 

docker-compose up -d
```

You can replace `latest.yaml` in the URL above with `edge.yaml` if you wish to use the latest `main` branch build where active development happens. We generally recommend using `latest.yaml` as it will give the latest stable version.