# Prose rhythm

On-demand procedure: rewrite narrative/user-facing prose for natural rhythm. Decision job - not expertise cosplay. **Not** always-on. **Not** a writing tutor.

## Goal

Rewrite narrative/user-facing prose so it has natural rhythm (mix of short, medium, long sentences) without purple prose or filler.

## Output

Chat or file edit with:

```
### essence ###
<one-liner>
### essence ###

### rewrite ###
<full rewrite>
### rewrite ###
```

Or when editing in-place: essence as a one-line comment in the Task handshake / chat; rewrite applied to the file.

## Good / Bad

- Good: full rewrite (not light paraphrase); short/medium/long sentence mix; engaging second pass while staying precise; US English; ASCII hyphen-minus only
- Bad: expertise cosplay; threat or career-stakes framing; forced chain-of-thought; lyrical rhythm on Verify bars, JSON schemas, CLI flags, Never/Always lists, gate checklists, or code; monotone same-length sentence stacks; purple prose or filler

## Verify

When run: chat or file contains `### essence ###` / `### rewrite ###` blocks OR in-place rewrite with essence one-liner in handshake. Does **not** block discuss clear or ship.

## When to use

- Rewriting user-facing prose, README hero, handoff narrative, announcement copy, mission story
- Human asks for engaging rewrite or rhythm polish on narrative text
- Optional second pass after `narrative-context.md` draft

## When skip

- Machine-checkable Verify bars, JSON schemas, CLI flags, rule Never/Always lists, gate checklists, code
- Already compressed exact prose that must stay literal
- Routine structural edits with no narrative intent

## Procedure

1. **Act** - explain nothing long. Extract essence (one line).
2. **Rewrite** - full rewrite; must not look like a light paraphrase of the source.
3. **Rhythm** - mix short / medium / long sentences. Create music, not monotone stacks of same-length lines.
4. **Second pass** - make the rewrite more engaging while staying precise and US English; ASCII hyphen-minus only.
5. **Scope** - do **not** apply to machine-checkable Verify bars, JSON schemas, CLI flags, rule Never/Always lists, gate checklists, or code. Those stay compressed and exact.

## Must / Must not

- **Must**: Extract essence before rewrite
- **Must**: Full rewrite with sentence-length variety
- **Must**: Second pass for engagement without losing precision
- **Must not**: Expertise cosplay ("as a writing coach…")
- **Must not**: Threat, tips, or career-stakes framing
- **Must not**: Forced chain-of-thought on reasoning models
- **Must not**: Apply rhythm craft to Verify, gates, JSON, CLI flags, or code
- **Must not**: Unicode em dash - ASCII hyphen-minus only

## Related

- Narrative context before draft: `narrative-context.md`
- Agent: `.cursor/agents/sc-writer.md`
- Inspiration (sentence-length rhythm): Gary Provost - vary sentence length for music, not monotone stacks; craft only, not required reading; strip cosplay/threat framing from source prompts
