#!/usr/bin/env bash
set -euo pipefail

echo "=== Sulthancom POS Deployment ==="

command -v docker >/dev/null 2>&1 || { echo "Error: Docker not installed"; exit 1; }
command -v docker compose >/dev/null 2>&1 || { echo "Error: Docker Compose not installed"; exit 1; }

if [ ! -f .env ]; then
  cp .env.production .env
  echo "Created .env — edit JWT_SECRET before deploying!"
fi

mkdir -p data/mongodb data/redis

echo "Building and starting services..."
docker compose up -d --build

echo ""
echo "=== Deployment complete ==="
echo "  Frontend : http://localhost:3000"
echo "  Backend  : http://localhost:8080"
echo ""
echo "First deploy — seed data:"
echo "  curl -X POST http://localhost:8080/api/seed"
echo ""
echo "Backup data (existing server):"
echo "  ./backup.sh"
echo ""
echo "Restore data (new server):"
echo "  ./restore.sh backup-YYYYMMDD_HHMMSS"
echo ""
echo "Logs:  docker compose logs -f"
