# Cloud Run function: Listen to Slack Events

This directory contains a sample Cloud Run Functions (2nd gen) service that receives Slack `app_mention` events and writes them to Cloud Logging.

## Prerequisites

- A GCP project (you have already run `gcloud config set project <PROJECT_ID>`)
- IAM roles such as Cloud Run Admin / Service Account Admin / Secret Manager Admin / Cloud Functions Admin
- Admin permissions for your Slack workspace
- `gcloud` CLI installed

## GCP resources to prepare

1. **Create a Secret Manager secret for the Slack Signing Secret**  
   Store the Signing Secret you get from the Slack App into Secret Manager.

2. **Create a service account for Cloud Run**  
   Create a service account for running the function and allow it to access the secret.

3. **Deploy the Cloud Run function**  
   Use the `deploy.sh` script in the `function/` directory to deploy the function.

   ```sh
   cd ./function
   ./deploy.sh
   ```

   After deployment, copy the Cloud Run URL (for example: `https://xxxxx.run.app`). You will use this in the Slack App settings.

## Slack configuration

1. **Create a Slack App**  
   Go to [Slack API](https://api.slack.com/apps) and create a new App.  
   If you use an App Manifest, you can start from the following JSON:

   ```json
   {
     "display_information": {
       "name": "slackapp-demo"
     },
     "features": {
       "bot_user": {
         "display_name": "slackapp-demo",
         "always_online": false
       }
     },
     "oauth_config": {
       "scopes": {
         "bot": [
           "app_mentions:read"
         ]
       }
     },
     "settings": {
       "event_subscriptions": {
         "request_url": "https://<cloudrun url>",
         "bot_events": [
           "app_mention"
         ]
       },
       "org_deploy_enabled": false,
       "socket_mode_enabled": false,
       "token_rotation_enabled": false
     }
   }
   ```

   - Replace `request_url` with the Cloud Run function URL you copied after deployment.

2. **Get the Signing Secret and store it in Secret Manager**  
   From the Slack App’s **Basic Information** page, copy the Signing Secret and store it in the `slack-signing-secret` secret created earlier.

3. **Enable Event Subscriptions**  
   - In the Slack App settings, turn **Event Subscriptions** ON and set the **Request URL** to the Cloud Run URL.
   - Under **Subscribe to bot events**, add `app_mention`.

4. **Install the Slack App into your workspace**  
   - From **Install App**, install the App into your workspace.
   - Invite the bot user to the target channel.

## How to test

1. In a Slack channel, mention the bot (for example: `@slackapp-demo hello`).
2. In Cloud Console, open **Logging > Logs Explorer** and check the logs of the Cloud Run function.

Example log entry:

```json
{
  "textPayload": "2025/12/27 22:44:34 Verified Event received: {\"type\":\"app_mention\",\"user\":\"UXXXXXXX\",\"ts\":\"1766875472.167609\",\"text\":\"<@UXXXXXXX> hello\",\"channel\":\"CYYYYYYY\" ...}"
}
```

If a log like this appears, the Cloud Run function is correctly receiving events from Slack.
