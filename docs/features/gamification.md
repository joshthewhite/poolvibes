# Gamification

PoolVibes includes gamification features on the dashboard to encourage consistent pool maintenance habits.

## Pool Health Score

A 0-100 score computed from four weighted components:

| Component | Weight | What it measures |
|-----------|--------|------------------|
| Testing Consistency | 30% | Tests in the last 14 days vs. expected (4) |
| Water Quality | 30% | Readings in range on most recent test (6 parameters) |
| Task Completion | 25% | Tasks completed on time in the last 30 days |
| Chemical Stock | 15% | Chemicals above their low-stock threshold |

The score is computed on each dashboard load — no historical score data is stored.

## Streaks

Two streak counters track consecutive weeks of good behavior:

- **Testing Streak** — Consecutive weeks with at least one water test logged
- **Task Streak** — Consecutive weeks with zero overdue tasks (capped at 52)

## Milestone Badges

Eight achievement badges that are earned once and persist in the database:

| Badge | Criteria |
|-------|----------|
| First Dip | Log your first water test |
| Balanced | All 6 readings in range on a single test |
| Consistent | 4-week testing streak |
| Devoted | 12-week testing streak |
| On It | Complete 10 tasks on time |
| Stocked Up | All chemicals above threshold at once |
| Clean Record | 4+ week task streak (zero overdue) |
| Pool Pro | Achieve a health score of 90+ |

Unearned badges appear dimmed. Newly earned badges glow briefly on the dashboard.

## Technical Details

- Health score and streaks are pure functions computed from existing data — no additional database tables needed
- Milestones are persisted in the `user_milestones` table and checked on each dashboard load
- New milestones are saved automatically when their criteria are met
- Demo user milestones are cleaned up when demo accounts expire
