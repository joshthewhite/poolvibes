# Mobile Responsive Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make PoolVibes mobile-friendly while maintaining the existing desktop layout.

**Architecture:** Pure CSS + Bulma responsive utilities + Datastar signals for expandable chemistry rows. No new dependencies. All changes are in `.templ` files (templates + inline CSS).

**Tech Stack:** Bulma 1.0.4 responsive classes, Datastar v1 signals, templ templates

---

### Task 1: Hamburger Navigation on Mobile

**Files:**
- Modify: `internal/interface/web/templates/layout.templ`

**Step 1: Add hamburger toggle and restructure navbar**

Replace the navbar and tab navigation sections (lines 474-523) with a combined hamburger nav. The key changes:
- Add a `navbar-burger` button in the `navbar-brand`
- Move tabs AND email/logout into a `navbar-menu` div that Bulma collapses on mobile
- Use a Datastar signal `_menuOpen` to toggle the `is-active` class
- Keep tabs as a separate `section` on desktop using CSS `is-hidden-touch` / `is-hidden-desktop`

In `layout.templ`, replace lines 473-523 (the data-signals div opening through end of tab section):

```templ
<div data-signals:tab="window._savedTab" data-signals:_loading="false" data-signals:_menuOpen="false" data-effect="localStorage.setItem('poolvibes_tab', $tab)">
    <!-- Navbar -->
    <nav class="navbar pv-navbar" role="navigation" aria-label="main navigation">
        <div class="container">
            <div class="navbar-brand">
                <a class="navbar-item" href="/">
                    <strong class="is-size-4">PoolVibes</strong>
                </a>
                <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false"
                    data-class:is-active="$_menuOpen"
                    data-on:click="$_menuOpen = !$_menuOpen">
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                </a>
            </div>
            <div class="navbar-menu" data-class:is-active="$_menuOpen">
                <!-- Mobile nav links (hidden on desktop) -->
                <div class="navbar-start is-hidden-desktop">
                    <a class="navbar-item" data-on:click="$tab = 'dashboard'; @get('/dashboard'); $_menuOpen = false">Dashboard</a>
                    <a class="navbar-item" data-on:click="$tab = 'chemistry'; @get('/chemistry'); $_menuOpen = false">Water Chemistry</a>
                    <a class="navbar-item" data-on:click="$tab = 'tasks'; @get('/tasks'); $_menuOpen = false">Tasks</a>
                    <a class="navbar-item" data-on:click="$tab = 'equipment'; @get('/equipment'); $_menuOpen = false">Equipment</a>
                    <a class="navbar-item" data-on:click="$tab = 'chemicals'; @get('/chemicals'); $_menuOpen = false">Chemicals</a>
                    <a class="navbar-item" data-on:click="$tab = 'settings'; @get('/settings'); $_menuOpen = false">Settings</a>
                    if isAdmin {
                        <a class="navbar-item" data-on:click="$tab = 'admin'; @get('/admin/users'); $_menuOpen = false">Admin</a>
                    }
                    <hr class="navbar-divider"/>
                </div>
                <div class="navbar-end">
                    <span class="navbar-item pv-email">{ email }</span>
                    <div class="navbar-item">
                        <form method="POST" action="/logout">
                            <button type="submit" class="button is-small pv-logout">Logout</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </nav>
    <!-- Tab Navigation (desktop only) -->
    <section class="section is-hidden-touch" style="padding-bottom: 0;">
        <div class="container">
            <div class="tabs pv-tabs">
                <ul>
                    <!-- same tab <li> elements as before -->
                </ul>
            </div>
        </div>
    </section>
```

**Step 2: Add mobile navbar CSS**

Add to the `<style>` block in layout.templ, after the existing navbar styles (around line 88):

```css
/* Navbar burger color */
.navbar.pv-navbar .navbar-burger {
    color: #f1eff8;
}
/* Mobile navbar menu background */
@media screen and (max-width: 1023px) {
    .navbar.pv-navbar .navbar-menu {
        background: var(--pv-navbar);
    }
    .navbar.pv-navbar .navbar-menu .navbar-item {
        color: #f1eff8;
    }
    .navbar.pv-navbar .navbar-menu .navbar-item:hover {
        background: #2a2640;
        color: #ffffff;
    }
    .navbar.pv-navbar .navbar-divider {
        background-color: #2a2640;
    }
}
```

