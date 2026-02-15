# Demo Mode

Demo mode lets potential customers sign up and immediately see the app populated with a year of realistic pool maintenance data. Demo users are automatically deleted after 24 hours unless an admin converts them to regular users first.

## Enabling Demo Mode

Start the server with the `--demo` flag:

```sh
./poolvibes serve --demo
```

Or set it in your config file:

```yaml
demo: true
```

Or via environment variable:

```sh
export DEMO=true
```

## Limiting Demo Users

By default, a maximum of 50 concurrent demo users are allowed. Configure with `--demo-max-users`:

```sh
./poolvibes serve --demo --demo-max-users 100
```

Set to `0` for unlimited. When the cap is reached, new signups get a "demo slots are full, please try again later" message.

## How It Works

1. **First user (admin)** is never a demo user -- they sign up normally
2. **Subsequent signups** are flagged as demo users with a 24-hour expiry (subject to the demo user cap)
3. Demo users see **seeded data** across all tabs immediately after signup
4. A **background cleanup job** runs every 15 minutes and deletes expired demo users along with all their data

## Seeded Data

When a demo user signs up, the following data is automatically created:

- **Chemistry logs** (~100 entries over 12 months): realistic pH, chlorine, alkalinity, CYA, hardness, and temperature readings with seasonal variation and occasional out-of-range values
- **Tasks** (6 recurring): water testing, skimmer cleaning, filter backwash, pump checks, wall brushing, equipment inspection -- with a mix of pending and overdue statuses
- **Equipment** (5 items): variable speed pump, sand filter, salt chlorinator, robotic cleaner, gas heater -- with realistic manufacturers, models, and warranty dates
- **Service records** (5 entries): filter cleaning, pump seal replacement, heater tune-up -- tied to equipment with realistic costs
- **Chemicals** (5 items): liquid chlorine, pH decreaser, alkalinity increaser, stabilizer, calcium hardness increaser -- with stock levels and alert thresholds

## Admin Management

Admins can manage demo users from the admin panel (`/admin/users`):

- **Demo badge**: Demo users show a "Demo" tag with an expiry countdown
- **Convert to regular user**: Edit the user and uncheck the "Demo" checkbox to convert them to a permanent account (their seeded data is preserved)
- **Manual deletion**: Disable or delete demo users manually if needed

## Cleanup

Expired demo users are cleaned up automatically:

- The cleanup service runs every **15 minutes**
- It finds all demo users whose `demo_expires_at` timestamp has passed
- All associated data is deleted (chemistry logs, tasks, equipment, service records, chemicals, sessions)
- The user account itself is then deleted
