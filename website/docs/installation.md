---
sidebar_position: 2
---

# Installation

gokakashi has the following software components:

- __Postgres Database__ : For all the data storage.
- __gokakashi Server__ : Serves the API and dashboard UI.
- __gokakashi CLI__ : Used for automations.

## Using Docker

The following software are required for this installation method:

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

Download the docker-compose configuration using curl:
```sh
curl -o docker-compose.yml https://raw.githubusercontent.com/shinobistack/gokakashi/refs/heads/main/docker/compose/releases/latest.yaml
```
or using wget:
```sh
wget -O docker-compose.yml https://raw.githubusercontent.com/shinobistack/gokakashi/refs/heads/main/docker/compose/releases/latest.yaml
```

The above instruction will install the latest stable version of gokakashi. You can replace `latest.yaml` in the URL with `edge.yaml` to use the latest `main` branch build, where active development is happening.

Change the `POSTGRES_PASSWORD` and `DB_PASSWORD` environment variable values in the docker-compose.yml file to secure values of your choice.

Bring up the containers using docker-compose.
```sh
docker compose up -d
```