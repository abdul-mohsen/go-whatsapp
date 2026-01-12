#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Meta WhatsApp Business API Registration and Token Setup Script

.DESCRIPTION
    This script guides you through the process of:
    1. Creating a Meta Developer Account
    2. Creating a Meta Business App
    3. Setting up WhatsApp Business API
    4. Generating and managing access tokens
    5. Configuring webhooks

.NOTES
    Author: WhatsApp Go Library
    Version: 1.0.0
    Requirements: PowerShell 7+ and a web browser
#>

# ================================
# ANSI Color Codes for Pretty Output
# ================================
$ESC = [char]27
$Green = "$ESC[32m"
$Yellow = "$ESC[33m"
$Blue = "$ESC[34m"
$Cyan = "$ESC[36m"
$Red = "$ESC[31m"
$Reset = "$ESC[0m"
$Bold = "$ESC[1m"

function Write-Header {
    param([string]$Text)
    Write-Host ""
    Write-Host "$Blue$Bold========================================$Reset"
    Write-Host "$Blue$Bold  $Text$Reset"
    Write-Host "$Blue$Bold========================================$Reset"
    Write-Host ""
}

function Write-Step {
    param([int]$Number, [string]$Text)
    Write-Host "$Green[$Number]$Reset $Bold$Text$Reset"
}

function Write-SubStep {
    param([string]$Text)
    Write-Host "    $Cyan→$Reset $Text"
}

function Write-Warning {
    param([string]$Text)
    Write-Host "$Yellow⚠ WARNING:$Reset $Text"
}

function Write-Info {
    param([string]$Text)
    Write-Host "$Cyan ℹ$Reset $Text"
}

function Write-Success {
    param([string]$Text)
    Write-Host "$Green ✓$Reset $Text"
}

function Wait-ForUser {
    param([string]$Message = "Press Enter to continue...")
    Write-Host ""
    Read-Host $Message
}

function Open-Browser {
    param([string]$Url)
    Write-Info "Opening: $Url"
    Start-Process $Url
}

# ================================
# Main Script
# ================================

Clear-Host
Write-Header "META WHATSAPP BUSINESS API SETUP"

Write-Host @"
This script will guide you through setting up the WhatsApp Business API
with Meta. You will need:

  • A Meta (Facebook) account
  • A business to register
  • A phone number for WhatsApp Business (not used with personal WhatsApp)

The process involves creating a Meta Developer account, setting up a
business app, and configuring the WhatsApp Business API.

"@

Wait-ForUser "Press Enter to begin the setup process..."

# ================================
# STEP 1: Meta Developer Account
# ================================
Write-Header "STEP 1: META DEVELOPER ACCOUNT"

Write-Step 1 "Create or access your Meta Developer account"
Write-Host ""
Write-SubStep "Go to Meta for Developers website"
Write-SubStep "Log in with your Facebook account"
Write-SubStep "Accept the Developer Terms if this is your first time"
Write-Host ""

$response = Read-Host "Do you already have a Meta Developer account? (y/n)"
if ($response -eq "n") {
    Open-Browser "https://developers.facebook.com/"
    Write-Host ""
    Write-Info "Create your developer account, then return here."
    Wait-ForUser "Press Enter when you have created your developer account..."
}

Write-Success "Meta Developer account ready!"

# ================================
# STEP 2: Create Business App
# ================================
Write-Header "STEP 2: CREATE A META BUSINESS APP"

Write-Step 2 "Create a new Business App in Meta for Developers"
Write-Host ""
Write-SubStep "Click 'My Apps' in the top right"
Write-SubStep "Click 'Create App'"
Write-SubStep "Select 'Business' as the app type"
Write-SubStep "Enter your app name (e.g., 'My WhatsApp Bot')"
Write-SubStep "Select or create a Business Portfolio"
Write-SubStep "Click 'Create App'"
Write-Host ""

Open-Browser "https://developers.facebook.com/apps/create/"
Wait-ForUser "Press Enter when you have created your app..."

Write-Success "Business App created!"

# ================================
# STEP 3: Add WhatsApp Product
# ================================
Write-Header "STEP 3: ADD WHATSAPP TO YOUR APP"

