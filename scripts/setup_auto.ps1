#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Automated WhatsApp Business API Setup - Direct Links & Token Retrieval

.DESCRIPTION
    This script provides direct links and automates credential gathering
    for the WhatsApp Business API setup.
#>

$ErrorActionPreference = "Continue"

# Colors
$ESC = [char]27
$Green = "$ESC[32m"
$Yellow = "$ESC[33m"
$Blue = "$ESC[34m"
$Cyan = "$ESC[36m"
$Red = "$ESC[31m"
$Reset = "$ESC[0m"
$Bold = "$ESC[1m"

function Write-Header { param([string]$Text)
    Write-Host "`n$Blue$Bold═══════════════════════════════════════════════════════════════$Reset"
    Write-Host "$Blue$Bold  $Text$Reset"
    Write-Host "$Blue$Bold═══════════════════════════════════════════════════════════════$Reset`n"
}

function Write-Link { param([string]$Description, [string]$Url)
    Write-Host "  $Cyan►$Reset $Description"
    Write-Host "    $Yellow$Url$Reset`n"
}

function Open-AndWait { param([string]$Url, [string]$Description)
    Write-Host "$Green→ Opening:$Reset $Description"
    Start-Process $Url
    Read-Host "  Press Enter after completing this step"
}

Clear-Host
Write-Header "WHATSAPP BUSINESS API - AUTOMATED SETUP"

Write-Host @"
$Bold This script will help you get all required credentials: $Reset

  1. WHATSAPP_BUSINESS_ACCOUNT_ID  - Your business account ID
  2. WHATSAPP_PHONE_NUMBER_ID      - Your phone number ID  
  3. WHATSAPP_ACCESS_TOKEN         - Permanent access token
  4. WHATSAPP_WEBHOOK_VERIFY_TOKEN - Auto-generated for you
  5. WHATSAPP_APP_SECRET           - From app settings

"@

# ═══════════════════════════════════════════════════════════════
# STEP 1: Create Meta Developer Account & App
# ═══════════════════════════════════════════════════════════════
Write-Header "STEP 1: META DEVELOPER ACCOUNT & APP"

Write-Host "$Bold Required: Create a Meta Business App with WhatsApp $Reset`n"

Write-Link "1.1 - Meta for Developers (Login/Register)" "https://developers.facebook.com/"
Write-Link "1.2 - Create New App (Select 'Business' type)" "https://developers.facebook.com/apps/create/"
Write-Link "1.3 - Meta Business Suite (Create Business if needed)" "https://business.facebook.com/"

Write-Host "$Yellow Steps:$Reset"
Write-Host "  • Log in with Facebook account"
Write-Host "  • Click 'Create App' → Select 'Business' → Name your app"
Write-Host "  • In app dashboard, find 'WhatsApp' → Click 'Set Up'"
Write-Host "  • Connect or create a Meta Business Account`n"

Open-AndWait "https://developers.facebook.com/apps/create/" "Create App page"

# ═══════════════════════════════════════════════════════════════
# STEP 2: Get Phone Number ID & Business Account ID
# ═══════════════════════════════════════════════════════════════
Write-Header "STEP 2: GET YOUR IDS (Phone Number ID & Business Account ID)"

Write-Host "$Bold These IDs are shown in the WhatsApp API Setup page $Reset`n"

Write-Link "WhatsApp API Setup Page" "https://developers.facebook.com/apps/?show_reminder=true"

Write-Host "$Yellow Where to find them:$Reset"
Write-Host "  • Go to your App → WhatsApp → API Setup (left sidebar)"
Write-Host "  • Phone Number ID: Under 'From' phone number dropdown"
Write-Host "  • Business Account ID: In the URL or API calls section"
Write-Host "  • Copy the 'Temporary access token' (we'll make it permanent later)`n"

Write-Host "$Cyan$Bold Look for:$Reset"
Write-Host @"
  ┌─────────────────────────────────────────────────────────────┐
  │  Phone number ID: $Green 123456789012345 $Reset                        │
  │  WhatsApp Business Account ID: $Green 109876543210987 $Reset           │
  │  Temporary access token: $Green EAAxxxxxxx... $Reset                   │
  └─────────────────────────────────────────────────────────────┘
"@

Open-AndWait "https://developers.facebook.com/apps/" "Your Apps (select your app → WhatsApp → API Setup)"

# Collect the IDs
Write-Host "`n$Bold Enter your IDs from the API Setup page:$Reset`n"
$phoneNumberId = Read-Host "  Phone Number ID"
$businessAccountId = Read-Host "  WhatsApp Business Account ID"
$tempToken = Read-Host "  Temporary Access Token"

# ═══════════════════════════════════════════════════════════════
# STEP 3: Get App Secret
# ═══════════════════════════════════════════════════════════════
Write-Header "STEP 3: GET APP SECRET"

