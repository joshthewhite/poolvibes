# Gamification Design: Pool Health Score + Milestones

## Goal

Add subtle, practical gamification to PoolVibes that encourages consistent pool care, educates users on good practices, and makes the app more engaging — without feeling game-y or over-the-top. Single-user only, no social features.

## Architecture

Three new elements on the dashboard, computed from existing data:

1. **Pool Health Score** — a 0-100 metric summarizing overall pool care
2. **Streaks** — consecutive-week counters for testing and task completion
3. **Milestones** — 8 achievement badges earned through good pool care

All scores and streaks are computed on-the-fly from existing data. Milestones require one new database table to persist unlock state.

## Pool Health Score

A single 0-100 number displayed prominently at the top of the dashboard. Calculated from 4 weighted components:

| Component | Weight | Measurement |
|-----------|--------|-------------|
| Testing Consistency | 30% | Tests logged in last 14 days vs. expected (4 tests = 100%) |
| Water Quality | 30% | % of 6 parameters in range on most recent test |
| Task Completion | 25% | % of tasks completed on time in last 30 days |
| Chemical Stock | 15% | % of tracked chemicals above low-stock threshold |

Color indicator: green (80+), yellow (50-79), red (<50).

Label text:
- 80+: "Your pool is in great shape"
- 50-79: "A few things need attention"
- <50: "Your pool needs some love"

Computed on each dashboard load. No new database tables.

## Streaks

Two streak counters displayed as small inline text below the Pool Health Score:

- **Testing Streak** — Consecutive weeks with at least one water test. Resets if a full week passes with no test.
- **Task Streak** — Consecutive weeks with zero overdue tasks. Resets the moment a task goes overdue.

Computed on-the-fly from existing `chemistry_log.tested_at` and `task.completed_at`/`task.due_date` timestamps. Hidden when streak is 0-1 weeks. No new database storage.

## Milestones

8 achievement badges that unlock through good pool care:

| Milestone | Icon | Criteria |
|-----------|------|----------|
| First Dip | `fa-droplet` | Log first water test |
| Balanced | `fa-scale-balanced` | All 6 readings in range on a single test |
| Consistent | `fa-calendar-check` | 4-week testing streak |
| Devoted | `fa-fire` | 12-week testing streak |
| On It | `fa-clipboard-check` | Complete 10 tasks on time |
| Stocked Up | `fa-boxes-stacked` | All chemicals above threshold at once |
| Clean Record | `fa-circle-check` | 30 days with zero overdue tasks |
| Pool Pro | `fa-trophy` | Pool Health Score reaches 90+ |

### Storage

New `user_milestones` table:

| Column | Type |
|--------|------|
| id | UUID PK |
| user_id | UUID FK (users) |
| milestone | TEXT (enum key) |
| earned_at | TIMESTAMP |
| UNIQUE | (user_id, milestone) |

Milestones are checked on each dashboard load. If criteria are met and the milestone isn't already in the table, it's inserted. This means milestones persist even if the underlying data later changes (e.g., deleting a chemistry log doesn't un-earn "First Dip").

### Display

Small rounded pill badges with Font Awesome icon + milestone name. Earned badges are teal, locked badges are gray. Newly earned badges get a brief CSS highlight animation (gentle glow or scale-up).

Icons via Font Awesome CDN (free tier, solid style).

## Dashboard Layout (top to bottom)

1. **Pool Health Score** — large number + color + label, full width
2. **Streaks** — small inline text ("3 week testing streak / 5 week task streak")
3. **Milestones** — row of 8 icon badges, wrapping on mobile
4. **Stat Cards** — existing 4 cards (unchanged)
5. **Charts & Lists** — existing trend charts and quick lists (unchanged)

## Tech Stack

- Font Awesome 6 Free (CDN, solid icons)
- New domain entity: `Milestone`
- New repository: `MilestoneRepository`
- New service methods on `DashboardService` for score/streak computation
- New migration for `user_milestones` table
- New templ components for score, streaks, and milestone badges

## What's NOT included

- No XP, levels, or points system
- No leaderboards or social comparison
- No notifications/toasts for milestone unlocks
- No dedicated achievements page
- No settings to disable gamification
