---
name: sc-planning
description: Convert a mission spec into a small executable plan with verifiable tasks
license: MIT
compatibility: opencode
---
- Create `plan.json`-ready output.
- Before producing `plan.json`, check `questions.md` and `decisions.md`.
- Do not fill gray areas with hidden assumptions.
- Record low-risk assumptions explicitly.
- Blocking product/design decisions require sc-clarify.
- Keep tasks <= 7.
- Each task needs `id`, `title`, `status`, `files`, `acceptance`, `verify`, and `evidence`.
- Status values: `pending`, `in-progress`, `done`, `blocked`.
- Avoid vague tasks.
- Include exact files when known.
- Prefer focused tests and build checks.
- Do not create broad architecture plans unless the mission explicitly requires it.
