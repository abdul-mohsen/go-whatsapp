#!/bin/bash
#
# Meta WhatsApp Business API Registration and Token Setup Script
#
# This script guides you through the process of:
# 1. Creating a Meta Developer Account
# 2. Creating a Meta Business App
# 3. Setting up WhatsApp Business API
# 4. Generating and managing access tokens
# 5. Configuring webhooks
#
# Author: WhatsApp Go Library
# Version: 1.0.0

# ================================
# Colors for pretty output
# ================================
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# ================================
# Helper Functions
# ================================
print_header() {
    echo ""
    echo -e "${BLUE}${BOLD}========================================"
    echo -e "  $1"
    echo -e "========================================${NC}"
    echo ""
}

print_step() {
    echo -e "${GREEN}[$1]${NC} ${BOLD}$2${NC}"
}

print_substep() {
    echo -e "    ${CYAN}→${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠ WARNING:${NC} $1"
}

print_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

wait_for_user() {
    echo ""
    read -p "${1:-Press Enter to continue...}" </dev/tty
}

open_browser() {
    print_info "Opening: $1"
    if command -v xdg-open &> /dev/null; then
        xdg-open "$1" 2>/dev/null
    elif command -v open &> /dev/null; then
        open "$1"
    else
        print_warning "Could not open browser. Please visit: $1"
    fi
}

# ================================
# Main Script
# ================================
clear
print_header "META WHATSAPP BUSINESS API SETUP"

cat << 'EOF'
This script will guide you through setting up the WhatsApp Business API
with Meta. You will need:

  • A Meta (Facebook) account
  • A business to register
  • A phone number for WhatsApp Business (not used with personal WhatsApp)

The process involves creating a Meta Developer account, setting up a
business app, and configuring the WhatsApp Business API.

EOF

wait_for_user "Press Enter to begin the setup process..."

# ================================
# STEP 1: Meta Developer Account
# ================================
print_header "STEP 1: META DEVELOPER ACCOUNT"

print_step 1 "Create or access your Meta Developer account"
echo ""
print_substep "Go to Meta for Developers website"
print_substep "Log in with your Facebook account"
print_substep "Accept the Developer Terms if this is your first time"
echo ""

read -p "Do you already have a Meta Developer account? (y/n): " response </dev/tty
if [ "$response" = "n" ]; then
    open_browser "https://developers.facebook.com/"
    echo ""
    print_info "Create your developer account, then return here."
    wait_for_user "Press Enter when you have created your developer account..."
fi

print_success "Meta Developer account ready!"

# ================================
# STEP 2: Create Business App
# ================================
print_header "STEP 2: CREATE A META BUSINESS APP"

print_step 2 "Create a new Business App in Meta for Developers"
echo ""
print_substep "Click 'My Apps' in the top right"
print_substep "Click 'Create App'"
print_substep "Select 'Business' as the app type"
print_substep "Enter your app name (e.g., 'My WhatsApp Bot')"
print_substep "Select or create a Business Portfolio"
print_substep "Click 'Create App'"
echo ""

open_browser "https://developers.facebook.com/apps/create/"
wait_for_user "Press Enter when you have created your app..."

print_success "Business App created!"

# ================================
# STEP 3: Add WhatsApp Product
# ================================
print_header "STEP 3: ADD WHATSAPP TO YOUR APP"

print_step 3 "Add the WhatsApp product to your app"
echo ""
print_substep "In your app dashboard, find 'Add products to your app'"
print_substep "Find 'WhatsApp' and click 'Set Up'"
print_substep "Select your Meta Business Account (or create one)"
print_substep "Click 'Continue'"
echo ""

print_info "You should now see WhatsApp in your app's left sidebar"
wait_for_user "Press Enter when WhatsApp is added to your app..."

print_success "WhatsApp product added!"

# ================================
# STEP 4: Get API Credentials
# ================================
print_header "STEP 4: COLLECT YOUR API CREDENTIALS"

print_step 4 "Gather your API credentials from the WhatsApp dashboard"
echo ""
print_substep "Go to WhatsApp > API Setup in the left sidebar"
print_substep "You'll see a temporary access token (valid for 24 hours)"
print_substep "Note your Phone Number ID"
print_substep "Note your WhatsApp Business Account ID"
echo ""

print_warning "The temporary token expires in 24 hours!"
print_info "We'll set up a permanent token in the next step."
echo ""

# Collect credentials
echo -e "${BOLD}Please enter your credentials:${NC}"
echo ""

read -p "Enter your Phone Number ID: " phone_number_id </dev/tty
read -p "Enter your WhatsApp Business Account ID: " business_account_id </dev/tty
read -p "Enter your temporary Access Token: " temp_token </dev/tty

# ================================
# STEP 5: Create System User Token
# ================================
print_header "STEP 5: CREATE A PERMANENT ACCESS TOKEN"

print_step 5 "Create a System User for permanent access"
echo ""
print_substep "Go to Meta Business Suite: business.facebook.com"
print_substep "Click Settings (gear icon) > Business Settings"
print_substep "Go to Users > System Users"
print_substep "Click 'Add' to create a new system user"
print_substep "Name: 'WhatsApp Bot' (or any name)"
print_substep "Role: Admin"
print_substep "Click 'Create System User'"
echo ""

open_browser "https://business.facebook.com/settings/system-users"
wait_for_user "Press Enter when you have created the system user..."

