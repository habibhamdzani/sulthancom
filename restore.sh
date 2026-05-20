#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <backup-directory>"
  echo "Example: $0 backup-20250520_113000"
  exit 1
fi

BACKUP_DIR="$1"

echo "=== Restore MongoDB ==="
if [ -d "$BACKUP_DIR/mongodb" ]; then
  docker cp "$BACKUP_DIR/mongodb" sulthancom-mongodb:/tmp/dump
  docker exec sulthancom-mongodb mongorestore \
    --username admin --password admin123 \
    --authenticationDatabase admin \
    --nsInclude="sulthancom_pos.*" \
    --drop \
    /tmp/dump
  docker exec sulthancom-mongodb rm -rf /tmp/dump
  echo "  MongoDB restored from $BACKUP_DIR/mongodb/"
else
  echo "  SKIP: $BACKUP_DIR/mongodb/ not found"
fi

echo ""
echo "=== Restore Redis ==="
if [ -f "$BACKUP_DIR/redis.rdb" ]; then
  docker compose stop redis
  docker cp "$BACKUP_DIR/redis.rdb" sulthancom-redis:/data/dump.rdb
  docker compose start redis
  echo "  Redis restored (service restarted)"
else
  echo "  SKIP: $BACKUP_DIR/redis.rdb not found (not critical, starts fresh)"
fi

echo ""
echo "Restore selesai."
echo ""
echo "Jika ada error di Redis, mulai ulang semua service:"
echo "  docker compose restart"