Write-Link "App Settings → Basic" "https://developers.facebook.com/apps/"

Write-Host "$Yellow Steps:$Reset"
Write-Host "  • Go to your App → Settings → Basic"
Write-Host "  • Click 'Show' next to 'App Secret'"
Write-Host "  • Enter your Facebook password if prompted`n"

Open-AndWait "https://developers.facebook.com/apps/" "App Settings (Your App → Settings → Basic)"

$appSecret = Read-Host "  Enter App Secret"

# ═══════════════════════════════════════════════════════════════
# STEP 4: Create Permanent Token (System User)
# ═══════════════════════════════════════════════════════════════
Write-Header "STEP 4: CREATE PERMANENT ACCESS TOKEN"

Write-Host "$Yellow$Bold The temporary token expires in 24 hours! Let's create a permanent one.$Reset`n"

Write-Link "4.1 - System Users Page" "https://business.facebook.com/settings/system-users"
Write-Link "4.2 - Business Settings" "https://business.facebook.com/settings/"

Write-Host "$Bold Step-by-step:$Reset"
Write-Host @"

  $Cyan A. Create System User: $Reset
     • Go to: Business Settings → Users → System Users
     • Click 'Add' button
     • Name: "WhatsApp Bot" 
     • Role: Admin
     • Click 'Create System User'

  $Cyan B. Assign Assets: $Reset
     • Click on your new system user
     • Click 'Add Assets' 
     • Go to 'Apps' tab
     • Find your WhatsApp app, toggle it ON
     • Set to 'Full Control'
     • Click 'Save Changes'

  $Cyan C. Generate Token: $Reset
     • Click 'Generate New Token'
     • Select your app
     • Select permissions:
       ✓ whatsapp_business_management
       ✓ whatsapp_business_messaging  
     • Click 'Generate Token'
     • $Red COPY THE TOKEN NOW - It won't be shown again! $Reset

"@

Open-AndWait "https://business.facebook.com/settings/system-users" "System Users page"

$permanentToken = Read-Host "  Enter Permanent Access Token (or press Enter to use temp token)"
if ([string]::IsNullOrWhiteSpace($permanentToken)) {
    $permanentToken = $tempToken
    Write-Host "`n  $Yellow⚠ Using temporary token - expires in 24 hours!$Reset"
}

# ═══════════════════════════════════════════════════════════════
# STEP 5: Generate Webhook Verify Token
# ═══════════════════════════════════════════════════════════════
Write-Header "STEP 5: WEBHOOK CONFIGURATION"

# Auto-generate verify token
$verifyToken = -join ((65..90) + (97..122) + (48..57) | Get-Random -Count 32 | ForEach-Object {[char]$_})

Write-Host "$Green✓ Auto-generated Webhook Verify Token:$Reset"
Write-Host "  $Cyan$verifyToken$Reset`n"

Write-Link "Webhook Configuration" "https://developers.facebook.com/apps/"

