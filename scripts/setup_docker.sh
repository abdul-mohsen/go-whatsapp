#!/bin/bash
#
# WhatsApp Business API - Automated Docker Setup
# Linux-first, containerized approach
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_ROOT/.env"

echo -e "${BLUE}${BOLD}"
cat << 'EOF'
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║   🐳 WhatsApp Business API - Docker Setup                        ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check Docker
check_docker() {
    echo -e "${CYAN}Checking Docker...${NC}"
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Docker not found! Install Docker first:${NC}"
        echo -e "  ${YELLOW}https://docs.docker.com/engine/install/${NC}"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        echo -e "${RED}Docker daemon not running! Start Docker first.${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Docker is ready${NC}\n"
}

# Check docker-compose
check_compose() {
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        echo -e "${RED}docker-compose not found!${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Using: $COMPOSE_CMD${NC}\n"
}

# Open URL cross-platform
open_url() {
    local url=$1
    echo -e "${CYAN}► Opening: ${url}${NC}"
    
    if command -v xdg-open &> /dev/null; then
        xdg-open "$url" 2>/dev/null &
    elif command -v open &> /dev/null; then
        open "$url" &
    elif command -v wslview &> /dev/null; then
        wslview "$url" &
    else
        echo -e "${YELLOW}  Please open manually: ${url}${NC}"
    fi
}

# Print step
step() {
    echo -e "\n${BLUE}${BOLD}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}${BOLD}  $1${NC}"
    echo -e "${BLUE}${BOLD}═══════════════════════════════════════════════════════════════${NC}\n"
}

# Wait for user
wait_user() {
    echo ""
    read -p "  Press Enter when done..." </dev/tty
}

