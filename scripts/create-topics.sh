#!/usr/bin/env bash
set -euo pipefail

BROKER="${BROKER:-localhost:9092}"
PARTITIONS="${PARTITIONS:-6}"
REPLICATION_FACTOR="${REPLICATION_FACTOR:-1}"
TOPICS=(events.main events.retry events.dlq)

for topic in "${TOPICS[@]}"; do
  docker exec "$(docker ps --filter name=kafka --format '{{.Names}}' | head -n 1)" \
    kafka-topics --bootstrap-server "$BROKER" \
    --create --if-not-exists \
    --topic "$topic" \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION_FACTOR"
done

echo "Created topics: ${TOPICS[*]}"
