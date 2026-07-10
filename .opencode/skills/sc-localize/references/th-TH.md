> Consult when: reviewing Thai translations, bilingual TH-EN UI copy, or Thai locale-specific localization quality.

# Thai localization reference (th-TH)

## Core principle

Thai expresses concepts through **noun phrases, verb serialization, and classifiers** — not through English-style compound nouns. When translating, rebuild the phrase from the concept, not the words.

## Collocation traps

These English→Thai word pairs look correct in a dictionary but are unnatural together:

### Compound nouns

| English pattern | Trap | Why it fails | Natural Thai |
|-----------------|------|--------------|--------------|
| `[Verb]er + [Noun]` (Flexible Hunter) | `ล่า + ช่วงเวลา` (hunt + period) | "ล่า" takes concrete prey (สัตว์, คนร้าย), not abstract concepts. "ล่าช่วงเวลา" sounds like you're chasing time itself. | Restructure as noun phrase: `นักเดินทางยืดหยุ่น` (flexible traveler), or verb phrase: `ตามหาทริปที่ใช่` (find the right trip) |
| `[Adjective] + [Noun]` | `ยืดหยุ่น + นักล่า` | Thai adjective-noun order is noun-adjective, but `นักล่ายืดหยุ่น` still sounds like an unnatural calque | Use classifier: `คนที่เดินทางแบบยืดหยุ่น` or re-concept: `ปรับแผนได้` |
| `[Noun] + [Noun]` (Travel Planner) | `วางแผน + เดินทาง` (plan + travel) as compound | Thai prefers verb phrases over compound nouns | `วางแผนเที่ยว` (plan trips) — verb-first, natural |
| `Quick [Noun]` | `เร็ว + [คำนาม]` | "เร็ว" as prefix modifier is English calque | Use `[คำนาม]ด่วน`, `[คำนาม]ทันใจ`, or restructure: `จองเร็ว` → `จองปุ๊บได้ปั๊บ` |

### Verb collisions

| English | Dictionary trap | Natural Thai | Why |
|---------|----------------|--------------|-----|
| "Explore deals" | `สำรวจดีล` | `ค้นหาดีล` or `ส่องดีล` | "สำรวจ" is for physical exploration (land, cave). For browsing/searching, use ค้นหา, ส่อง, ดู |
| "Book now" | `จองตอนนี้` | `จองเลย` | "ตอนนี้" is clinical/instructional. "เลย" adds immediacy naturally. |
| "Save trip" | `บันทึกทริป` | `เซฟทริปไว้` or `เก็บทริป` | "บันทึก" is formal (save document). For casual UI: เซฟ or เก็บ |

### Adjective mismatches

| English | Trap | Natural Thai |
|---------|------|--------------|
| "Cheap flights" | `เที่ยวบินถูก` | `ตั๋วถูก` or `ไฟลต์ราคาเบา ๆ` |
| "Smart search" | `ค้นหาฉลาด` | `ค้นหาอัจฉริยะ` or `ค้นหาที่ใช่` |
| "Best price" | `ราคาดีที่สุด` | `ราคาดีที่สุด` is OK, but `ถูกสุด` or `คุ้มสุด` often better contextually |

## Register (ระดับภาษา)

Thai has distinct register levels. Wrong register = wrong product feel.

| Level | Use for | Particles | Example tone |
|-------|---------|-----------|--------------|
| ทางการ (formal) | banking, gov, legal | ครับ/ค่ะ, polite endings | `ดำเนินการ`, `โปรด`, `กรุณา` |
| กึ่งทางการ (semi-formal) | e-commerce, SaaS, booking | ครับ/ค่ะ, sometimes dropped | `เลือก`, `ค้นหา`, `ดู` |
| กันเอง (casual) | chat, gaming, social | เลย, นะ, อะ, ป่ะ | `ดูนี่`, `กดเลย`, `จัดไป` |
| วัยรุ่น (slang) | youth apps, Gen Z | ปัง, เก๋, จึ้ง, ฟิน | `ปังมาก`, `ฟินสุด ๆ` |

