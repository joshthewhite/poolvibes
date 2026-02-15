# Mobile Responsive Design

## Goal

Make PoolVibes mobile-friendly while maintaining the existing desktop layout.

## Current State

- Bulma 1.0.4 via CDN provides responsive utilities
- Viewport meta tag already present
- Most grids missing explicit mobile breakpoints
- No hamburger menu on mobile
- 9-column chemistry table unusable on phones
- Tabs overflow horizontally on small screens

## Design Decisions

### Navigation: Hamburger Dropdown

- Implement Bulma's built-in `navbar-burger` toggle
- On mobile (<1024px): collapse tabs + email/logout into hamburger dropdown
- On desktop: keep current layout (navbar + horizontal tabs below)
- Use Datastar for toggle state (no extra JS)

### Chemistry Table: Expandable Rows

- On mobile (<769px): hide less-important columns (CC, TA, CYA, CH, Temp) via `is-hidden-mobile`
- Show Date, pH, FC, and a chevron button
- Clicking chevron reveals a detail row below with hidden values
- Toggle via Datastar signals
- On desktop: all columns visible, no chevron

### Admin Table

- Horizontal scroll via `table-container` wrapper (fewer columns, simpler fix)

### Grid Layouts

| Component | Desktop | Tablet | Mobile |
|-----------|---------|--------|--------|
| Dashboard cards | `is-3-desktop` | `is-6-tablet` | `is-12-mobile` |
| Dashboard charts | `is-half` | `is-half` | `is-12-mobile` |
| Chemicals grid | `is-one-third` | `is-half-tablet` | `is-12-mobile` |
| Equipment grid | `is-half` | `is-half` | `is-12-mobile` |
| Quick lists | `is-half` | `is-half` | `is-12-mobile` |

### Forms & Modals

- Form columns: add `is-12-mobile` so fields stack on phones
- Modals: ensure adequate padding on small screens

### Filter Bar (Chemistry)

- Stack filters vertically on mobile
- Date inputs and buttons go full-width

### Touch & Spacing

- Ensure buttons meet 44px minimum tap target
- Adequate spacing between interactive elements on mobile