Write-Step 3 "Add the WhatsApp product to your app"
Write-Host ""
Write-SubStep "In your app dashboard, find 'Add products to your app'"
Write-SubStep "Find 'WhatsApp' and click 'Set Up'"
Write-SubStep "Select your Meta Business Account (or create one)"
Write-SubStep "Click 'Continue'"
Write-Host ""

Write-Info "You should now see WhatsApp in your app's left sidebar"
Wait-ForUser "Press Enter when WhatsApp is added to your app..."

Write-Success "WhatsApp product added!"

# ================================
# STEP 4: Get API Credentials
# ================================
Write-Header "STEP 4: COLLECT YOUR API CREDENTIALS"

Write-Step 4 "Gather your API credentials from the WhatsApp dashboard"
Write-Host ""
Write-SubStep "Go to WhatsApp > API Setup in the left sidebar"
Write-SubStep "You'll see a temporary access token (valid for 24 hours)"
Write-SubStep "Note your Phone Number ID"
Write-SubStep "Note your WhatsApp Business Account ID"
Write-Host ""

Write-Warning "The temporary token expires in 24 hours!"
Write-Info "We'll set up a permanent token in the next step."
Write-Host ""

# Collect credentials
Write-Host "$Bold Please enter your credentials:$Reset"
Write-Host ""

$phoneNumberId = Read-Host "Enter your Phone Number ID"
$businessAccountId = Read-Host "Enter your WhatsApp Business Account ID"
$tempToken = Read-Host "Enter your temporary Access Token"

# ================================
# STEP 5: Create System User Token
# ================================
Write-Header "STEP 5: CREATE A PERMANENT ACCESS TOKEN"

Write-Step 5 "Create a System User for permanent access"
Write-Host ""
Write-SubStep "Go to Meta Business Suite: business.facebook.com"
Write-SubStep "Click Settings (gear icon) > Business Settings"
Write-SubStep "Go to Users > System Users"
Write-SubStep "Click 'Add' to create a new system user"
Write-SubStep "Name: 'WhatsApp Bot' (or any name)"
Write-SubStep "Role: Admin"
Write-SubStep "Click 'Create System User'"
Write-Host ""

Open-Browser "https://business.facebook.com/settings/system-users"
Wait-ForUser "Press Enter when you have created the system user..."

Write-Host ""
Write-Step 6 "Add assets to the system user"
Write-SubStep "Click on your new system user"
Write-SubStep "Click 'Add Assets'"
Write-SubStep "Select 'Apps' tab"
Write-SubStep "Find your WhatsApp app and toggle it on"
Write-SubStep "Enable 'Full Control'"
Write-SubStep "Click 'Save Changes'"
Write-Host ""

Wait-ForUser "Press Enter when you have assigned the app to the system user..."

Write-Host ""
Write-Step 7 "Generate permanent access token"
Write-SubStep "Click 'Generate New Token'"
Write-SubStep "Select your WhatsApp app"
Write-SubStep "Select these permissions:"
Write-Host "       - whatsapp_business_management"
Write-Host "       - whatsapp_business_messaging"
Write-SubStep "Click 'Generate Token'"
Write-SubStep "Copy the token (it won't be shown again!)"
Write-Host ""

$permanentToken = Read-Host "Enter your permanent Access Token (or press Enter to use temp token)"
if ([string]::IsNullOrWhiteSpace($permanentToken)) {
    $permanentToken = $tempToken
    Write-Warning "Using temporary token. Remember to replace it within 24 hours!"
}

# ================================
# STEP 6: Get App Secret
# ================================
Write-Header "STEP 6: GET YOUR APP SECRET"

Write-Step 8 "Get your App Secret for webhook verification"
Write-Host ""
Write-SubStep "Go to your app in Meta for Developers"
Write-SubStep "Click Settings > Basic in the left sidebar"
Write-SubStep "Click 'Show' next to App Secret"
Write-SubStep "Enter your Facebook password if prompted"
Write-Host ""

Open-Browser "https://developers.facebook.com/apps/"
Write-Info "Navigate to your app > Settings > Basic"

$appSecret = Read-Host "Enter your App Secret"

