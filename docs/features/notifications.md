# Notifications

PoolVibes can send email, SMS, and push notifications to alert you when maintenance tasks are due.

## How It Works

A background scheduler runs on a configurable interval (default: 1 hour) and checks for pending tasks due today. All due tasks for a user are batched into a single notification per channel (email/SMS/push), sent at most once per day. If you have multiple tasks due, you'll receive one message listing all of them.

## Channels

### Email (Resend)

Email notifications are sent via the [Resend](https://resend.com) API. To enable, configure your Resend API key in the config file or environment variables. See [Configuration](../configuration.md) for details.

### SMS (Twilio)

SMS notifications are sent via the [Twilio](https://www.twilio.com) API. To enable, configure your Twilio account SID, auth token, and sender phone number. See [Configuration](../configuration.md) for details.

### Web Push (VAPID)

Push notifications are sent via the [Web Push protocol](https://web.dev/push-notifications-overview/) using VAPID keys. This allows instant browser/device notifications without requiring a native app. To enable, generate a VAPID key pair and configure the public key, private key, and contact email. See [Configuration](../configuration.md) for details.

Push subscriptions are stored per-user and per-device. Expired subscriptions (e.g., when a user revokes browser permission) are automatically cleaned up when the push service returns a 410 response.

## User Settings

Each user can configure their notification preferences from the **Settings** tab:

- **Phone Number** — Required for SMS notifications (include country code, e.g., `+15551234567`)
- **Email Notifications** — Toggle email alerts on/off (enabled by default)
- **SMS Notifications** — Toggle SMS alerts on/off (disabled by default)
- **Push Notifications** — Toggle push alerts on/off. Enabling this will prompt the browser for notification permission.

## Batching & Duplicate Prevention

Notifications are batched so that each user receives at most **one notification per channel per day**. A `task_notifications` table tracks sent batches by user, channel, and date. If the scheduler runs multiple times per day, duplicate notifications are prevented by this uniqueness constraint.
