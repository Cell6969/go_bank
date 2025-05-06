#!/bin/sh

set -e


echo "=== FILES IN /app ===" && ls -la /app && echo "=== END FILE LIST ==="
echo "start the app"
exec "$@"
