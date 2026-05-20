#!/usr/bin/env bash
set -euo pipefail

BACKUP_DIR="backup-$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "=== Backup MongoDB ==="
docker exec sulthancom-mongodb mongodump \
  --username admin --password admin123 \
  --authenticationDatabase admin \
  --db sulthancom_pos \
  --out /tmp/dump

docker cp sulthancom-mongodb:/tmp/dump ./"$BACKUP_DIR"/mongodb
docker exec sulthancom-mongodb rm -rf /tmp/dump
echo "  -> $BACKUP_DIR/mongodb/"

echo "=== Backup Redis ==="
docker exec sulthancom-redis redis-cli SAVE
docker cp sulthancom-redis:/data/dump.rdb ./"$BACKUP_DIR"/redis.rdb
echo "  -> $BACKUP_DIR/redis.rdb"

echo ""
echo "Backup selesai: $BACKUP_DIR/"
echo "Total: $(du -sh "$BACKUP_DIR" | cut -f1)"
