> Consult when: reviewing bilingual content for locale pairs that lack a dedicated reference, or when establishing baseline quality gates for any translation.

# Universal localization rules

Rules that apply to any language pair regardless of specific locale knowledge.

## The concept-first rule

Before checking any translation, answer: **"What does the source text mean to a native speaker?"**

1. Strip all idioms, metaphors, and cultural references
2. Identify the core concept, emotion, or action
3. Find how the target culture naturally expresses that same concept
4. Rebuild the phrase in the target language's natural structure

**Anti-pattern**: translating word-by-word, then rearranging. Start from the concept, not the words.

## Structure independence

Every language has its own:
- Word order (SVO, SOV, VSO, free order)
- Preferred phrase types (noun phrases vs verb phrases vs clauses)
- Information hierarchy (what comes first matters differently)
- Dropping rules (subjects, objects, articles implied by context)

**Rule**: Do not preserve source language sentence structure. Rebuild in target language structure.

| Source pattern | Don't | Do |
|---------------|-------|-----|
| English compound noun ("Travel Planner") | Force a noun-noun compound | Use the target language's natural equivalent (verb phrase, clause, classifier construction) |
| English -er/-or agent ("Finder", "Hunter") | Invent an agent noun | Use descriptive phrase or reconceptualize |
| Adjective-Noun order | Preserve source order | Follow target language's adjective placement |

## Collocation check

Collocation = words that naturally pair together in a language. They differ between languages and can't be predicted from dictionaries.

**Check**: For any `[word A] + [word B]` pair in the translation, ask:
- Do native speakers put these words together?
- Would a different word B be more natural with word A?
- Is this a calque (literal borrowing) of the source language's collocation?

**Red flag words** - these often signal calque when paired:
- "get", "make", "take", "have" + noun (English light verbs)
- "smart", "quick", "fast", "easy" (English prefix-modifiers)
- "explore", "discover", "experience" (tourism/marketing verbs that don't always map cleanly)

## Register alignment

| Register mismatch | Example | Fix |
|------------------|---------|-----|
| Source is casual, translation is formal | "Hey, check this out!" → stiff formal equivalent | Match the casual tone in target language's casual register |
| Source is professional, translation is slangy | Banking app using chat slang | Match professional register |
| Source has brand voice, translation is generic | Playful SaaS becomes neutral | Preserve playfulness using target language's playful conventions |

**Rule**: If you can't match the register exactly, note it and recommend the closest natural register.

## Idiom and metaphor handling

| Situation | Action |
|-----------|--------|
| Idiom has direct equivalent in target language | Use it (rare) |
| Idiom has different equivalent with same meaning | Use the target language's idiom |
| Idiom has no equivalent | Translate the meaning plainly. Do not invent a new idiom. |
| Metaphor uses culture-specific reference (sports, food, landmarks) | Replace with a metaphor from the target culture that carries the same feeling - or drop the metaphor and state the meaning directly |

## Length budget

Many languages are longer than English for equivalent meaning:
- German: +20-35%
- French: +15-25%
- Spanish: +15-25%
- Japanese: -10% to +20% (kanji short, kana long)
- **Thai: -15% to +10%** (often shorter, but literal translations blow up)

For UI labels, always test the translation in a component with the real font and width.

## Machine translation smell detection

Signs a translation is machine-generated and needs human review:

1. **Overly literal**: Every English word has a 1:1 target word, in English order
2. **Dictionary-default choices**: Uses the first dictionary entry for each word, not the context-appropriate one
3. **Missing particles/fillers**: Languages with sentence-final particles (Thai, Japanese, Korean) have them stripped
4. **No elision**: Every pronoun, article, and function word is explicitly translated even when the target language drops them
5. **Register flat**: All text is the same neutral/formal register regardless of context

When any 3 of these signs are present, flag the entire text for human review - don't attempt spot-fixes.

## When to escalate

- The translation involves legal, medical, financial, or safety-critical text - never rely on AI-only review
- The target language has significant dialect or regional variation (e.g., Arabic, Chinese, Spanish) - clarify which region
- The product has existing translation memory or glossary - respect it; report conflicts, don't override

## Spacecraft integration

- Universal rules are the fallback when no locale reference exists
- Each locale reference overrides universal rules where they conflict
- Flag evidence with `locale-review:universal` when no locale-specific reference was used
