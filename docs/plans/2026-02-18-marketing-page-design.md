# Marketing Landing Page Design

**Date:** 2026-02-18
**Status:** Approved

## Overview

Add a marketing landing page as the first screen unauthenticated visitors see. The page showcases PoolVibes features with bold, playful copy and drives visitors to sign up.

## Audience

Consumer pool owners. Focus on simplicity, peace of mind, and saving money — not technical details.

## Tone

Bold and playful. Eye-catching, personality-forward. Makes pool maintenance feel less like a chore.

## Architecture

Integrated into the existing app server. `GET /` checks auth status:

- **Authenticated users** see the dashboard (no change)
- **Unauthenticated users** see the marketing page

The marketing page has its own navbar: logo/brand on the left, "Log In" and "Sign Up" buttons on the right.

## Page Structure: Single-Page Hero Scroll

### Section 1: Hero

- Full-width dark background (`#13111C`), matching the app navbar
- Large bold headline in white Inter Tight, e.g.: **"Your pool called. It wants better care."**
- Subline in lighter text: "Track water chemistry, schedule maintenance, manage equipment — and actually enjoy pool season."
- Primary CTA: **"Get Started"** button (teal, large, links to `/signup`)
- Secondary: "Already have an account? Log in"
- Below text: a stylized CSS-rendered mockup of the dashboard using the app's design tokens (neumorphic cards, teal accents, status colors) — not a screenshot
- Subtle gradient or wave transition at the bottom blending into the light body (`#f8f7fc`)

### Section 2: Feature Highlights

Five alternating sections (odd: text left / visual right; even: flipped). Each has a playful heading, short description, and a styled HTML/CSS visual using the app's design tokens.

**Feature 1: Water Chemistry**
- Heading: "Know your water like never before"
- Copy: Track pH, chlorine, alkalinity, and more. Out-of-range values get flagged instantly — no more guessing if your water is safe. Get a step-by-step treatment plan telling you exactly what to add and how much.
- Visual: Styled table snippet showing green/red value badges

**Feature 2: Maintenance Tasks**
- Heading: "Never forget to skim again"
- Copy: Set up recurring tasks that reschedule themselves when you mark them done. See what's overdue, what's due today, and what's coming up — all at a glance.
- Visual: Task cards with due-date tags in green/amber/red

**Feature 3: Chemical Inventory**
- Heading: "Running low? You'll know first."
- Copy: Track your chemical stock with quick-adjust buttons. Get alerts before you run out mid-season. No more emergency trips to the pool store.
- Visual: Chemical cards with stock levels and a "Low Stock" badge

**Feature 4: Equipment Tracking**
- Heading: "Your pump has a story. Keep it."
- Copy: Log every piece of equipment with warranty dates and service history. Know exactly when your filter was last serviced and whether that heater is still under warranty.
- Visual: Equipment card with warranty badge and service record entries

**Feature 5: Smart Notifications**
- Heading: "A nudge before things go sideways"
- Copy: Get email or SMS reminders when tasks are due. PoolVibes checks on a schedule so you don't have to.
- Visual: Stylized notification bell or phone mockup with notification bubble

### Section 3: Stats Bar

Horizontal band with a light teal-tinted background. Three bold stats in a row:

- **"6 key parameters"** — tracked with every water test
- **"Zero guesswork"** — treatment plans tell you exactly what to add
- **"Free & open source"** — built with love, no strings attached

Large keywords in Inter Tight bold with small descriptor text below.

### Section 4: Final CTA

- Dark background (`#13111C`) to bookend with the hero (dark-light-dark sandwich)
- Heading: **"Your pool deserves better than a spreadsheet."**
- Subline: "Sign up in seconds. Start tracking today."
- Large teal "Get Started" CTA button
- Secondary "Log in" link

## Visual Design

Carries over the app's existing design system:

- **Colors:** Deep purple-black (`#13111C`) for dark sections, lavender-gray (`#f8f7fc`) for light sections, teal (`#0d9488`) as primary accent
- **Typography:** Inter for body, Inter Tight (600-800 weight) for headings
- **Cards:** Neumorphic shadow style from the app
- **Status colors:** Green (good), amber (attention), red (action needed)
- **Feature visuals:** Built with HTML/CSS using the app's design tokens — not screenshots — so they stay crisp and consistent
- **Responsive:** Full mobile support with Bulma grid, stacked layouts on small screens
- **Dark mode:** Supported via `prefers-color-scheme`

## Primary CTA

"Get Started" links to `/signup` (account creation on this instance).