**Step 3: Build and verify**

Run: `task templ && task build`

**Step 4: Commit**

```bash
git add internal/interface/web/templates/layout.templ
git commit -m "feat: add hamburger navigation for mobile"
```

---

### Task 2: Mobile Section Padding & General Responsive CSS

**Files:**
- Modify: `internal/interface/web/templates/layout.templ`

**Step 1: Add mobile responsive CSS**

Add to the `<style>` block in layout.templ, before the closing `</style>`:

```css
/* Mobile responsive adjustments */
@media screen and (max-width: 768px) {
    .section {
        padding-left: 0.75rem;
        padding-right: 0.75rem;
    }
    /* Ensure buttons have adequate tap targets */
    .button.is-small {
        min-height: 2.25rem;
        padding-left: 0.75rem;
        padding-right: 0.75rem;
    }
}
```

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/layout.templ
git commit -m "feat: add mobile padding and tap target CSS"
```

---

### Task 3: Dashboard Grid Mobile Breakpoints

**Files:**
- Modify: `internal/interface/web/templates/dashboard.templ`

**Step 1: Add `is-12-mobile` to all column divs**

Summary cards (lines 14, 30, 46, 66): change `is-3-desktop is-6-tablet` to `is-3-desktop is-6-tablet is-12-mobile`

Chart columns (lines 89, 97): change `is-half` to `is-half-desktop is-12-mobile`

Quick list columns (lines 116, 129): change `is-half` to `is-half-desktop is-12-mobile`

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/dashboard.templ
git commit -m "feat: add mobile breakpoints to dashboard grid"
```

---

### Task 4: Chemistry Table Expandable Rows

**Files:**
- Modify: `internal/interface/web/templates/chemistry.templ`
- Modify: `internal/interface/web/templates/layout.templ` (CSS)

**Step 1: Add mobile CSS for expandable rows**

In `layout.templ`, add to the `<style>` block:

```css
/* Chemistry table mobile expandable rows */
@media screen and (max-width: 768px) {
    .pv-hidden-mobile {
        display: none !important;
    }
    .pv-expand-btn {
        display: inline-flex !important;
    }
    .pv-detail-row td {
        padding-top: 0;
        border-top: none;
    }
}
@media screen and (min-width: 769px) {
    .pv-expand-btn {
        display: none !important;
    }
    .pv-detail-row {
        display: none !important;
    }
}
```

**Step 2: Update chemistry table header**

In `chemistry.templ`, modify the `<thead>` (lines 28-38) to add `pv-hidden-mobile` class to less-important columns and add an expand column header:

```templ
<thead>
    <tr>
        @sortableHeader("Date", "tested_at", data.SortBy, data.SortDir)
        @sortableHeader("pH", "ph", data.SortBy, data.SortDir)
        @sortableHeader("FC", "free_chlorine", data.SortBy, data.SortDir)
        <th class="pv-hidden-mobile">CC</th>
        @sortableHeaderWithClass("TA", "total_alkalinity", data.SortBy, data.SortDir, "pv-hidden-mobile")
        @sortableHeaderWithClass("CYA", "cya", data.SortBy, data.SortDir, "pv-hidden-mobile")
        <th class="pv-hidden-mobile">CH</th>
        <th class="pv-hidden-mobile">Temp</th>
        <th class="has-text-right">Actions</th>
    </tr>
</thead>
```

Note: We need to create a `sortableHeaderWithClass` variant or just add the class to the existing `sortableHeader`. The simplest approach: add `pv-hidden-mobile` as a wrapping class on the `<th>` elements that should be hidden. Since `sortableHeader` generates the `<th>`, we either:
- Option A: Create a new `sortableHeaderMobile` templ function
- Option B: Manually write the sortable th with the extra class for those 2 columns

Best approach: add an optional class param. Create `sortableHeaderHiddenMobile` that wraps sortableHeader behavior with the extra class:

