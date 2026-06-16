# UI Design System & Reference — "Luminous Depth"

> Audience: AI agents and developers. This document is the single source of
> truth for the visual/UI layer of `frontend/`. It defines the design language,
> tokens, typography, components, theming, motion, accessibility, and the admin
> design system. It complements `FRONTEND-DESIGN.md` (app architecture/data)
> and `PLAN-DESIGN.md` (migration plan). Where this doc and older docs disagree
> about visuals, THIS doc wins. It supersedes every earlier direction
> (the amber/cyan "Systems/Telemetry" look AND the teal/azure/magenta "aurora
> field" look) entirely.

---

## 1. Design language: the thesis

The portfolio belongs to a **polished, role-neutral senior software engineer**
(backend, fullstack, application security). The UI must read as premium,
modern, and confident — never templated, never a terminal. The direction is
**"Luminous Depth"**:

- **Clean surfaces, luminous accent.** A bright, airy light theme (white page,
  near-white raised surfaces) and a deep-navy dark theme (`#0b1020` base). Soft
  hairline borders, gentle radii, restrained depth.
- **One luminous signature.** A single iridescent gradient — **indigo ->
  violet -> blue** — is the brand. It is **not** a background animation. It is
  spent on gradient **text** (the hero name, stat values, footer heading), thin
  **dividers**, the soft **section glow** behind each title, the **timeline**
  markers/line, the hero **icon chips**, the **logo** wordmark, the
  **scroll-progress** bar, and `.btn-primary`. Everything else stays quiet.
- **Expressive grotesque type.** **Space Grotesk** is the display/heading face;
  **Inter** is the body. **Clash Display** is reserved for ONE thing — the logo
  wordmark. Monospace (**JetBrains Mono**) is demoted to code, dates, counts,
  IDs, and small metadata only — headings are NOT monospace and NOT Clash.
- **Refined structure.** Sections are separated by a **centered, zero-padded
  number kicker** + centered title + glow, reinforced by **faint brand-tinted
  alternating bands**. Eyebrows are quiet uppercase tracked labels (no `//`,
  no `~/path`, no `[status]`, no terminal flavor).
- **A dynamic hero, not a busy page.** The hero's moving element is a cluster of
  three **floating role cards** (Backend / Security / Fullstack), not a canvas
  field or a code/box panel.

### Anti-patterns (do NOT reintroduce)

- No **teal / azure / magenta** signature, and no amber `#f5a623` / cyan
  `#57c4c2` two-signal palette. Indigo/violet/blue is now THE palette — do not
  treat indigo `#6366f1` / violet `#8b5cf6` as forbidden (an earlier doc did;
  that rule is reversed).
- No **aurora field** / animated `<canvas>` background, no `AuroraField`, no
  `NodeNetwork`, no `ServiceGraph` / wireframe topology, no three.js.
- No schematic **grid or dot background** on `body`. The page background is a
  flat token color; section rhythm comes from numbered titles + glow + bands.
- No **node-network / box / terminal-window** hero visual.
- No **JetBrains Mono headings**; no terminal eyebrows (`//`, `~/`, `edu.01`).
- No **blinking-caret logo**; the logo is a gradient wordmark only.
- No emoji as UI icons (use the inline SVGs in `src/icons/` or stroke SVGs).
- No `glightbox` (replaced by the in-house `Lightbox.svelte`).
- No hardcoded hex in components — use tokens / Tailwind utilities. Only
  translucent shadows/scrims (`rgba(0,0,0,.x)`) and the error red `#e5614b`
  are exempt.

---

## 2. How styling works (Tailwind-first)

Order of preference:

1. **Tailwind v4 utilities in markup** are the default. Color/spacing/layout/
   typography are expressed as utilities (including arbitrary values like
   `min-[980px]:grid-cols-[1.05fr_0.95fr]`). Per-component scoped `<style>`
   blocks are not used for layout.
2. **Component recipes** in `src/styles/global.css` `@layer components`
   (`.eyebrow`/`.eyebrow-accent`, `.aurora-text`, `.aurora-rule`,
   `.section-glow`, `.card`/`.card-hover`, `.panel`, `.btn-primary`,
   `.btn-ghost`, `.chip`). Reuse these for repeated patterns instead of
   re-deriving them.
3. **Tokens** (CSS custom properties) are the source of truth for color, font,
   and radius. They are declared in `src/layouts/Layout.astro`
   `<style is:global>` on `:root` (light) and `:root[data-theme="dark"]` (dark),
   and bridged into Tailwind via `@theme inline` in `global.css`.
4. **Scoped `<style>` is allowed only** in three genuinely special leaves where
   utilities can't express the effect cleanly: `Logo.svelte`,
   `icons/Logo.astro` (the shimmering wordmark, self-contained so it also runs
   in `/admin`), and `Skills.astro` (a `.skill-icon :global(svg)` recolor rule
   that forces inline SVGs to `var(--color-icon)` while preserving
   `svg[fill="none"]` outline icons). Everything else is Tailwind-first.

Stack: **Astro 6** (SSR, `output: "server"`, Node standalone) renders pages and
static sections; **Svelte 5** (runes) powers interactive islands; **Tailwind
v4** via `@tailwindcss/vite`. Fonts: Fontshare (Clash Display) + Google (Inter +
Space Grotesk + JetBrains Mono) `<link>`s in `Layout.astro` (and mirrored in
`AdminLayout.astro`).

---

## 3. Design tokens (public site)

Declared in `src/layouts/Layout.astro` -> `<style is:global>`. **Light is
`:root`; dark overrides under `:root[data-theme="dark"]`.**

| Token                   | Light                   | Dark                 | Meaning                                      |
| ----------------------- | ----------------------- | -------------------- | -------------------------------------------- |
| `--grad-1`              | `#6366f1`               | `#818cf8`            | Gradient stop 1 (indigo)                     |
| `--grad-2`              | `#8b5cf6`               | `#a78bfa`            | Gradient stop 2 (violet)                     |
| `--grad-3`              | `#3b82f6`               | `#60a5fa`            | Gradient stop 3 (blue)                       |
| `--color-primary`       | `#6366f1`               | `#818cf8`            | Indigo — primary actions/links               |
| `--color-primary-dark`  | `#4f46e5`               | `#6366f1`            | Hover/darker primary                         |
| `--color-secondary`     | `#8b5cf6`               | `#a78bfa`            | Violet — accent, eyebrow-accent, "completed" |
| `--color-accent`        | `#3b82f6`               | `#60a5fa`            | Blue — sparing pops                          |
| `--color-text`          | `#1f2937`               | `#e5e7eb`            | Primary text                                 |
| `--color-text-light`    | `#6b7280`               | `#9ca3af`            | Muted/secondary text                         |
| `--color-bg`            | `#ffffff`               | `#0b1020`            | Page background                              |
| `--color-bg-alt`        | `#f9fafb`               | `#0f172a`            | Chips, badges, bands, raised bars            |
| `--color-surface`       | `#ffffff`               | `#111827`            | Card/panel surface                           |
| `--color-surface-2`     | `#f3f4f6`               | `#1f2937`            | Raised surface                               |
| `--color-border`        | `#e5e7eb`               | `#1f2937`            | Hairline borders                             |
| `--color-border-strong` | `#d1d5db`               | `#374151`            | Stronger divider/scrollbar                   |
| `--color-icon`          | `#1f2937`               | `#e5e7eb`            | Skill icon color                             |
| `--nav-surface`         | `rgba(255,255,255,.85)` | `rgba(17,24,39,.85)` | Scrolled navbar bg (blurred)                 |

Footer tokens (separate dark-on-light/dark band):

| Token                         | Light     | Dark      | Meaning                       |
| ----------------------------- | --------- | --------- | ----------------------------- |
| `--footer-grad-start`         | `#1e1b4b` | `#0b1020` | Footer background             |
| `--footer-grad-end`           | `#312e81` | `#151a2f` | Footer gradient end (defined) |
| `--footer-text`               | `#ffffff` | `#e5e7eb` | Footer text                   |
| `--footer-subtle`             | `#a5b4fc` | `#9ca3af` | Footer muted text             |
| `--footer-heading-grad-start` | `#ffffff` | `#e5e7eb` | Footer heading gradient start |
| `--footer-heading-grad-end`   | `#c7d2fe` | `#a5b4fc` | Footer heading gradient end   |

Non-color tokens (theme-independent):

| Token            | Value                                            | Meaning                        |
| ---------------- | ------------------------------------------------ | ------------------------------ |
| `--font-sans`    | `"Inter", system-ui, …`                          | Body font                      |
| `--font-display` | `"Space Grotesk", system-ui, sans-serif`         | Display/heading font           |
| `--font-logo`    | `"Clash Display", "Space Grotesk", system-ui, …` | **Logo wordmark ONLY**         |
| `--font-mono`    | `"JetBrains Mono", ui-monospace, …`              | Code + small metadata          |
| `--radius`       | `14px`                                           | Standard card/panel radius     |
| `--radius-sm`    | `8px`                                            | Chips, badges, buttons, inputs |
| `--radius-lg`    | `20px`                                           | Panels, modals                 |
| `--transition`   | `all .3s cubic-bezier(.4,0,.2,1)`                | Default transition             |

### Tailwind utilities exposed (via `@theme inline`)

`bg-bg`, `bg-bg-alt`, `bg-surface`, `bg-surface-2`, `bg-primary`,
`bg-secondary`, `bg-accent`; `text-text`, `text-text-light`, `text-primary`,
`text-secondary`, `text-accent`; `border-border`, `border-border-strong`;
`text-grad-1/-2/-3`; `font-sans`, `font-display`, `font-mono`; `rounded`
(=`--radius`), `rounded-sm`, `rounded-lg`. (`--font-logo` is NOT bridged into a
utility — it is consumed only by the logo's scoped `<style>`.) Dark mode keeps
working through utilities because they resolve to the live variables.

---

## 4. Typography

- **Display / headings:** Space Grotesk (`--font-display`), weight **700**. All
  `h1–h6` use it globally (set in `global.css @layer base`), `line-height 1.2`,
  `letter-spacing -0.02em`, `text-wrap: balance`.
  - `h1`: `clamp(2.6rem, 6vw, 4.5rem)`
  - `h2`: `clamp(2rem, 4vw, 3rem)` (Section overrides to `clamp(1.9rem,4vw,2.6rem)`)
  - `h3`: `clamp(1.4rem, 3vw, 1.9rem)`
- **Body:** Inter (`--font-sans`). `p` is `--color-text-light`, `1rem`,
  `line-height 1.7`, `text-wrap: pretty`. `html` carries `letter-spacing
-0.011em`.
- **Logo:** Clash Display (`--font-logo`), weight 700 — used ONLY by the
  `.logo` wordmark. Never for headings or body.
- **Mono:** JetBrains Mono (`--font-mono`) — code, dates, counts, IDs, small
  tags/labels only. Never headings or the hero name.
- **Gradient type:** apply `.aurora-text` to a heading/word to fill it with the
  drifting indigo->violet->blue gradient (the hero name, About stat values,
  footer heading).

### Type roles cheat-sheet

| Use                 | Class/Font                                                          |
| ------------------- | ------------------------------------------------------------------- |
| Section/hero titles | h2/h1 (Space Grotesk); `.aurora-text` for the signature word        |
| Logo wordmark       | `.logo` (Clash Display, gradient, shimmer)                          |
| Eyebrows / labels   | `.eyebrow` (uppercase, tracked, muted) / `.eyebrow-accent` (violet) |
| Dates, counts, tags | `font-mono`, muted                                                  |
| Body copy           | Inter, `text-text-light`                                            |

---

## 5. Logo & iconography

### Logo — gradient wordmark `JPC`

Two mirrored implementations: `src/icons/Logo.astro` (Astro) and
`src/components/Logo.svelte` (Svelte, used in the nav island and `/admin`). The
text `JPC` is set in **Clash Display 700 via `--font-logo`** (this is the lone
consumer of that token) and filled with the luminous gradient
(`--grad-1/-2/-3`) which **shimmers** via the `logo-shimmer` keyframe (8s
linear; speeds to 2.4s + `brightness(1.18)` on hover). Each implementation
carries its OWN scoped `@keyframes logo-shimmer` so the effect runs in `/admin`
(which doesn't load `global.css`). Theme-aware (reads the gradient tokens).
Motion stills under `prefers-reduced-motion`. Sized by the consumer via
`font-size` on `class`. Accessible name "Juan Pablo Cano". No SVG, no caret.

### Icons — `src/icons/*.astro`

~30 inline SVG components (tech/brand icons + Mail, LinkedIn, GitHub, Phone,
etc.). They render inline so CSS recolors them via `currentColor` /
`var(--color-icon)`. Inline stroke SVGs scattered through markup (close, arrows,
sun/moon, view, server/shield/code chips, graduation cap, check) use
`stroke="currentColor"`.

---

## 6. The signature: luminous gradient + glow (NOT a background)

The single bold element is the **indigo -> violet -> blue gradient**, applied to
foreground accents — there is no animated background and no `<canvas>`. Where it
shows up:

- **`.aurora-text`** — gradient-clipped, drifting text (`aurora-text-drift`,
  9s, `linear-gradient(100deg, grad-1 0%, grad-2 45%, grad-3 90%)`,
  `background-size: 220% 100%`). Used on the hero name, About stat values, and
  the footer "Juan Pablo Cano" heading.
- **`.aurora-rule`** — 2px gradient divider
  (`linear-gradient(90deg, grad-1, grad-2 40%, grad-3 70%, transparent)`),
  `opacity .85`. Under section titles (constrained `max-w-[80px] mx-auto`) and
  along the top of the footer.
- **`.section-glow`** — a blurred radial brand glow centered behind each section
  title (`width: min(420px,80%)`, `height: 220px`,
  `transform: translate(-50%,-62%)`, `filter: blur(44px)`, `opacity .55`,
  radial of violet 30% -> indigo 16% -> transparent). Adds depth without a busy
  background.
- **Section bands** — `main > section.content-section:nth-of-type(even)` gets
  `color-mix(in srgb, var(--grad-2) 6%, var(--color-bg-alt))`: a faint violet
  wash so alternating sections read as distinct blocks.
- **Timeline** (Experience), **icon chips** (Hero / Education), the **logo**,
  the **ScrollProgress** bar, and **`.btn-primary`** all use the gradient too.

---

## 7. Layout primitives & global UI

### `Layout.astro` (public shell)

`<head>`: charset, title, `<ClientRouter />`, `<ThemeProvider />`, meta,
Fontshare (Clash Display) + Google (Inter / Space Grotesk / JetBrains Mono)
font links, then `@/styles/global.css`. `<body>`: just `<slot />` (no background
field wrapper). The `<style is:global>` holds ONLY the token `:root` /
`[data-theme="dark"]` blocks. Base/reset/heading rules and the alternating-band
rule live in `global.css @layer base`.

### `Section.astro` (feature-section wrapper)

Props `{ title?, id?, index? }`. Tailwind-first. Renders
`section.content-section.fade-in#id` (`relative px-5 py-16 md:px-8 md:py-24`) ->
centered container (`mx-auto w-full max-w-[1180px]`). When `title` is set, a
**centered** `header` (`relative mb-10 text-center md:mb-14`) contains, in z
order:

1. `<span class="section-glow">` (the blurred brand glow, `aria-hidden`),
2. when `index != null`, a **zero-padded number kicker**
   `<span class="eyebrow eyebrow-accent ...">{"01"}</span>`
   (`String(index).padStart(2,"0")`),
3. a centered `<h2 class="text-[clamp(1.9rem,4vw,2.6rem)]">{title}`,
4. a centered `<span class="aurora-rule ... mx-auto mt-6 block max-w-[80px]">`.

Then `<slot />`. The section is **not** transparent and there is no per-id
eyebrow text. Section numbering (via `index`): **About = 1, Experience = 2,
Education = 3, Projects = 4, Skills = 5, Certifications = 6** (Hero is not a
`Section` and has no number).

### Breakpoints

Tailwind defaults (mobile-first; `md` = 768px, `lg` = 1024px) plus arbitrary
stops used in markup where the layout needs them — notably `min-[980px]:` /
`max-[980px]:` (Hero two-column collapse), `max-[820px]:` (nav off-canvas), and
`max-[768px]:`.

### Container widths

Public content max-width **1180px**; navbar inner container **1320px**.

### Z-index ladder

ScrollProgress `9999`; Navbar `1000`; Modal/Lightbox overlays `1000`; `main`
`z-10`; admin sidebar `100`; admin header `50`; admin overlay `99`. (No negative
background layer — the field is gone.)

---

## 8. Theming mechanism (unchanged contract)

- Source of truth: `localStorage["theme"]` (`"light"|"dark"`), falling back to
  OS `prefers-color-scheme`; applied as `<html data-theme="…">`.
- **No FOUC:** `ThemeProvider.astro` runs a render-blocking inline `<script>` in
  `<head>` that resolves/applies the theme before paint, manages
  `<meta name="color-scheme">`, exposes `toggle()`, and dispatches/listens to
  events.
- **Event bus:** `theme:toggle` (request a flip, from `ThemeToggle.svelte`);
  `theme:changed` (emitted after a flip; listened to by `ThemeToggle.svelte` to
  re-read the resolved theme). Cross-tab via the `storage` event.
- **Admin is exempt:** admin pages don't load `ThemeProvider`; `admin.css`
  defines its own fixed dark tokens under `.admin-root`.

---

## 9. Component recipes (reuse these)

From `global.css @layer components` — all render in the indigo/violet/blue
palette:

- **`.eyebrow`** / **`.eyebrow-accent`** — uppercase, tracked (`.2em`), 600,
  muted label; `-accent` recolors to `--color-secondary` (violet). Replaces all
  terminal eyebrows; also used as the section number kicker.
- **`.aurora-text`** — luminous gradient-clipped text that drifts (stills under
  reduced motion).
- **`.aurora-rule`** — 2px luminous gradient divider. Make it `block` and
  constrain width (e.g. `max-w-[80px]`).
- **`.section-glow`** — absolutely-positioned blurred radial brand glow for
  section-title depth (see §6).
- **`.card`** + **`.card-hover`** — surface + hairline + `--radius`; `-hover`
  adds a `-4px` lift, a primary-tinted border, and a layered shadow on hover.
- **`.panel`** — glassy translucent panel (`backdrop-filter: blur(14px)`) for
  the ProjectModal and Lightbox dialogs.
- **`.btn-primary`** — luminous **violet->blue** gradient fill
  (`linear-gradient(100deg, grad-2, grad-3)`), white text, Space Grotesk,
  position-shift + lift on hover.
- **`.btn-ghost`** — bordered translucent control (Space Grotesk 500); hover ->
  primary border+text + lift. Used by ThemeToggle and the nav "cv" button.
- **`.chip`** — small mono metadata chip (`bg-alt`, hairline, muted).

Keyframes available globally (in `global.css`): `aurora-text-drift`,
`logo-shimmer`, `pulse-dot` (status dots), `nudge` (scroll cue), and the v1
hero set — **`float`** (floating cards), **`slide-in-left`** (hero text
entrance), **`fade-in-soft`** (hero visual entrance).

---

## 10. Public component catalog (current look)

Composed in `src/pages/index.astro`:
`Hero -> About -> Experience -> Education -> Projects -> Skills ->
Certifications`, inside `<main class="relative z-10">`, with
`<ScrollProgress client:load>`, `<Navigation>`, and `<Footer>`. A page script
adds `.visible` to `.fade-in` on scroll via `IntersectionObserver` (re-runs on
`astro:after-swap`).

- **Navigation** (`Navigation.astro` -> `NavigationBar.svelte`): fixed bar,
  inner `max-w-[1320px]`. Left: `<Logo>` (Clash wordmark). Center: links
  (About / Experience / Projects / Skills) in **Inter `text-[1rem]` medium**,
  `text-text-light` -> `text-text` on hover, plus a **subtle outlined**
  "Contact" CTA (`border border-primary/45 text-primary hover:bg-primary/10`,
  `rounded-[var(--radius-sm)]`) — **not** `.btn-primary`. Right: `<ThemeToggle>`
  (`.btn-ghost`, ~40px square), a `.btn-ghost` "cv" download, and a hamburger
  button. `scrolled` (>80px) state adds `--nav-surface` blur + bottom hairline
  and tightens padding. Mobile off-canvas menu at `<=820px`.
- **ThemeToggle.svelte:** `.btn-ghost` 40×40 (`w-10 h-10 p-0`) with sun/moon
  stroke SVG; dispatches `theme:toggle`; mirrors theme via `theme:changed` /
  `storage`; `aria-pressed` + `aria-label`.
- **ScrollProgress.svelte:** fixed 3px top bar (`z-[9999]`); inner bar width =
  scroll %, filled with `linear-gradient(90deg, grad-1, grad-2, grad-3)`.
- **Hero.astro:** full-height section, **no panel and no field**. Two-column
  grid `grid-cols-[1.05fr_0.95fr]` that collapses to one column `<=980px`.
  - **Left** (`animate-[slide-in-left_0.7s_ease]`): `.eyebrow.eyebrow-accent`
    reading **"Backend · Security · Fullstack"**, an `h1` name with
    **`.aurora-text`**, a `text-secondary` role label, a short description, a
    mono definition list bordered-left (`Role` / `Stack` / `Based`), then
    `<ContactLinks>`.
  - **Right** (`animate-[fade-in-soft_0.8s_ease_0.4s_backwards]`): three
    **floating `.card`s**, absolutely positioned on desktop and
    `static flex-wrap` `<=980px`. Each has a gradient **icon chip** (rounded
    `bg-gradient-to-br`, white inline stroke SVG — server / shield / code) plus
    a `font-display` label: **Backend**, **Security**, **Fullstack**. Each
    floats via `animate-[float_3s_ease-in-out_infinite]` with delays
    `0` / `1s` / `0.5s`.
  - Bottom: a mono "scroll" cue with a `nudge`-animated down arrow (shown
    `min-[980px]:flex`).
- **About.astro** (`index={1}`): centered block (`mx-auto max-w-[920px]
text-center`) — centered summary paragraph, then a stat grid
  (`auto-fit minmax(220px,1fr)`) of `.card .card-hover` tiles, each a big
  **`.aurora-text`** value over a mono lowercase label.
- **Experience.astro** (`index={2}`): SSR `api.get<Experiences>`; `<ApiError>`
  then a **vertical timeline** spanning full width. Each row is a grid
  `grid-cols-[34px_1fr]` (`md:grid-cols-[50px_1fr]`): a gradient **dot marker**
  (`bg-gradient-to-br from-primary to-secondary`, ringed via box-shadow,
  scales on group-hover) above a gradient **connecting line** (hidden on the
  last item), beside a `.card .card-hover` with a company link, position,
  `location • type`, a **rounded-full date pill** (`bg-bg-alt`, start–end
  years), description, and a stretched **"View details ->"** link to
  `/experience/{id}`.
- **Education.astro** (`index={3}`): full-width grid
  `auto-fit minmax(min(100%,300px),1fr) gap-8` of `.card .card-hover` articles.
  Each opens with a gradient **graduation-cap icon chip** (`h-12 w-12`,
  `bg-gradient-to-br from-primary to-secondary`, scales on group-hover), the
  institution (optional external link), `"{studyType} in {area}"`, and a year
  **badge pill** (`text-primary bg-primary/10 rounded-full`).
- **Projects.astro + ProjectModal.svelte** (`index={4}`): `<ApiError>` then a
  `<ul>` grid `auto-fit minmax(min(320px,100%),1fr)` of `.card .card-hover`
  articles (`role="button"`, `tabindex="0"`) carrying a `data-project` JSON
  blob. Each shows a title, a `.chip` status — **active** = `text-primary` + a
  pulsing dot (`pulse-dot`) / **completed** = `text-secondary` + a check SVG —
  a clamped description, and `.chip` tech tags. A vanilla script dispatches a
  `project:open` CustomEvent on click/Enter/Space (re-init on
  `astro:after-swap`). `ProjectModal.svelte` (`client:idle`): one shared modal
  listening for `project:open`; scrim `rgba(0,0,0,.6)` + blur, a `.panel`
  dialog (`role="dialog" aria-modal`), title + status chip, description, tech
  chips, a started/completed mono block, and an optional `.btn-primary` "Visit
  Project"; closes on close button / backdrop / Escape, with body scroll lock
  and reset on `astro:after-swap`.
- **Skills.astro** (`index={5}`): grid `auto-fit minmax(min(280px,100%),1fr)` of
  4 hardcoded category `.card`s. Each card header has an `.eyebrow` title + a
  mono count; the body is a flex-wrap list of pill items
  (`bg-bg-alt`, hairline, hover -> primary border + lift) with a `.skill-icon`
  (inline SVG forced to `var(--color-icon)` via the one allowed scoped
  `:global(svg)` rule, preserving `svg[fill="none"]` outlines) + a mono name.
- **Certifications.astro + Lightbox.svelte** (`index={6}`): SSR
  `api.get<CareerCertifications>` mapped to `{ url, title, issuer, date }` ->
  `<Lightbox client:visible>`. The lightbox renders the cert-card grid
  (`auto-fill minmax(min(300px,100%),1fr)`; `.card .card-hover` with image +
  hover/focus "view" overlay chip + title/issuer/date) and a full-screen
  `.panel` dialog (image + prev/next + counter). Fully accessible:
  `role="dialog" aria-modal`, focus trap, Esc/Arrows, backdrop close, scroll
  lock, focus restore; resets on `astro:after-swap`.
- **Footer.astro + ContactLinks.astro:** `footer#contact`, background
  `var(--footer-grad-start)`, with an `.aurora-rule` along the top edge. Three
  columns (Contact: `.eyebrow` + `.aurora-text` name + bio; Navigation links;
  Connect = `<ContactLinks>`) over a hairline copyright row (mono).
  `ContactLinks` = a mailto chip + social icon links — 40×40 `bg-surface`
  hairline squares, hover -> primary border/text + lift.
- **experience/[id].astro:** SSR `GET /experiences/:id`, `zod safeParse` with
  `experienceDetailSchema`, rendered in `Layout` with `<Navigation>` +
  `<Footer>`. Mono back-link, detail header (eyebrow `location • type`, `h1`,
  `text-secondary` company link, mono period, summary), then client `.card`s
  (sorted by client `startDate`) with name + `.chip` period, description, and
  Responsibilities / Achievements bullet lists (secondary `▹` markers via
  `before:`) + `.chip` tech tags.
- **ApiError.astro:** shared error/empty state. Error -> soft red panel
  (`border-[#e5614b]/40 bg-[#e5614b]/10`; `#e5614b` is the one allowed semantic
  red), raw details only when `ENVIRONMENT === "local"`; empty -> muted "No data
  to display yet."

---

## 11. Admin design system (`/admin`)

A **self-contained dark "work mode"** in the same Luminous Depth language,
denser. It does NOT use `ThemeProvider`; `src/styles/admin.css` declares its own
tokens under `.admin-root` (imported by `AdminLayout.astro`). **Class names are
unchanged** so the wired Svelte managers need no logic changes — this is a
repaint only, to the indigo/violet/blue palette.

### Admin tokens (under `.admin-root`, fixed dark)

| Token                    | Value                                                       | Meaning                        |
| ------------------------ | ----------------------------------------------------------- | ------------------------------ |
| `--admin-bg`             | `#0b1020`                                                   | App background (deep navy)     |
| `--admin-surface`        | `#111827`                                                   | Cards, sidebar, header         |
| `--admin-surface-2`      | `#1f2937`                                                   | Raised surface                 |
| `--admin-border`         | `#1f2937`                                                   | Hairlines                      |
| `--admin-border-strong`  | `#374151`                                                   | Stronger divider               |
| `--admin-text`           | `#e5e7eb`                                                   | Text                           |
| `--admin-text-muted`     | `#9ca3af`                                                   | Muted text                     |
| `--admin-text-faint`     | `#6b7280`                                                   | Faint text                     |
| `--admin-primary`        | `#818cf8`                                                   | Indigo — primary buttons/links |
| `--admin-primary-hover`  | `#a5b4fc`                                                   | Hover                          |
| `--admin-primary-ink`    | `#0b1020`                                                   | Ink on primary                 |
| `--admin-secondary`      | `#a78bfa`                                                   | Violet accent                  |
| `--admin-accent`         | `#60a5fa`                                                   | Blue accent                    |
| `--admin-danger`         | `#e5614b`                                                   | Destructive                    |
| `--admin-success`        | `#34d399`                                                   | Success                        |
| `--admin-warning`        | `#fbbf24`                                                   | Warning                        |
| `--admin-input-bg`       | `#0d1426`                                                   | Input background               |
| `--admin-input-border`   | `#2a3550`                                                   | Input border                   |
| `--admin-primary-subtle` | `rgba(129,140,248,0.16)`                                    | Focus ring tint                |
| `--admin-radius` / `-sm` | `14px` / `8px`                                              | Radius                         |
| layout                   | `--admin-sidebar-width 260px`, `--admin-header-height 64px` |                                |

`--admin-gradient` = `linear-gradient(100deg, secondary 0%, primary 50%, accent
100%)` = **violet -> indigo -> blue**, used on `.admin-btn-primary`, the header
title tick, stat-card top borders, and the login card top stripe. The shared
`<Logo>` works in admin because `--grad-1/-2/-3` are aliased to
secondary/primary/accent. **Fonts** (Inter / Space Grotesk / Clash logo /
JetBrains Mono) and legacy aliases (old `--admin-amber*` / `--admin-cyan*`
names) are re-declared under `.admin-root` and map onto the new tokens so
existing markup keeps working.

### Admin class vocabulary (reuse; don't reinvent)

Shell: `.admin-root`, `.admin-sidebar` (`-header/-logo/-nav/-footer`),
`.admin-nav-section`/`-title`, `.admin-nav-link`(`.active`),
`.admin-user-info`(`-avatar/-email`), `.admin-main`, `.admin-header`(`-title`),
`.admin-content`. Cards/stats: `.admin-card`(`-header/-title`),
`.admin-stats-grid`, `.admin-stat-card`(`-label/-value`), `.admin-stat-link`,
`.admin-quick-actions`/`-action`(`-icon/-content`). Buttons: `.admin-btn`
(`-primary/-danger/-ghost/-sm/-icon`). Forms: `.admin-form`,
`.admin-form-group/-row`, `.admin-label`(`-required`), `.admin-input`,
`.admin-textarea`, `.admin-select` (+ `-error/-success`), `.admin-error-text`,
`.admin-help-text`, `.admin-input-wrapper/-icon`, `.admin-form-actions`,
`.admin-array-row`, `.admin-link`. Tables: `.admin-table-container`,
`.admin-table`, `.admin-table-actions`, cell helpers
(`.admin-cell-sub/-ellipsis/-period/-tech`). Modal: `.admin-modal-overlay`,
`.admin-modal`(`-header/-title/-close/-body/-footer`), `.admin-confirm-text`
(`-warning`). Feedback: `.admin-alert`(`-error/-success`), `.admin-badge`
(`-primary/-success/-warning/-secondary`), `.admin-spinner`,
`.admin-empty-state`, `.admin-upload-zone` (+ `-label/-icon/-text/-hint`),
`.admin-login-*`, mobile drawer (`.admin-menu-toggle`, `.admin-sidebar-overlay`,
`<=768px`).

### Admin components

`AdminLayout.astro` (sidebar + main + mobile drawer, `robots: noindex`, its own
font links), `Modal.astro` (`window.openModal/closeModal`, Escape, re-init on
`astro:after-swap`), and Svelte managers (`ExperiencesManager.svelte`,
`ProjectsManager.svelte`, `ExperienceClientsEditor.svelte`, `client:load`,
logic untouched). Pages: `admin/index` (stat/quick-action tiles), `admin/login`
(focused card), `admin/experiences` (+ `[id]` clients), `admin/projects`,
`admin/certifications`.

---

## 12. Motion & accessibility rules

- **Reduced motion:** a global rule in `global.css @layer base` neutralizes
  animations/transitions and pins `.fade-in` visible; `admin.css` has its own
  equivalent under `.admin-root`. `.aurora-text` and the logo shimmer also stop
  explicitly. JS-driven motion (modal/lightbox transitions) is short and
  CSS-driven, so it collapses too. Add the same handling for any new JS motion.
- **Focus:** global `:focus-visible { outline: 2px solid var(--color-primary) }`
  — don't remove it. Modal/Lightbox implement focus traps + restore on close.
- **Keyboard:** project cards handle Enter/Space; overlays handle Escape;
  the lightbox handles Arrow keys.
- **Semantics:** dialogs use `role="dialog" aria-modal="true"` + label;
  icon-only controls have `aria-label`; the logo has an accessible name;
  decorative glow/arrows are `aria-hidden`.
- **Color contrast:** primary/secondary/accent are tuned per theme (deeper on
  light, brighter on dark) — keep that split.

---

## 13. Do / Don't

DO:

- Style with Tailwind utilities + the `@layer components` recipes; use tokens
  for any custom color/radius/font.
- Keep Space Grotesk for headings, Inter for body, Clash Display for the logo
  only, mono for code/metadata only.
- Spend boldness only on the indigo->violet->blue signature (gradient text,
  dividers, section glow, timeline, icon chips, logo, `.btn-primary`); keep
  everything else quiet.
- Separate sections with the centered number kicker + glow + alternating band
  (don't hand-roll new section chrome).
- Add `prefers-reduced-motion` handling + visible focus to anything interactive.
- Put new interactivity in a Svelte island with proper `$effect` teardown.

DON'T:

- Reintroduce teal/azure/magenta or amber/cyan, the aurora field / any animated
  background, a grid or dot background, a node-network / box / terminal hero,
  mono headings, terminal eyebrows, the ServiceGraph, the caret logo, three.js,
  or glightbox.
- Add per-component scoped `<style>` for layout (only the 3 leaves — `Logo`
  ×2 and `Skills` — keep one). Don't hardcode hex (except shadow/scrim rgba and
  the `#e5614b` error).
- Break the cookie/SSR/data flow or the event contracts (`theme:*`,
  `project:open`, `window.openModal`) when restyling.

---

## 14. File map (UI-relevant)

```
src/
├── layouts/
│   ├── Layout.astro            # tokens (incl. --font-logo), fonts, ClientRouter, ThemeProvider
│   └── AdminLayout.astro       # admin shell (imports admin.css; own font links)
├── styles/
│   ├── global.css              # @import tailwindcss + @theme bridge + @layer base/components
│   │                           #   (alternating bands, .section-glow) + keyframes
│   │                           #   (float / slide-in-left / fade-in-soft / aurora-text-drift / …)
│   └── admin.css               # admin design system (.admin-* tokens + classes)
├── components/
│   ├── Logo.svelte / icons/Logo.astro  # gradient wordmark "JPC" via --font-logo (mirrored)
│   ├── Navigation.astro / NavigationBar.svelte   # subtle-outline Contact CTA, 1320px inner
│   ├── ThemeToggle.svelte / ScrollProgress.svelte
│   ├── ProjectModal.svelte / Lightbox.svelte     # .panel dialogs
│   ├── Section.astro           # centered numbered title + glow + aurora-rule (index 1..6)
│   ├── Footer.astro / ContactLinks.astro / ApiError.astro
│   ├── ThemeProvider.astro
│   └── admin/                  # Modal.astro, ExperiencesManager / ProjectsManager / ExperienceClientsEditor (.svelte)
├── features/                   # Hero (floating cards), About, Experience (timeline),
│                               #   Education (card grid), Projects, Skills, Certifications (.astro)
├── icons/                      # ~30 inline SVG icon components (+ Logo.astro)
└── pages/
    ├── index.astro             # composes the home sections + IntersectionObserver reveal
    ├── experience/[id].astro   # experience detail page
    └── admin/**                # admin pages
```

### Removed (do not resurrect)

`AuroraField.svelte` and `NodeNetwork.svelte` (created then deleted — neither
exists), `ServiceGraph.svelte`, the schematic grid/dot background, the
blinking-caret logo, the teal/azure/magenta AND amber/cyan token sets,
JetBrains-Mono headings, and terminal eyebrow devices. (`three`,
`@types/three`, `glightbox`, and the `logo-*.webp` assets are gone.)