echo ""
print_step 6 "Add assets to the system user"
print_substep "Click on your new system user"
print_substep "Click 'Add Assets'"
print_substep "Select 'Apps' tab"
print_substep "Find your WhatsApp app and toggle it on"
print_substep "Enable 'Full Control'"
print_substep "Click 'Save Changes'"
echo ""

wait_for_user "Press Enter when you have assigned the app to the system user..."

echo ""
print_step 7 "Generate permanent access token"
print_substep "Click 'Generate New Token'"
print_substep "Select your WhatsApp app"
print_substep "Select these permissions:"
echo "       - whatsapp_business_management"
echo "       - whatsapp_business_messaging"
print_substep "Click 'Generate Token'"
print_substep "Copy the token (it won't be shown again!)"
echo ""

read -p "Enter your permanent Access Token (or press Enter to use temp token): " permanent_token </dev/tty
if [ -z "$permanent_token" ]; then
    permanent_token="$temp_token"
    print_warning "Using temporary token. Remember to replace it within 24 hours!"
fi

# ================================
# STEP 6: Get App Secret
# ================================
print_header "STEP 6: GET YOUR APP SECRET"

print_step 8 "Get your App Secret for webhook verification"
echo ""
print_substep "Go to your app in Meta for Developers"
print_substep "Click Settings > Basic in the left sidebar"
print_substep "Click 'Show' next to App Secret"
print_substep "Enter your Facebook password if prompted"
echo ""

open_browser "https://developers.facebook.com/apps/"
print_info "Navigate to your app > Settings > Basic"

read -p "Enter your App Secret: " app_secret </dev/tty

# ================================
# STEP 7: Configure Webhook
# ================================
print_header "STEP 7: WEBHOOK CONFIGURATION"

# Generate a random verify token
verify_token=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

print_step 9 "Set up your webhook endpoint"
echo ""
print_substep "Go to WhatsApp > Configuration in the left sidebar"
print_substep "Click 'Edit' next to Webhook"
print_substep "Enter your webhook URL (must be HTTPS)"
echo -e "    ${CYAN}→${NC} Enter this verify token: ${CYAN}$verify_token${NC}"
print_substep "Click 'Verify and Save'"
echo ""

print_info "Your webhook must be publicly accessible with HTTPS."
print_info "You can use ngrok for local development: ngrok http 8080"
echo ""

print_step 10 "Subscribe to webhook fields"
print_substep "After verifying, click 'Manage'"
print_substep "Subscribe to 'messages' field"
print_substep "Optionally subscribe to other fields as needed"
echo ""

wait_for_user "Press Enter when webhook is configured..."

# ================================
# STEP 8: Generate .env File
# ================================
print_header "GENERATING CONFIGURATION FILE"

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_PATH="$PROJECT_ROOT/.env"

cat > "$ENV_PATH" << EOF
# WhatsApp Business API Configuration
# Generated on: $(date '+%Y-%m-%d %H:%M:%S')
# ====================================

# Your WhatsApp Business Account ID
WHATSAPP_BUSINESS_ACCOUNT_ID=$business_account_id

# Phone Number ID (from Meta Business Suite)
WHATSAPP_PHONE_NUMBER_ID=$phone_number_id

# Access Token (System User Token - permanent)
WHATSAPP_ACCESS_TOKEN=$permanent_token

# Webhook Verification Token (generated for you)
WHATSAPP_WEBHOOK_VERIFY_TOKEN=$verify_token

# App Secret (from Meta for Developers - App Settings)
WHATSAPP_APP_SECRET=$app_secret

# API Version (current stable version)
WHATSAPP_API_VERSION=v18.0

# Webhook Port (for local development)
WEBHOOK_PORT=8080
EOF

print_success "Configuration file created: $ENV_PATH"

# ================================
# STEP 9: Test Phone Number
# ================================
print_header "STEP 8: TEST YOUR SETUP"

print_step 11 "Add a test phone number"
echo ""
print_substep "In WhatsApp > API Setup, find 'To' field"
print_substep "Click 'Manage phone number list'"
print_substep "Add your personal phone number for testing"
print_substep "You'll receive a verification code via WhatsApp"
echo ""

wait_for_user "Press Enter when you have added a test number..."

print_step 12 "Send a test message"
print_substep "Use the 'Send Message' section in API Setup"
print_substep "Select your test phone number"
print_substep "Click 'Send Message'"
print_substep "You should receive a test template message"
echo ""

wait_for_user "Press Enter when you have received the test message..."

# ================================
# Summary
# ================================
print_header "SETUP COMPLETE!"

cat << EOF
$(echo -e "${GREEN}✓${NC}") Your WhatsApp Business API is now configured!

$(echo -e "${BOLD}Configuration Summary:${NC}")
  • Phone Number ID:     $phone_number_id
  • Business Account ID: $business_account_id
  • Webhook Token:       $verify_token
  • Config File:         $ENV_PATH

$(echo -e "${BOLD}Next Steps:${NC}")
  1. Start your Go application with the webhook handler
  2. Use ngrok or similar for local development
  3. Send test messages using the API

$(echo -e "${BOLD}Important Notes:${NC}")
  • Keep your App Secret and Access Token secure
  • Use environment variables, never commit secrets
  • The sandbox has a 24-hour message window limit
  • For production, complete Business Verification

$(echo -e "${BOLD}Resources:${NC}")
  • API Documentation: https://developers.facebook.com/docs/whatsapp
  • Error Codes: https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes
  • Pricing: https://developers.facebook.com/docs/whatsapp/pricing

EOF

echo -e "${BOLD}Happy messaging! 🚀${NC}"
echo ""
