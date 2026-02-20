# Notifications

PoolVibes can send email and SMS notifications to alert you when maintenance tasks are due.

## How It Works

A background scheduler runs on a configurable interval (default: 1 hour) and checks for pending tasks due today. All due tasks for a user are batched into a single notification per channel (email/SMS), sent at most once per day. If you have multiple tasks due, you'll receive one message listing all of them.

## Channels

### Email (Resend)

Email notifications are sent via the [Resend](https://resend.com) API. To enable, configure your Resend API key in the config file or environment variables. See [Configuration](../configuration.md) for details.

### SMS (Twilio)

SMS notifications are sent via the [Twilio](https://www.twilio.com) API. To enable, configure your Twilio account SID, auth token, and sender phone number. See [Configuration](../configuration.md) for details.

## User Settings

Each user can configure their notification preferences from the **Settings** tab:

- **Phone Number** — Required for SMS notifications (include country code, e.g., `+15551234567`)
- **Email Notifications** — Toggle email alerts on/off (enabled by default)
- **SMS Notifications** — Toggle SMS alerts on/off (disabled by default)

## Batching & Duplicate Prevention

Notifications are batched so that each user receives at most **one notification per channel per day**. A `task_notifications` table tracks sent batches by user, channel, and date. If the scheduler runs multiple times per day, duplicate notifications are prevented by this uniqueness constraint.
