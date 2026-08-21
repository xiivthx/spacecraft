# UX / UI catalog

Item SoT for discuss, designer, and browser-probe. **Do not load the whole tree.** Match id(s), Read those files, score each `- [ ]` line. Tips (`💡`) are hints, not gates.

Always-on rules do not load this catalog.

## Lanes

| Lane | Files | Adapter |
|------|-------|---------|
| `/sc-discuss` + `sc-designer` | **One** primary id | `../surface-checklist.md` |
| `/sc-browser-probe` | Inventory (visible surfaces + on-screen chrome) | `../../../sc-browser-probe/references/surface-match.md` |
| `design-system/` | Probe chrome only (not a discuss primary id) | files in this tree |

## Record

```
UX checklist: <id>
UX checklist: none - <reason>
```

`<id>` = alias below, or `category/slug` from the file's `id:` frontmatter.

## Platform pick

1. In-app web → `web-app/`
2. Marketing site → `website/`
3. Native mobile → `mobile-app/`
4. Multi-step interaction → `flow/`
5. Component chrome → `design-system/` (probe, not discuss primary)

Load the variant that matches the product, not every platform file.

## Aliases

| Id | Path | Match when |
|----|------|------------|
| `login` | `web-app/login.md` | Sign-in / session start |
| `sign-up` | `website/sign-up.md` | Create-account / register |
| `form-submit` | `flow/submitting-a-form.md` | Primary job is submit (create, save, contact) |
| `input-error` | `flow/showing-input-error.md` | Validation-heavy form |
| `empty-state` | `web-app/empty-state.md` | Collection first-run / no rows |
| `settings` | `web-app/settings.md` | Preferences / account settings |
| `search-results` | `web-app/search-results.md` | Search or filter results |
| `multi-step` | `web-app/multi-step-form.md` | Wizard / stepped form |
| `save-dirty` | `flow/saving-changes.md` | Edit + save existing record |
| `filter` | `flow/filtering-items.md` | Collection filters |
| `upload` | `flow/uploading-media.md` | File / media upload |
| `auth-recovery` | `flow/resetting-password.md` | Forgot / reset password |
| `destructive` | `flow/deleting-account.md` | Delete account or irreversible destroy |
| `notifications` | `web-app/notifications.md` | In-app notification list |
| `onboarding` | `web-app/onboarding.md` | First-run setup |
| `detail` | `web-app/single-item-detail.md` | Single item detail |
| `not-found` | `website/404.md` | Unknown URL / 404 |
| `nav-tabs` | `mobile-app/tab-bar-navigation.md` | Persistent tab / bottom nav |
| `overlay` | `design-system/modal.md` | Modal, drawer, or action sheet |
| `table` | `design-system/table.md` | Data table |

Any other file: use its `id:` (example `web-app/kanban-board`). `website/pricing` → `website/price.md`.

## Score

Each `- [ ]` title is one item. `n/a` when the product lacks that capability.

- Item marked `(state)` → missing/fail = **critical**
- Else chrome/path → missing/fail = **important**

Discuss/designer emit `present` | `missing` | `n/a`. Probe emits `ok` | `fail` | `n/a` | `deferred`.

Do not use an external checklist site or its AI review as a gate.

## Catalog

112 files: `website/` (23), `web-app/` (27), `design-system/` (29), `mobile-app/` (20), `flow/` (13).

### website/

`404` `about` `affiliate` `billing` `blog` `blog-post` `careers` `cart` `compare` `contact-us` `faq` `features` `login` `press-media` `price` `privacy` `search` `security` `sign-up` `status` `team` `testimonials` `waitlist`

### web-app/

`2fa` `account` `admin-panel` `analytics` `api-keys` `billing` `chat` `comments` `empty-state` `feed` `help-center` `integrations` `kanban-board` `login` `maintenance` `multi-step-form` `notification-settings` `notifications` `onboarding` `pricing` `public-profile` `search-results` `settings` `single-item-detail` `timeline-gantt-view` `user-management` `version-history`

### design-system/

`accessibility` `accordion` `avatar` `badge` `banner` `button` `card` `carousel` `checkbox` `color-system` `date-picker` `drawer` `dropdown-menu` `icon` `input-field` `loading` `modal` `radio` `searchbar` `skeleton` `slider` `spacing-grid` `table` `tabs` `toast` `toggle` `tokens` `tooltip` `typography`

### mobile-app/

`account` `action-sheet` `billing` `camera` `cart` `chat` `checkout` `gesture-navigation` `in-app-browser` `in-app-notifications` `invite` `login` `map-view` `onboarding` `onboarding-checklist` `paywall` `search` `settings` `splash-screen` `tab-bar-navigation`

### flow/

`adding-to-cart` `canceling-subscription` `contacting-support` `deleting-account` `entering-promo-code` `filtering-items` `making-a-card-payment` `resetting-password` `saving-changes` `showing-input-error` `submitting-a-form` `uploading-media` `verifying-account`
