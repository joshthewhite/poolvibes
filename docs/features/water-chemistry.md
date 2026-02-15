# Water Chemistry

Log water test results and track your pool's chemical balance over time.

## Parameters

Each chemistry log records the following parameters:

| Parameter | Unit | Ideal Range |
|-----------|------|-------------|
| pH | — | 7.2 – 7.6 |
| Free Chlorine (FC) | ppm | 1.0 – 3.0 |
| Combined Chlorine (CC) | ppm | 0 – 0.5 |
| Total Alkalinity (TA) | ppm | 80 – 120 |
| Cyanuric Acid (CYA) | ppm | 30 – 50 |
| Calcium Hardness (CH) | ppm | 200 – 400 |
| Temperature | °F | — |

Each log entry also supports an optional **Notes** field for recording observations or context.

## In-Range Highlighting

Values outside their ideal range are highlighted automatically in the chemistry log list. This makes it easy to spot readings that need attention without comparing numbers manually.

## Treatment Plans

After logging a chemistry test, click the **Plan** button on any row to generate a treatment plan. The plan calculates specific chemical dosages to bring out-of-range readings back to ideal levels, scaled to your pool's volume.

Treatment plans use generic chemical names (muriatic acid, cal-hypo, baking soda, etc.) and cover corrections for: high/low pH, low free chlorine, high combined chlorine, high/low total alkalinity, low CYA, and low calcium hardness.

To get accurate dosages, set your pool's gallon size in **Settings** under "Pool Details."

## Pagination, Sorting & Filtering

The chemistry log table uses server-side pagination to handle large numbers of entries efficiently.

- **Pagination** — 25 rows per page with previous/next and numbered page controls
- **Sortable columns** — Click column headers (Date, pH, FC, TA, CYA) to sort ascending or descending
- **Date range filter** — Filter logs to a specific date range using From/To date inputs
- **Out of range filter** — Toggle to show only entries with at least one parameter outside its ideal range
- **Persistent state** — Sorting, filtering, and page position are preserved across create, edit, and delete operations

## Operations

- **Create** — Log a new water test with any combination of parameters
- **Edit** — Update a previously recorded test
- **Delete** — Remove a log entry
- **Plan** — Generate a treatment plan with chemical dosages
- **List** — View paginated chemistry logs with sorting and filtering