# ================================
# STEP 7: Configure Webhook
# ================================
Write-Header "STEP 7: WEBHOOK CONFIGURATION"

# Generate a random verify token
$verifyToken = -join ((65..90) + (97..122) + (48..57) | Get-Random -Count 32 | ForEach-Object {[char]$_})

Write-Step 9 "Set up your webhook endpoint"
Write-Host ""
Write-SubStep "Go to WhatsApp > Configuration in the left sidebar"
Write-SubStep "Click 'Edit' next to Webhook"
Write-SubStep "Enter your webhook URL (must be HTTPS)"
Write-SubStep "Enter this verify token: $Cyan$verifyToken$Reset"
Write-SubStep "Click 'Verify and Save'"
Write-Host ""

Write-Info "Your webhook must be publicly accessible with HTTPS."
Write-Info "You can use ngrok for local development: ngrok http 8080"
Write-Host ""

Write-Step 10 "Subscribe to webhook fields"
Write-SubStep "After verifying, click 'Manage'"
Write-SubStep "Subscribe to 'messages' field"
Write-SubStep "Optionally subscribe to other fields as needed"
Write-Host ""

Wait-ForUser "Press Enter when webhook is configured..."

# ================================
# STEP 8: Generate .env File
# ================================
Write-Header "GENERATING CONFIGURATION FILE"

$envContent = @"
# WhatsApp Business API Configuration
# Generated on: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
# ====================================

# Your WhatsApp Business Account ID
WHATSAPP_BUSINESS_ACCOUNT_ID=$businessAccountId

# Phone Number ID (from Meta Business Suite)
WHATSAPP_PHONE_NUMBER_ID=$phoneNumberId

# Access Token (System User Token - permanent)
WHATSAPP_ACCESS_TOKEN=$permanentToken

# Webhook Verification Token (generated for you)
WHATSAPP_WEBHOOK_VERIFY_TOKEN=$verifyToken

# App Secret (from Meta for Developers - App Settings)
WHATSAPP_APP_SECRET=$appSecret

# API Version (current stable version)
WHATSAPP_API_VERSION=v18.0

# Webhook Port (for local development)
WEBHOOK_PORT=8080
"@

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$envPath = Join-Path $projectRoot ".env"

$envContent | Out-File -FilePath $envPath -Encoding UTF8

Write-Success "Configuration file created: $envPath"

# ================================
# STEP 9: Test Phone Number
# ================================
Write-Header "STEP 8: TEST YOUR SETUP"

Write-Step 11 "Add a test phone number"
Write-Host ""
Write-SubStep "In WhatsApp > API Setup, find 'To' field"
Write-SubStep "Click 'Manage phone number list'"
Write-SubStep "Add your personal phone number for testing"
Write-SubStep "You'll receive a verification code via WhatsApp"
Write-Host ""

Wait-ForUser "Press Enter when you have added a test number..."

Write-Step 12 "Send a test message"
Write-SubStep "Use the 'Send Message' section in API Setup"
Write-SubStep "Select your test phone number"
Write-SubStep "Click 'Send Message'"
Write-SubStep "You should receive a test template message"
Write-Host ""

Wait-ForUser "Press Enter when you have received the test message..."

# ================================
# Summary
# ================================
Write-Header "SETUP COMPLETE!"

Write-Host @"
$Green✓$Reset Your WhatsApp Business API is now configured!

$Bold Configuration Summary:$Reset
  • Phone Number ID:     $phoneNumberId
  • Business Account ID: $businessAccountId
  • Webhook Token:       $verifyToken
  • Config File:         $envPath

$Bold Next Steps:$Reset
  1. Start your Go application with the webhook handler
  2. Use ngrok or similar for local development
  3. Send test messages using the API

$Bold Important Notes:$Reset
  • Keep your App Secret and Access Token secure
  • Use environment variables, never commit secrets
  • The sandbox has a 24-hour message window limit
  • For production, complete Business Verification

$Bold Resources:$Reset
  • API Documentation: https://developers.facebook.com/docs/whatsapp
  • Error Codes: https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes
  • Pricing: https://developers.facebook.com/docs/whatsapp/pricing

"@

Write-Host "$Bold Happy messaging! 🚀$Reset"
Write-Host ""
