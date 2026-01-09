#!/bin/bash
#
# Quick start script - run everything
#

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "No .env file found. Running setup..."
    bash scripts/setup_docker.sh
    exit 0
fi

# Determine compose command
if command -v docker-compose &> /dev/null; then
    COMPOSE="docker-compose"
else
    COMPOSE="docker compose"
fi

case "${1:-up}" in
    up|start)
        echo "🚀 Starting WhatsApp bot..."
        $COMPOSE up -d --build
        echo "✓ Bot running at http://localhost:8080"
        echo "  Logs: $COMPOSE logs -f"
        ;;
    down|stop)
        echo "🛑 Stopping..."
        $COMPOSE down
        ;;
    logs)
        $COMPOSE logs -f
        ;;
    dev)
        echo "🚀 Starting with ngrok tunnel..."
        $COMPOSE --profile dev up -d --build
        sleep 3
        NGROK_URL=$(curl -s http://localhost:4040/api/tunnels | python3 -c "import sys,json; print(json.load(sys.stdin)['tunnels'][0]['public_url'])" 2>/dev/null)
        echo "✓ Bot running!"
        echo "  Local:   http://localhost:8080/webhook"
        echo "  Public:  ${NGROK_URL}/webhook"
        echo "  Ngrok:   http://localhost:4040"
        ;;
    test)
        bash scripts/send_test.sh "$2"
        ;;
    setup)
        bash scripts/setup_docker.sh
        ;;
    *)
        echo "Usage: $0 {up|down|logs|dev|test <phone>|setup}"
        ;;
esac
