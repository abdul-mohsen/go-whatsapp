#!/bin/bash
#
# Send a test WhatsApp message
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Load .env
if [ -f "$PROJECT_ROOT/.env" ]; then
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

RECIPIENT="${1:-}"

if [ -z "$RECIPIENT" ]; then
    echo "Usage: $0 <phone_number>"
    echo "  Phone number with country code, no + or spaces"
    echo "  Example: $0 14155551234"
    exit 1
fi

echo "Sending test message to $RECIPIENT..."

RESPONSE=$(curl -s -X POST \
    "https://graph.facebook.com/${WHATSAPP_API_VERSION:-v18.0}/${WHATSAPP_PHONE_NUMBER_ID}/messages" \
    -H "Authorization: Bearer $WHATSAPP_ACCESS_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"messaging_product\": \"whatsapp\",
        \"to\": \"$RECIPIENT\",
        \"type\": \"template\",
        \"template\": {
            \"name\": \"hello_world\",
            \"language\": { \"code\": \"en_US\" }
        }
    }")

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

if echo "$RESPONSE" | grep -q '"messages"'; then
    echo -e "\n✓ Message sent successfully!"
else
    echo -e "\n✗ Failed to send message"
    echo "  Make sure $RECIPIENT is in your test phone numbers!"
fi
