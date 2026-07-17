---
name: sc-localize
description: "Review bilingual copy for cultural fit. Catch literal translations that break in the target culture. Activate on bilingual content, translation quality, multilingual UI, or locale-specific copy."
---

# sc-localize

Review bilingual/multilingual copy for cultural appropriateness. Catch direct translations that are grammatically correct but culturally broken, and suggest natural alternatives.

## When to use

Activate when:

- Bilingual or multilingual content is shown, pasted, or written - especially UI copy, labels, navigation, or marketing text
- User asks "does this translation sound natural", "is this correct in [language]", or similar quality checks
- UI contains side-by-side language labels, i18n resource keys, or locale files
- Reviewing a diff that adds or changes translated strings
- User reports that a translation "feels weird" or "sounds like machine translation"

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Identify languages** - Confirm the source language and target language(s). Note the domain context (travel, e-commerce, healthcare, gaming, etc.).

2. **Extract intent** - What does the source text mean to its native audience? What concept, emotion, or action is it trying to convey? Strip idioms and metaphors - what is the underlying idea?

3. **Check cultural fit** - For each translated string, verify against locale-specific rules (load `references/<locale>.md` if available). Check:
   - Does the collocation feel natural? (words that pair together differ by language)
   - Is the register appropriate? (formal/casual/friendly matches the product tone)
   - Does the metaphor or idiom carry across? (many don't - replace with local equivalent)
   - Is the sentence structure natural in the target language? (word order, subject-dropping, particles)
   - Does the translation convey the same intent, not just the same words?

4. **Flag issues** - For each problem found, report:
   - **Severity**: `broken` (meaningless/wrong), `awkward` (grammatical but unnatural), `drift` (meaning shifted)
   - **What's wrong**: specific diagnosis (wrong collocation, literal idiom, register mismatch, etc.)
   - **Why it fails**: explain the cultural/linguistic reason
   - **Suggestion**: 1-3 natural alternatives with explanation of why they work

5. **When no locale reference exists** - Apply universal rules. Note uncertainty. Consult a native speaker when the decision affects user-facing copy.

## Rules

### Core principle

- **Must**: Translate the concept, not the words. Start from "what does this mean to a native speaker?" - then find how the target culture expresses that same idea.
- **Must not**: Accept a translation just because it's grammatically correct. Natural ≠ grammatical.
- **Must not**: Preserve English sentence structure in the target language. Rebuild the sentence in the target language's natural rhythm.

### Collocation

- **Must**: Check that word pairs feel natural together in the target language. A word that's correct in isolation may be wrong in combination.
- **Must**: When flagging a collocation issue, show the problematic pair and explain why it's unnatural.

### Register and tone

- **Must**: Verify that the formality level matches the product context. A casual app should not sound like government documents. A banking app should not sound like chat slang.
- **Must**: Check that gendered or hierarchical language choices (where applicable) match the product's audience.

### Metaphors and idioms

- **Must**: Flag idioms and metaphors that don't carry across languages. Suggest a local equivalent or a plain rephrasing.
- **Must not**: Invent new compound words that don't exist in the target language just to match the English structure.

### UI-specific

- **Must**: Check that translated UI labels fit within expected character limits. Some languages are 30-50% longer than English.
- **Must**: Verify that the translated label is scannable and actionable - users should understand it at a glance.
- **Must**: For navigation, buttons, and CTAs, prefer verbs/actions that feel natural in the target language's UI conventions.

### Evidence

- **Must**: For each flagged issue, cite the specific linguistic rule or locale reference that supports the finding.
- **Prefer**: When available, reference real-world usage examples (app UI, signage, publication) rather than dictionary definitions.

## Out of scope

This skill does NOT handle:

- Machine translation setup or API configuration - use web service skills
- i18n framework wiring (i18next, react-intl, etc.) - use web frontend/backend skills
- Full-site translation - this reviews quality, it does not produce bulk translations
- Locale-specific formatting (dates, numbers, currency) - use web frontend/backend skills
- Accessibility or RTL layout - use sc-ux-design

## Output format

```
## Localization review: [source lang] → [target lang]

Domain: [travel / e-commerce / healthcare / ...]
Product tone: [casual / professional / playful / ...]

### Issues found: [N]

| # | Severity | String | Issue | Why | Suggestion |
|---|----------|--------|-------|-----|------------|
| 1 | broken   | "ล่าช่วงเวลา" | literal compound | "ล่า" + "ช่วงเวลา" is unnatural in Thai; sounds like hunting time itself | "นักเดินทางยืดหยุ่น" (Flexible Traveler) - conveys the concept through a natural Thai noun phrase |

> Example shown for th-TH locale. Replace with locale-appropriate examples for other language pairs.

### Clean strings: [N]
[strings that passed review, with brief confirmation]

### Locale notes
[any locale-specific context the team should know for future translations]
```

## Checklist

Before claiming localization review is done:

- [ ] Source and target languages confirmed
- [ ] Domain context noted
- [ ] Intent of each source string understood before checking translation
- [ ] Locale references loaded (if available)
- [ ] Collocations checked - word pairs feel natural
- [ ] Register matches product tone
- [ ] Idioms/metaphors flagged and alternatives suggested
- [ ] UI labels checked for scanability and fit
- [ ] Every flagged issue includes: severity, diagnosis, cultural reason, and suggestion
- [ ] No accepted translations based on grammar alone - naturalness verified

---

## References

Load details on demand from:

- `references/th-TH.md` - Thai localization pitfalls: collocations, register, UI conventions, common translation traps
- `references/universal.md` - universal rules applicable to any locale pair
