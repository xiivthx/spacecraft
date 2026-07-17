---
name: sc-pathfinder
description: "Chart a map of tickets for large, multi-session work. Invoke on \"pathfinder\", \"chart a map\", or \"too big for one session\". Not auto-triggered."
---

# sc-pathfinder

Chart a path through a large, foggy effort by creating a shared map of investigation tickets on the repo's issue tracker, then resolving them one at a time until the destination is clear. Pathfinding is about finding the way - not charging at the destination.

## When to use

Activate when the user explicitly invokes:

- **"pathfinder" / "chart a map"** - a loose idea too big for one session
- **"too big for one session"** - work that spans multiple agent invocations
- **Resume a map** - continue working an existing `wayfinder:map` issue

## Core concepts

### The map
A single issue on the issue tracker labeled `wayfinder:map`. It is an **index**, not a store - it lists decisions and links to the child tickets that hold the detail.

**Map body structure:**
```markdown
## Destination
<one or two lines - what reaching the end of this map looks like>

## Notes
<domain; skills to consult; standing preferences>

## Decisions so far
- [<closed ticket title>](link) - <one-line gist of the answer>

## Not yet specified
<in-scope fog that can't be ticketed yet; graduates as the frontier advances>

## Out of scope
<work ruled beyond the destination; never graduates>
```

### Tickets
Each ticket is a **child issue** of the map. Body is exactly:
```markdown
## Question
<the decision or investigation this ticket resolves>
```

Each ticket carries a `wayfinder:<type>` label - one of:
- **research** (AFK) - reading docs, APIs, or knowledge bases. Produces a linked summary.
- **prototype** (HITL) - cheap concrete artifact to react to (outline, stub, rough UI). Links the artifact.
- **grilling** (HITL) - human conversation, one question at a time. Default case.
- **task** (AFK or HITL) - manual work that unblocks a decision (signing up, provisioning, moving data).

**HITL** = human in the loop, worked with a human who speaks for themselves.
**AFK** = agent-driven alone.

### Claiming, blocking, and the frontier
- **Claim** a ticket by assigning it before any work - concurrent sessions skip assigned tickets.
- **Blocking** uses the tracker's native dependency relationship. A ticket is unblocked when every blocker is closed.
- The **frontier** is open, unblocked, unclaimed children - the edge of the known.

### Fog of war
Don't chart what you can't yet see. Beyond the frontier lies fog - decisions you sense but can't yet pin down. Resolving a ticket clears fog ahead. The ticket/fog test: can you state the question precisely *now*? Ticket if yes; **Not yet specified** if no.

### Plan, don't do
Pathfinding is planning by default - produce decisions, not deliverables. The pull to "just do the work" signals you've reached the edge of the map.

## Workflow

### Chart a map
User provides a loose idea.

1. **Name the destination.** Pin down what this map is finding its way to - the spec, decision, or change. The destination fixes the scope.
2. **Map the frontier.** Fan out breadth-first across the space, surfacing open decisions and first steps. **If no fog surfaces** - the way is already clear - stop. Don't create a map. Ask how to proceed.
3. **Create the map** (label `wayfinder:map`): Destination and Notes filled in, Decisions-so-far empty, fog sketched into Not yet specified.
4. **Create tickets** as child issues, then **wire blocking edges** in a second pass (issues need ids first).
5. Stop - charting is one session. Don't resolve tickets yet.

### Work through the map
User invokes with a map URL or number.

1. **Load the map** - the low-res index, not every ticket body.
2. **Choose and claim** the ticket. If user named one, use it. Otherwise take the first frontier ticket. Assign it before any work.
3. **Resolve it** - fetch related ticket bodies on demand. Use `/grilling` and `/domain-modeling` if needed. Consult skills listed in the map's Notes.
4. **Record the resolution**: post the answer as a comment, close the issue, append to Decisions-so-far on the map.
5. **Graduate fog**: create newly-surfaced tickets, clear graduated fog from Not yet specified. If a ticket sits beyond the destination, rule it out of scope and close it.
6. **Never resolve more than one ticket per session.**

## Rules

- **Must**: Invoke only on explicit user request. Never auto-trigger.
- **Must**: Claim a ticket (assign it) before any work on it.
- **Must**: Never resolve more than one ticket per session.
- **Must**: Produce decisions, not deliverables - unless the map's Notes override this.
- **Must**: Refer to tickets by name (title), never by bare id or number.
- **Must**: Record resolution as a comment on the closed ticket; gist it in the map's Decisions-so-far.
- **Must**: Ticket a decision only when the question can be stated precisely. Everything else stays in fog.
- **Must not**: Chart past the frontier. Don't pre-ticket what can't yet be seen.
- **Must not**: Answer own questions in a HITL ticket. The human speaks for themselves.
- **Must not**: Charge at the destination - each ticket resolves one decision, nothing more.

## Out of scope

This skill does NOT handle:

- Small, single-session tasks - use the normal mission lane instead
- Implementation or code delivery - this is planning, not building
- Debugging - use sc-debug
- Ambiguity resolution within a mission - use sc-clarify (pathfinder is for multi-session scoping efforts)
- Git operations - use sc-git

## Output format

A `wayfinder:map` issue and its child tickets on the tracker. Each closed ticket has a resolution comment. The map's Decisions-so-far grows one line per resolved ticket:

```markdown
- [<closed ticket title>](link) - <one-line gist of the answer>
```

## Checklist

Before claiming a pathfinding session is complete:

- [ ] Destination named and fixed in the map
- [ ] Ticket claimed (assigned) before any work
- [ ] Exactly one ticket resolved this session
- [ ] Resolution posted as a comment on the closed ticket
- [ ] Decisions-so-far updated with gist + link
- [ ] Newly-surfaced tickets created and blocked correctly
- [ ] Graduated fog cleared from Not yet specified
- [ ] Out-of-scope tickets closed and noted in Out of scope
- [ ] Tickets referred to by name, never bare id

---

## References

- `references/tracker.md` - tracker-specific operations (creating child issues, labeling, blocking edges, querying frontier)
- `references/examples.md` - worked examples of maps and ticket resolutions
