set dotenv-load := true
set dotenv-filename := "docker-compose/.env"

# Run development docker compose
# Usage: just dev
dev:
    docker compose -f docker-compose/base.yaml -f docker-compose/dev.yaml up --build

# Clean up docker compose containers and volumes
# Usage: just clean
clean:
    docker compose -f docker-compose/base.yaml -f docker-compose/dev.yaml down -v

dev-agent:
    go run main.go agent start --server=http://localhost:5556 --token="${GOKAKASHI_API_TOKEN}" --experiments="v2_agents"