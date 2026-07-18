#!/usr/bin/env bash

set -euo pipefail

NOTICOEL_URL="${NOTICOEL_URL:-http://localhost:8080}"
NOTICOEL_TOKEN="${NOTICOEL_TOKEN:-change-me}"

EVENT_FILE="${1:-../events/workflow-success.json}"

if [ ! -f "$EVENT_FILE" ]; then
    echo "Event file not found: $EVENT_FILE"
    echo ""
    echo "Usage:"
    echo "  ./send.sh ../events/workflow-success.json"
    echo "  ./send.sh ../events/workflow-failure.json"
    echo "  ./send.sh ../events/release.json"
    exit 1
fi

echo "Sending event: $EVENT_FILE"

curl \
    --fail \
    --silent \
    --show-error \
    --request POST \
    --url "${NOTICOEL_URL}/api/v1/events" \
    --header "Authorization: Bearer ${NOTICOEL_TOKEN}" \
    --header "Content-Type: application/json" \
    --data "@${EVENT_FILE}"

echo ""
echo "✓ Event sent successfully"