```templ
templ sortableHeaderHiddenMobile(label, col, currentSortBy, currentSortDir string) {
    <th class="pv-hidden-mobile" style="cursor: pointer; user-select: none;" data-on:click={ sortAction(col, currentSortBy, currentSortDir) }>
        { label }{ sortIndicator(col, currentSortBy, currentSortDir) }
    </th>
}
```

**Step 3: Update chemistryRow to support expand/collapse**

Replace the `chemistryRow` templ (lines 128-146) with:

```templ
templ chemistryRow(l entities.ChemistryLog, idx int) {
    <tr>
        <td>{ l.TestedAt.Format("Jan 2, 3:04 PM") }</td>
        <td><span class={ valueClass(l.PHInRange()) }>{ fmtFloat(l.PH, 1) }</span></td>
        <td><span class={ valueClass(l.FreeChlorineInRange()) }>{ fmtFloat(l.FreeChlorine, 1) }</span></td>
        <td class="pv-hidden-mobile"><span class={ valueClass(l.CombinedChlorineInRange()) }>{ fmtFloat(l.CombinedChlorine, 1) }</span></td>
        <td class="pv-hidden-mobile"><span class={ valueClass(l.TotalAlkalinityInRange()) }>{ fmtFloat(l.TotalAlkalinity, 0) }</span></td>
        <td class="pv-hidden-mobile"><span class={ valueClass(l.CYAInRange()) }>{ fmtFloat(l.CYA, 0) }</span></td>
        <td class="pv-hidden-mobile"><span class={ valueClass(l.CalciumHardnessInRange()) }>{ fmtFloat(l.CalciumHardness, 0) }</span></td>
        <td class="pv-hidden-mobile">{ fmt.Sprintf("%.0f", l.Temperature) }&deg;F</td>
        <td class="has-text-right">
            <div class="buttons is-right are-small" style="flex-wrap: nowrap;">
                <button data-on:click={ fmt.Sprintf("$_chemExpand%d = !$_chemExpand%d", idx, idx) } class="button is-small pv-expand-btn">
                    <span data-show={ fmt.Sprintf("!$_chemExpand%d", idx) }>&darr;</span>
                    <span data-show={ fmt.Sprintf("$_chemExpand%d", idx) }>&uarr;</span>
                </button>
                <button data-on:click={ "@get('/chemistry/" + l.ID.String() + "/plan')" } class="button is-info is-outlined is-small">Plan</button>
                <button data-on:click={ "@get('/chemistry/" + l.ID.String() + "/edit')" } class="button is-primary is-outlined is-small">Edit</button>
                <button data-on:click={ "@delete('/chemistry/" + l.ID.String() + "')" } class="button is-danger is-outlined is-small">Delete</button>
            </div>
        </td>
    </tr>
    <tr class="pv-detail-row" data-show={ fmt.Sprintf("$_chemExpand%d", idx) }>
        <td colspan="4">
            <div class="columns is-mobile is-multiline is-size-7">
                <div class="column is-half">
                    <strong>CC:</strong> <span class={ valueClass(l.CombinedChlorineInRange()) }>{ fmtFloat(l.CombinedChlorine, 1) }</span>
                </div>
                <div class="column is-half">
                    <strong>TA:</strong> <span class={ valueClass(l.TotalAlkalinityInRange()) }>{ fmtFloat(l.TotalAlkalinity, 0) }</span>
                </div>
                <div class="column is-half">
                    <strong>CYA:</strong> <span class={ valueClass(l.CYAInRange()) }>{ fmtFloat(l.CYA, 0) }</span>
                </div>
                <div class="column is-half">
                    <strong>CH:</strong> <span class={ valueClass(l.CalciumHardnessInRange()) }>{ fmtFloat(l.CalciumHardness, 0) }</span>
                </div>
                <div class="column is-half">
                    <strong>Temp:</strong> { fmt.Sprintf("%.0f", l.Temperature) }&deg;F
                </div>
            </div>
        </td>
    </tr>
}
```

**Step 4: Update the row loop to pass index**

In the `ChemistryList` templ, change the loop (line 41) from:
```templ
for _, l := range data.Result.Items {
    @chemistryRow(l)
}
```
to:
```templ
for i, l := range data.Result.Items {
    @chemistryRow(l, i)
}
```

**Step 5: Add expand signal initializations**