# Main setup
main() {
    check_docker
    check_compose
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 1: META DEVELOPER SETUP"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "${BOLD}You need to create a Meta Business App with WhatsApp.${NC}\n"
    
    echo -e "  ${CYAN}1.1${NC} Create Developer Account:"
    echo -e "      ${YELLOW}https://developers.facebook.com/${NC}\n"
    
    echo -e "  ${CYAN}1.2${NC} Create Business App:"
    echo -e "      ${YELLOW}https://developers.facebook.com/apps/create/${NC}"
    echo -e "      → Select 'Business' → Add 'WhatsApp' product\n"
    
    open_url "https://developers.facebook.com/apps/create/"
    wait_user
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 2: COLLECT YOUR CREDENTIALS"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "${BOLD}Get these from: Your App → WhatsApp → API Setup${NC}\n"
    open_url "https://developers.facebook.com/apps/"
    
    echo ""
    read -p "  Enter Phone Number ID: " PHONE_NUMBER_ID </dev/tty
    read -p "  Enter Business Account ID: " BUSINESS_ACCOUNT_ID </dev/tty
    read -p "  Enter Temporary Access Token: " TEMP_TOKEN </dev/tty
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 3: GET APP SECRET"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "${BOLD}Go to: Your App → Settings → Basic → Show App Secret${NC}\n"
    open_url "https://developers.facebook.com/apps/"
    
    echo ""
    read -p "  Enter App Secret: " APP_SECRET </dev/tty
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 4: CREATE PERMANENT TOKEN"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "${YELLOW}The temporary token expires in 24 hours!${NC}\n"
    
    echo -e "${BOLD}Create a permanent token:${NC}"
    echo -e "  1. Go to: ${YELLOW}https://business.facebook.com/settings/system-users${NC}"
    echo -e "  2. Create System User (Admin role)"
    echo -e "  3. Add Assets → Select your App → Full Control"
    echo -e "  4. Generate Token with permissions:"
    echo -e "     ${GREEN}✓ whatsapp_business_management${NC}"
    echo -e "     ${GREEN}✓ whatsapp_business_messaging${NC}"
    echo -e "  5. ${RED}COPY THE TOKEN - shown only once!${NC}\n"
    
    open_url "https://business.facebook.com/settings/system-users"
    
    echo ""
    read -p "  Enter Permanent Token (or Enter for temp): " PERMANENT_TOKEN </dev/tty
    
    if [ -z "$PERMANENT_TOKEN" ]; then
        PERMANENT_TOKEN="$TEMP_TOKEN"
        echo -e "  ${YELLOW}⚠ Using temporary token (24hr expiry)${NC}"
    fi
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 5: GENERATE CONFIGURATION"
    # ═══════════════════════════════════════════════════════════════
    
    # Generate random webhook verify token
    VERIFY_TOKEN=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
    
    echo -e "${GREEN}✓ Generated Webhook Verify Token:${NC} ${CYAN}$VERIFY_TOKEN${NC}\n"
    
    # Create .env file
    cat > "$ENV_FILE" << EOF
# WhatsApp Business API Configuration
# Generated: $(date '+%Y-%m-%d %H:%M:%S')
# ═══════════════════════════════════════════════════════════════

# Meta API Credentials
WHATSAPP_BUSINESS_ACCOUNT_ID=$BUSINESS_ACCOUNT_ID
WHATSAPP_PHONE_NUMBER_ID=$PHONE_NUMBER_ID
WHATSAPP_ACCESS_TOKEN=$PERMANENT_TOKEN
WHATSAPP_APP_SECRET=$APP_SECRET

# Webhook Configuration
WHATSAPP_WEBHOOK_VERIFY_TOKEN=$VERIFY_TOKEN
WEBHOOK_PORT=8080

# API Version
WHATSAPP_API_VERSION=v18.0

# Optional: Ngrok auth token for tunneling
# Get from: https://dashboard.ngrok.com/get-started/your-authtoken
NGROK_AUTHTOKEN=
EOF

    echo -e "${GREEN}✓ Created .env file:${NC} $ENV_FILE\n"
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 6: VERIFY API CONNECTION"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "Testing API connection...\n"
    
    RESPONSE=$(curl -s -w "\n%{http_code}" \
        -H "Authorization: Bearer $PERMANENT_TOKEN" \
        "https://graph.facebook.com/v18.0/$PHONE_NUMBER_ID")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓ API Connection Successful!${NC}"
        echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
    else
        echo -e "${YELLOW}⚠ API test returned code: $HTTP_CODE${NC}"
        echo -e "  This might be normal for new setups\n"
    fi
    
    # ═══════════════════════════════════════════════════════════════
    step "STEP 7: BUILD & RUN DOCKER"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "${BOLD}Building Docker image...${NC}\n"
    
    cd "$PROJECT_ROOT"
    
    # Build
    $COMPOSE_CMD build
    
    echo -e "\n${GREEN}✓ Docker image built successfully!${NC}\n"
    
    # Ask to start
    read -p "Start the bot now? (y/n): " START_BOT </dev/tty
    
    if [ "$START_BOT" = "y" ]; then
        echo -e "\n${CYAN}Starting WhatsApp bot...${NC}\n"
        $COMPOSE_CMD up -d
        
        echo -e "\n${GREEN}✓ Bot is running!${NC}"
        echo -e "  Webhook: ${CYAN}http://localhost:8080/webhook${NC}"
        echo -e "  Logs:    ${CYAN}$COMPOSE_CMD logs -f${NC}\n"
        
        # Ask about ngrok
        read -p "Start ngrok tunnel for webhooks? (y/n): " START_NGROK </dev/tty
        
        if [ "$START_NGROK" = "y" ]; then
            echo -e "\n${CYAN}Starting ngrok...${NC}"
            $COMPOSE_CMD --profile dev up -d ngrok
            
            sleep 3
            
            # Get ngrok URL
            NGROK_URL=$(curl -s http://localhost:4040/api/tunnels | python3 -c "import sys,json; print(json.load(sys.stdin)['tunnels'][0]['public_url'])" 2>/dev/null || echo "")
            
            if [ -n "$NGROK_URL" ]; then
                echo -e "\n${GREEN}✓ Ngrok tunnel active!${NC}"
                echo -e "  ${BOLD}Webhook URL:${NC} ${CYAN}${NGROK_URL}/webhook${NC}"
                echo -e "  Dashboard:  ${CYAN}http://localhost:4040${NC}\n"
                echo -e "${YELLOW}Use this URL in Meta webhook configuration!${NC}"
            else
                echo -e "\n${YELLOW}Check ngrok dashboard: http://localhost:4040${NC}"
            fi
        fi
    fi
    
    # ═══════════════════════════════════════════════════════════════
    step "🎉 SETUP COMPLETE!"
    # ═══════════════════════════════════════════════════════════════
    
    echo -e "${BOLD}Quick Commands:${NC}"
    echo -e "  Start:   ${CYAN}$COMPOSE_CMD up -d${NC}"
    echo -e "  Stop:    ${CYAN}$COMPOSE_CMD down${NC}"
    echo -e "  Logs:    ${CYAN}$COMPOSE_CMD logs -f${NC}"
    echo -e "  Rebuild: ${CYAN}$COMPOSE_CMD up -d --build${NC}"
    echo -e "  Ngrok:   ${CYAN}$COMPOSE_CMD --profile dev up -d${NC}\n"
    
    echo -e "${BOLD}Test sending a message:${NC}"
    echo -e "  ${CYAN}./scripts/send_test.sh <phone_number>${NC}\n"
    
    echo -e "${GREEN}${BOLD}Happy messaging! 🚀${NC}\n"
}

main "$@"
