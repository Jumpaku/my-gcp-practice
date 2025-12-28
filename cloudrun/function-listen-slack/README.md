

1. Create Slack App
2. Add Signing Secret to Secret Manager 
3. Add and Configure Service Account
4. Deploy Cloud Run Functions
5. Configure Slack App Event Subscriptions
6. Install the Slack App
7. Add the Slack App to a channel and test it out!

Slack App Manifest
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


Cloud Logging entry

```json
{
  "textPayload": "2025/12/27 22:44:34 Verified Event received: {\"type\":\"app_mention\",\"ts\":\"1766875472.167609\",\"text\":\"<@U0A5WTJ6Z6G>\",\"channel\":\"C0A5WTMKWG4\", ...}"
}
```