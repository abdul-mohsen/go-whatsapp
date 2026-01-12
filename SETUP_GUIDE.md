# WhatsApp Business API - Quick Setup Guide

## 🔗 Direct Links (Open These in Order)

### Step 1: Create Meta Developer Account & App
| Action | Link |
|--------|------|
| Login/Create Developer Account | https://developers.facebook.com/ |
| Create New App | https://developers.facebook.com/apps/create/ |
| Meta Business Suite | https://business.facebook.com/ |

**Instructions:**
1. Click "Create App" 
2. Select **"Business"** as app type
3. Name your app (e.g., "My WhatsApp Bot")
4. Select or create a Business Portfolio
5. In the app dashboard, find **"WhatsApp"** and click **"Set Up"**

---

### Step 2: Get Phone Number ID & Business Account ID
| Action | Link |
|--------|------|
| Your Apps Dashboard | https://developers.facebook.com/apps/ |
| WhatsApp API Setup | *(Your App → WhatsApp → API Setup)* |

**Where to find:**
- **Phone Number ID**: Shown under the "From" phone number dropdown
- **Business Account ID**: Shown in the API calls section or URL
- **Temporary Token**: Displayed on the page (valid 24 hours)

---

### Step 3: Get App Secret
| Action | Link |
|--------|------|
| App Settings | https://developers.facebook.com/apps/ → Your App → Settings → Basic |

**Instructions:**
1. Go to Settings → Basic
2. Click "Show" next to App Secret
3. Enter Facebook password if prompted
4. Copy the secret

---

### Step 4: Create Permanent Access Token
| Action | Link |
|--------|------|
| System Users | https://business.facebook.com/settings/system-users |
| Business Settings | https://business.facebook.com/settings/ |

**Instructions:**

**A. Create System User:**
1. Go to Business Settings → Users → System Users
2. Click "Add"
3. Name: "WhatsApp Bot", Role: Admin
4. Click "Create System User"

**B. Assign Assets:**
1. Click on the system user
2. Click "Add Assets"
3. Select "Apps" tab
4. Find your app, toggle ON, set "Full Control"
5. Save Changes

**C. Generate Token:**
1. Click "Generate New Token"
2. Select your app
3. Check these permissions:
   - ✅ `whatsapp_business_management`
   - ✅ `whatsapp_business_messaging`
4. Click "Generate Token"
5. **⚠️ COPY THE TOKEN NOW** - it won't be shown again!

---

### Step 5: Add Test Phone Numbers
| Action | Link |
|--------|------|
| API Setup | *(Your App → WhatsApp → API Setup)* |

**Instructions:**
1. Find the "To" field
2. Click "Manage phone number list"
3. Add your personal WhatsApp number (with country code)
4. Enter verification code sent to your WhatsApp

---

### Step 6: Configure Webhook (Optional - for receiving messages)
| Action | Link |
|--------|------|
| WhatsApp Configuration | *(Your App → WhatsApp → Configuration)* |

**Instructions:**
1. Click "Edit" next to Callback URL
2. Enter your HTTPS webhook URL
3. Enter your verify token (generate any random string)
4. Click "Verify and Save"
5. Subscribe to "messages" field

**For local development:**
```bash
# Install ngrok: https://ngrok.com/download
ngrok http 8080
# Use the HTTPS URL as your webhook
```

---

## 📝 Your .env File Template

```env
# Copy this to .env and fill in your values

WHATSAPP_BUSINESS_ACCOUNT_ID=your_id_here
WHATSAPP_PHONE_NUMBER_ID=your_id_here
WHATSAPP_ACCESS_TOKEN=your_permanent_token_here
WHATSAPP_WEBHOOK_VERIFY_TOKEN=any_random_string_you_create
WHATSAPP_APP_SECRET=your_app_secret_here
WHATSAPP_API_VERSION=v18.0
WEBHOOK_PORT=8080
```

---

## 🧪 Test Your Setup

### Using PowerShell:
```powershell
# Test API connection
$token = "YOUR_ACCESS_TOKEN"
$phoneId = "YOUR_PHONE_NUMBER_ID"

Invoke-RestMethod -Uri "https://graph.facebook.com/v18.0/$phoneId" `
    -Headers @{ Authorization = "Bearer $token" }
```

### Send a test message:
```powershell
$body = @{
    messaging_product = "whatsapp"
    to = "RECIPIENT_PHONE"  # e.g., "14155551234"
    type = "template"
    template = @{
        name = "hello_world"
        language = @{ code = "en_US" }
    }
} | ConvertTo-Json

Invoke-RestMethod -Uri "https://graph.facebook.com/v18.0/$phoneId/messages" `
    -Headers @{ Authorization = "Bearer $token"; "Content-Type" = "application/json" } `
    -Method Post -Body $body
```

---

## 📚 Documentation Links

| Resource | Link |
|----------|------|
| Getting Started | https://developers.facebook.com/docs/whatsapp/cloud-api/get-started |
| API Reference | https://developers.facebook.com/docs/whatsapp/cloud-api/reference |
| Webhooks Guide | https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks |
| Error Codes | https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes |
| Message Templates | https://developers.facebook.com/docs/whatsapp/cloud-api/guides/send-message-templates |
| Pricing | https://developers.facebook.com/docs/whatsapp/pricing |

---

## ⚡ Quick Start Commands

```powershell
# Run the automated setup script
.\scripts\setup_auto.ps1

# Initialize Go module
go mod tidy

# Run the example bot
go run examples/simple_bot/main.go
```
