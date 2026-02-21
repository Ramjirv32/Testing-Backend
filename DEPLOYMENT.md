# Deployment Guide for Render

## Environment Variables Setup

When deploying to Render, set these environment variables in your service settings:

### Required Variables

1. **ENV**
   Set to `production`.

2. **PORT**
   Render sets this automatically, but you can set it to `10000` if needed.

3. **FIREBASE_CREDENTIALS**
   - Open your `ticpin-fa6d2-firebase-adminsdk-fbsvc-53d16fed36.json` file.
   - Copy the **entire contents** of the JSON file.
   - Paste it as the value for this environment variable in Render.
   - This allows the backend to initialize Firebase without needing the physical file in the repository.

### Third-Party Service Keys

4. **GROQ_API_KEY**
   Your Groq API key for AI features.

5. **CASHFREE_CLIENT_ID** & **CASHFREE_CLIENT_SECRET**
   Required for PAN verification and payment features.

### Email SMTP Variables (for Booking Confirmations)

6. **PLAY_APP_PASSWORD**, **DINING_APP_PASSWORD**, **EVENTS_APP_PASSWORD**, **ADMIN_APP_PASSWORD**
   Gmail App Passwords for the respective email accounts.

## Deployment Steps on Render

1. **Create Web Service**: Connect your GitHub repository.
2. **Runtime**: Select `Go`.
3. **Build Command**: `go build -o main .`
4. **Start Command**: `./main`
5. **Add Environment Variables**: Add all the variables listed above.

## Handle Service Account JSON

You have two options for the Firebase Service Account:

### Option A: Environment Variable (Recommended)
Add the entire JSON content to a `FIREBASE_CREDENTIALS` environment variable in Render. This is the cleanest way as you don't need to commit the secret file to Git.

### Option B: Secret File
In Render's "Advanced" or "Environment" tab, use "Secret Files" to upload the `ticpin-fa6d2-firebase-adminsdk-fbsvc-53d16fed36.json` file. 

If you upload it to a different path (e.g., `/etc/secrets/firebase.json`), you MUST set the **`FIREBASE_KEY_PATH`** environment variable to that exact path:
- **Key**: `FIREBASE_KEY_PATH`
- **Value**: `/etc/secrets/firebase.json`

## Verify Deployment

After setting environment variables and deploying, check the logs for:

```
✅ Firebase Auth initialized
✅ Firestore initialized
✅ Firebase Storage initialized
🚀 Firebase Services Sync Complete
```

If you see errors, check:
1. All three Firebase variables are set
2. Private key includes BEGIN/END markers
3. Private key uses `\n` (not actual newlines)
4. Client email matches your project
5. Project ID is correct

## Testing

Test the login endpoint:
```bash
curl -X POST https://your-app.onrender.com/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone": "1234567890"}'
```

You should see a 200 response with "OTP sent successfully".
