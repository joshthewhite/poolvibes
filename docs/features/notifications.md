# Notifications

PoolVibes can send email and SMS notifications to alert you when maintenance tasks are due.

## How It Works

A background scheduler runs on a configurable interval (default: 1 hour) and checks for pending tasks due today. For each due task, it sends notifications based on user preferences — at most once per task per day per channel (email/SMS).

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

## Duplicate Prevention

A `task_notifications` table tracks every notification sent. Before sending, the scheduler checks whether a notification has already been sent for the same task, channel, and due date. This prevents duplicate notifications even if the scheduler runs multiple times per day.
