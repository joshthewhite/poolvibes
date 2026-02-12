# Tasks

Schedule and track recurring pool maintenance tasks.

## Recurrence Options

Tasks can be configured with three recurrence frequencies:

| Frequency | Example |
|-----------|---------|
| **Daily** | Every N days (e.g., check chlorine every 2 days) |
| **Weekly** | Every N weeks (e.g., backwash filter every week) |
| **Monthly** | Every N months (e.g., inspect equipment every 3 months) |

The interval is configurable — "every 2 weeks" or "every 3 days" are both valid.

## Auto-Rescheduling

When you mark a task as completed, PoolVibes automatically creates the next occurrence based on the recurrence pattern. The new task's due date is calculated from the current due date plus the recurrence interval.

For example, completing a weekly task due on Monday will create the next occurrence due the following Monday.

## Status Tracking

Tasks have three statuses:

| Status | Meaning |
|--------|---------|
| **Pending** | Not yet due or currently due |
| **Completed** | Marked as done (triggers auto-rescheduling) |
| **Overdue** | Past the due date and not yet completed |

## Operations

- **Create** — Add a new recurring task with name, description, recurrence, and due date
- **Edit** — Modify a task's details or recurrence pattern
- **Complete** — Mark as done and auto-generate the next occurrence
- **Delete** — Remove a task entirely
- **List** — View all tasks with their status and due dates