In the `ChemistryList` templ, on the outer div (lines 11-19), we need signals for expand state. Since we don't know how many rows there are at compile time, we'll use individual signals initialized per row. Actually, a simpler approach: just use `data-signals__ifmissing` on each row's expand button. Datastar will initialize signals that don't exist. So no changes needed on the outer div — the `data-show` on the detail row will default to hidden (falsy).

Actually, we need to initialize the signals. The simplest approach: add `data-signals:_chemExpand0__ifmissing="false"` etc. per row. But this is dynamic. Better: Datastar signals that aren't initialized default to undefined which is falsy, so `data-show="$_chemExpand0"` will hide the row by default. We just need to make sure the signal exists before toggling. Use `__ifmissing` on the button click:

Change the expand button to:
```templ
<button data-signals:={ fmt.Sprintf("_chemExpand%d__ifmissing", idx) + "=\"false\"" } data-on:click={ fmt.Sprintf("$_chemExpand%d = !$_chemExpand%d", idx, idx) } class="button is-small pv-expand-btn">
```

Actually this is getting complex. Simpler: just initialize all signals as false in the data-show attribute — Datastar will treat undefined signals as falsy, so `data-show="$_chemExpand0"` will correctly hide the row when the signal doesn't exist. The toggle `$_chemExpand0 = !$_chemExpand0` will create the signal on first click. This should work fine.

**Step 6: Update filter bar for mobile**

In `chemistryFilterBar`, the `is-narrow` columns should be full-width on mobile. Change each `column is-narrow` to `column is-narrow-tablet`. Actually Bulma doesn't have `is-narrow-tablet`. Instead, we can just let the `is-multiline` handle it — on mobile the narrow columns will already stack. But we should also ensure the filter inputs are usable. The current layout should be OK since `is-multiline` wraps.

No changes needed for the filter bar — `is-multiline` + `is-variable` already handles mobile stacking.

**Step 7: Build and verify**

Run: `task templ && task build`

**Step 8: Commit**

```bash
git add internal/interface/web/templates/chemistry.templ internal/interface/web/templates/layout.templ
git commit -m "feat: add expandable rows to chemistry table on mobile"
```

---

### Task 5: Chemicals Grid Mobile Breakpoints

**Files:**
- Modify: `internal/interface/web/templates/chemicals.templ`

**Step 1: Add mobile breakpoints to ChemicalCard column**

In `chemicals.templ` line 24, change:
```templ
<div class="column is-one-third">
```
to:
```templ
<div class="column is-one-third-desktop is-half-tablet is-12-mobile">
```

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/chemicals.templ
git commit -m "feat: add mobile breakpoints to chemicals grid"
```

---

### Task 6: Equipment Grid Mobile Breakpoints

**Files:**
- Modify: `internal/interface/web/templates/equipment.templ`

**Step 1: Add mobile breakpoint to EquipmentCard column**

In `equipment.templ` line 24, change:
```templ
<div class="column is-half">
```
to:
```templ
<div class="column is-half-desktop is-12-mobile">
```

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/equipment.templ
git commit -m "feat: add mobile breakpoints to equipment grid"
```

---

### Task 7: Admin Table Scroll Container

**Files:**
- Modify: `internal/interface/web/templates/admin.templ`

**Step 1: Wrap table in scroll container**

In `admin.templ`, wrap the `<table>` (line 12) in a `<div class="table-container">`:

```templ
<div class="table-container">
    <table class="table is-fullwidth is-striped is-hoverable">
        ...
    </table>
</div>
```

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/admin.templ
git commit -m "feat: add scroll container to admin table for mobile"
```

---

### Task 8: PageHeader Mobile Layout

**Files:**
- Modify: `internal/interface/web/templates/shared.templ`

**Step 1: Make PageHeader use `is-mobile` level**

In `shared.templ` line 34, change:
```templ
<div class="level">
```
to:
```templ
<div class="level is-mobile">
```

This ensures the title and button stay side-by-side even on mobile instead of stacking.

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/shared.templ
git commit -m "feat: make page headers mobile-friendly"
```

---

### Task 9: Form Column Mobile Stacking