Write-Host "$Bold Configure Webhook (optional - needed for receiving messages):$Reset"
Write-Host @"

  • Go to: Your App → WhatsApp → Configuration
  • Click 'Edit' next to Callback URL
  • Callback URL: Your HTTPS endpoint (e.g., https://yourdomain.com/webhook)
  • Verify Token: $Cyan$verifyToken$Reset
  • Click 'Verify and Save'
  • Subscribe to 'messages' webhook field

  $Yellow For local development, use ngrok:$Reset
    ngrok http 8080
    Then use the HTTPS URL from ngrok

"@

Read-Host "  Press Enter to continue"

# ═══════════════════════════════════════════════════════════════
# STEP 6: Add Test Phone Number
# ═══════════════════════════════════════════════════════════════
Write-Header "STEP 6: ADD TEST PHONE NUMBER"

Write-Link "API Setup - Test Numbers" "https://developers.facebook.com/apps/"

Write-Host "$Bold Add your phone number for testing:$Reset"
Write-Host @"

  • Go to: Your App → WhatsApp → API Setup
  • Find 'To' field → Click 'Manage phone number list'
  • Click 'Add phone number'
  • Enter your personal WhatsApp number (with country code)
  • You'll receive a verification code via WhatsApp
  • Enter the code to verify

  $Yellow Note: In sandbox mode, you can only message verified numbers$Reset

"@

Read-Host "  Press Enter when done"

# ═══════════════════════════════════════════════════════════════
# GENERATE .ENV FILE
# ═══════════════════════════════════════════════════════════════
Write-Header "GENERATING .ENV FILE"

$envContent = @"
# WhatsApp Business API Configuration
# Generated: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
# ═══════════════════════════════════════════════════════════════

# Your WhatsApp Business Account ID
# Found in: App Dashboard → WhatsApp → API Setup
WHATSAPP_BUSINESS_ACCOUNT_ID=$businessAccountId

# Phone Number ID (the ID of your WhatsApp Business phone number)
# Found in: App Dashboard → WhatsApp → API Setup → Phone Number dropdown
WHATSAPP_PHONE_NUMBER_ID=$phoneNumberId

# Access Token (Permanent System User Token)
# Created in: Business Settings → System Users → Generate Token
WHATSAPP_ACCESS_TOKEN=$permanentToken

# Webhook Verification Token (auto-generated, use this when setting up webhook)
WHATSAPP_WEBHOOK_VERIFY_TOKEN=$verifyToken

# App Secret (for webhook signature verification)
# Found in: App Dashboard → Settings → Basic → App Secret
WHATSAPP_APP_SECRET=$appSecret

# API Version
WHATSAPP_API_VERSION=v18.0

# Webhook Server Port
WEBHOOK_PORT=8080
"@

$scriptPath = $PSScriptRoot
if (-not $scriptPath) { $scriptPath = Get-Location }
$projectRoot = Split-Path -Parent $scriptPath
if ($projectRoot -eq "") { $projectRoot = Get-Location }
$envPath = Join-Path $projectRoot ".env"

$envContent | Out-File -FilePath $envPath -Encoding UTF8 -Force

Write-Host "$Green✓ Configuration saved to:$Reset $envPath`n"

# ═══════════════════════════════════════════════════════════════
# VERIFY CONFIGURATION
# ═══════════════════════════════════════════════════════════════
Write-Header "VERIFYING YOUR CONFIGURATION"

Write-Host "Testing API connection...`n"

try {
    $headers = @{
        "Authorization" = "Bearer $permanentToken"
    }
    $testUrl = "https://graph.facebook.com/v18.0/$phoneNumberId"
    $response = Invoke-RestMethod -Uri $testUrl -Headers $headers -Method Get -ErrorAction Stop
    
    Write-Host "$Green✓ API Connection Successful!$Reset"
    Write-Host "  Phone: $($response.display_phone_number)"
    Write-Host "  Verified Name: $($response.verified_name)"
    Write-Host "  Quality Rating: $($response.quality_rating)`n"
}
catch {
    Write-Host "$Yellow⚠ Could not verify API connection$Reset"
    Write-Host "  Error: $($_.Exception.Message)"
    Write-Host "  This might be normal if using a temporary token or new setup`n"
}

# ═══════════════════════════════════════════════════════════════
# SUMMARY
# ═══════════════════════════════════════════════════════════════
Write-Header "SETUP COMPLETE! 🎉"

Write-Host @"
$Bold Your .env file is ready at:$Reset $envPath

$Bold Quick Reference Links:$Reset
"@

Write-Link "App Dashboard" "https://developers.facebook.com/apps/"
Write-Link "API Setup & Testing" "https://developers.facebook.com/docs/whatsapp/cloud-api/get-started"
Write-Link "Send Test Message" "https://developers.facebook.com/apps/$appId/whatsapp-business/wa-dev-console/"
Write-Link "Business Settings" "https://business.facebook.com/settings/"
Write-Link "Webhook Docs" "https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks"
Write-Link "API Reference" "https://developers.facebook.com/docs/whatsapp/cloud-api/reference"
Write-Link "Error Codes" "https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes"

Write-Host @"

$Bold Next Steps:$Reset
  1. Run: $Cyan cd $projectRoot && go mod tidy $Reset
  2. Run: $Cyan go run examples/simple_bot/main.go $Reset
  3. For webhooks, use ngrok: $Cyan ngrok http 8080 $Reset

$Bold Test sending a message:$Reset
"@

# Offer to send test message
$sendTest = Read-Host "`nWould you like to send a test message? (y/n)"
if ($sendTest -eq "y") {
    $testPhone = Read-Host "Enter recipient phone number (with country code, no + or spaces, e.g., 14155551234)"
    
    Write-Host "`nSending test message..."
    
    $body = @{
        messaging_product = "whatsapp"
        to = $testPhone
        type = "template"
        template = @{
            name = "hello_world"
            language = @{
                code = "en_US"
            }
        }
    } | ConvertTo-Json -Depth 10
    
    try {
        $response = Invoke-RestMethod `
            -Uri "https://graph.facebook.com/v18.0/$phoneNumberId/messages" `
            -Headers @{ "Authorization" = "Bearer $permanentToken"; "Content-Type" = "application/json" } `
            -Method Post `
            -Body $body
        
        Write-Host "`n$Green✓ Message sent successfully!$Reset"
        Write-Host "  Message ID: $($response.messages[0].id)"
        Write-Host "  Check your WhatsApp for the message!`n"
    }
    catch {
        Write-Host "`n$Red✗ Failed to send message$Reset"
        Write-Host "  Error: $($_.Exception.Message)"
        Write-Host "`n  Make sure the recipient number is in your test numbers list!`n"
    }
}

Write-Host "`n$Green$Bold Happy messaging! 🚀$Reset`n"