**Rule**: Match the product's voice. A travel booking app is usually semi-formal to casual (not formal, not slang). Banking must be formal. Gaming can be casual/slang.

## UI conventions

### Navigation labels

| English | Natural Thai | Why |
|---------|-------------|-----|
| Home | `หน้าหลัก` | Not `บ้าน` (that's a building) |
| Search | `ค้นหา` | Not `ค้น` (too short/abrupt) |
| My Trips | `ทริปของฉัน` | Not `ทริปฉัน` (possessive without `ของ` feels clipped) |
| Settings | `ตั้งค่า` | Not `การตั้งค่า` (nominalized form is too formal for UI) |
| Log in | `เข้าสู่ระบบ` | Standard. `ล็อกอิน` is common but informal. |

### Button/CTA

| English | Natural Thai | Note |
|---------|-------------|------|
| Get Started | `เริ่มเลย` or `เริ่มต้นใช้งาน` | Avoid `เริ่ม` alone (too abrupt) |
| Learn More | `ดูรายละเอียด` | Not `เรียนรู้เพิ่มเติม` (sounds like homework) |
| Try Free | `ทดลองใช้ฟรี` | Natural compound |
| Continue | `ต่อไป` | OK but `ถัดไป` sometimes better for pagination |

### Character budget

Thai text is typically 15-30% shorter than English for equivalent meaning. However, verbose literal translations can blow up. Keep UI labels under 20 Thai characters when possible.

## Common translation anti-patterns

### 1. English genitive as Thai possessive

`[X]'s [Y]` → `[Y]ของ[X]` is correct, but often overused. Many English possessives are better as Thai descriptive phrases:
- "User's trips" → `ทริปของคุณ` (not `ทริปของผู้ใช้` — too formal/literal)
- "Today's deals" → `ดีลวันนี้` (not `ดีลของวันนี้`)

### 2. Translating English prepositions literally

- "Search for flights" → `ค้นหาเที่ยวบิน` (not `ค้นหาสำหรับเที่ยวบิน`)
- "Deals from Bangkok" → `ดีลจากกรุงเทพ` is OK, but `ดีลออกจากกรุงเทพ` is natural for flights

### 3. Preserving English word order

English: "Find the cheapest flights to anywhere"
Trap: `หาเที่ยวบินถูกที่สุดไปทุกที่` (word-by-word)
Natural: `หาตั๋วถูกไปไหนก็ได้` (Thai structure: find-ticket-cheap-go-anywhere)

### 4. Using English-style plural/singular

Thai doesn't mark plural on nouns. Don't force it with classifiers unless needed for clarity:
- "Trips" → `ทริป` (not `หลายทริป` or `ทริปต่างๆ` unless disambiguation is needed)

### 5. Direct pronoun translation

- "You" → `คุณ` is fine in most UI. `ท่าน` is too formal (royal/government). Dropping the pronoun is often more natural.
- "We" / "Our" → `เรา` is standard for brand voice in Thai apps.

## When to use transliteration vs translation

| Use transliteration (ทับศัพท์) when | Use translation when |
|------------------------------------|---------------------|
| The English word is already standard in Thai tech lexicon: `ดีล`, `ทริป`, `แอป`, `ฟีเจอร์`, `ล็อกอิน` | A natural Thai equivalent exists and is common: `ตั้งค่า` (not `เซ็ตติ้ง`), `ค้นหา` (not `เสิร์ช`) |
| The concept is new/borrowed: `บล็อกเชน`, `คลาวด์` | The translated word is shorter or more scannable |
| The English brand name is used as-is in Thai market | The audience is broad and may not know the English term |

Hybrid is common and natural: `แอปทริป`, `ฟีเจอร์ใหม่`, `ดีลสุดคุ้ม`

## Spacecraft integration

- Flag issues in `evidence.jsonl` with label `locale-review:<lang>`
- Record locale-specific decisions in mission `decisions.md`
- Single source of truth per locale; no duplicated rules across missions