**Files:**
- Modify: `internal/interface/web/templates/chemistry.templ`
- Modify: `internal/interface/web/templates/tasks.templ`
- Modify: `internal/interface/web/templates/equipment.templ`
- Modify: `internal/interface/web/templates/chemicals.templ`

**Step 1: Chemistry form — add is-12-mobile to is-half columns**

In `ChemistryFormFields` (lines 149-222), each `column is-half` should become `column is-half is-12-mobile`. There are 8 such columns.

**Step 2: Task form — add is-12-mobile to columns**

In `TaskFormFields` (lines 99-131), each bare `column` in the `.columns` div should become `column is-12-mobile` so frequency/interval/due date stack on mobile.

**Step 3: Equipment form — add is-12-mobile to columns**

In `EquipmentFormFields` (lines 129-146 and 153-170), each bare `column` should become `column is-12-mobile`.

In `serviceRecordNewFormContent` (lines 249-265), each bare `column` should become `column is-12-mobile`.

**Step 4: Chemical form — add is-12-mobile to columns**

In `ChemicalFormFields` (lines 97-130), each bare `column` should become `column is-12-mobile`.

**Step 5: Build and verify**

Run: `task templ && task build`

**Step 6: Commit**

```bash
git add internal/interface/web/templates/chemistry.templ internal/interface/web/templates/tasks.templ internal/interface/web/templates/equipment.templ internal/interface/web/templates/chemicals.templ
git commit -m "feat: stack form fields on mobile"
```

---

### Task 10: Settings Page Responsive Width

**Files:**
- Modify: `internal/interface/web/templates/settings.templ`

**Step 1: Remove fixed max-width from settings boxes**

In `settings.templ`, change the two box divs with `style="max-width: 500px;"` (lines 20 and 30) to use a Bulma column approach instead:

Replace inline `style="max-width: 500px;"` with `style="max-width: 500px; width: 100%;"` — this keeps the max-width on desktop but allows full width on mobile since the container is already responsive.

Actually the simpler fix: just keep the max-width but ensure it's responsive. `max-width: 500px` with no explicit `width` means the box will be at most 500px but shrink on smaller screens. The box is already inside a `.container` which is responsive. So this should actually be fine already.

No changes needed — `max-width: 500px` already works responsively (the box shrinks on screens < 500px).

**Step 1 (revised): Skip this task — settings are already responsive.**

---

### Task 11: Treatment Plan Modal Mobile Layout

**Files:**
- Modify: `internal/interface/web/templates/chemistry.templ`

**Step 1: Add mobile breakpoints to treatment plan columns**

In `treatmentPlanContent` (lines 310-323), the `.columns.is-multiline` contains `is-half` and `is-one-quarter` columns. Add mobile breakpoints:

Change `column is-half` to `column is-half is-12-mobile`
Change `column is-one-quarter` to `column is-one-quarter is-half-mobile`

**Step 2: Build and verify**

Run: `task templ && task build`

**Step 3: Commit**

```bash
git add internal/interface/web/templates/chemistry.templ
git commit -m "feat: add mobile breakpoints to treatment plan modal"
```

---

### Task 12: Visual Verification

**Step 1: Start dev server**

Run: `task dev`

**Step 2: Test in browser at mobile viewport**

Open the app and resize browser to ~375px width (iPhone SE). Verify:
- Hamburger menu appears and works
- Desktop tabs are hidden on mobile
- Dashboard cards stack to single column
- Chemistry table hides extra columns, chevron reveals details
- Chemical cards are full-width
- Equipment cards are full-width
- Admin table scrolls horizontally
- All forms stack fields on mobile
- Modals are usable
- Buttons have adequate tap targets

**Step 3: Test at tablet viewport (~768px)**

Verify:
- Tabs may or may not show (Bulma touch breakpoint is <1024px)
- Dashboard cards are 2-column
- Chemical cards are 2-column
- Everything looks balanced

**Step 4: Test at desktop viewport (~1280px)**

Verify nothing changed from the current desktop layout:
- Horizontal tabs visible
- No hamburger visible
- All table columns visible, no chevrons
- Grids use full column widths

**Step 5: Commit any fixes**

```bash
git add -A && git commit -m "fix: mobile layout adjustments from visual testing"
```
