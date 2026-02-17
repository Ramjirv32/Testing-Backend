# Deployment Guide for Render

## Environment Variables Setup

When deploying to Render, set these environment variables in your service settings:

### Required Firebase Variables

1. **FIREBASE_PROJECT_ID**
   ```
   ticpin-website
   ```

2. **FIREBASE_CLIENT_EMAIL**
   ```
   firebase-adminsdk-fbsvc@ticpin-website.iam.gserviceaccount.com
   ```

3. **FIREBASE_PRIVATE_KEY**
   
   ⚠️ **IMPORTANT**: Copy the entire private key from your service account JSON file INCLUDING the `-----BEGIN PRIVATE KEY-----` and `-----END PRIVATE KEY-----` markers.
   
   The key should be in this format (with literal `\n` for newlines):
   ```
   -----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDoH4...\n...\n-----END PRIVATE KEY-----\n
   ```

   **How to format it correctly:**
   - Open your `ticpin-website-firebase-adminsdk-fbsvc-79b256dff0.json` file
   - Copy the ENTIRE value of `private_key` field (including quotes)
   - Remove the outer quotes
   - The newlines should be literal `\n` characters (two characters: backslash + n)
   - Paste it into Render as a single line

### Optional Email Variables

4. **DINING_APP_PASSWORD** - Gmail app password for dining@ticpin.in
5. **EVENTS_APP_PASSWORD** - Gmail app password for events@ticpin.in  
6. **PLAY_APP_PASSWORD** - Gmail app password for play@ticpin.in
7. **ADMIN_APP_PASSWORD** - Gmail app password for admin@ticpin.in

### Other Variables

8. **PORT** (optional, Render sets this automatically)
9. **ENV**
   ```
   production
   ```

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
