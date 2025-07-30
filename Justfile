# Justfile for dev and clean docker compose commands

# Run development docker compose
# Usage: just dev
dev:
    docker compose -f docker-compose/base.yaml -f docker-compose/dev.yaml up --build

# Clean up docker compose containers and volumes
# Usage: just clean
clean:
    docker compose -f docker-compose/base.yaml -f docker-compose/dev.yaml down -v
