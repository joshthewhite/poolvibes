# Collapsible Filter Bar on Mobile

## Problem

The chemistry table filter bar takes up significant vertical space on mobile, pushing the actual data below the fold. Filters are used occasionally, not on every page load.

## Design

Collapse the filter bar on mobile by default. Add a "Filters ▼/▲" toggle header visible only on mobile. Desktop is unchanged.

### Implementation

- Add `_chemFiltersOpen` Datastar signal (default `false`) to `ChemistryList`
- Add a clickable toggle header (`pv-filter-toggle`) inside the filter box, showing "Filters ▼/▲"
- Wrap existing filter content in a `pv-filter-content` div, hidden on mobile via signal
- CSS: toggle visible on mobile only, filter content always visible on desktop

### Files

- `internal/interface/web/templates/chemistry.templ` — signal, toggle header, content wrapper
- `internal/interface/web/templates/layout.templ` — CSS for toggle/content visibility
