#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="${1:-go-consumer}"
ACTION="${2:-restart}"

case "$ACTION" in
  kill)
    docker compose kill "$SERVICE_NAME"
    ;;
  stop)
    docker compose stop "$SERVICE_NAME"
    ;;
  restart)
    docker compose restart "$SERVICE_NAME"
    ;;
  *)
    echo "Usage: bash scripts/fault-injection.sh [service] [kill|stop|restart]"
    exit 1
    ;;
esac

echo "Injected fault: $ACTION on $SERVICE_NAME"
