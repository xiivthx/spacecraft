# Chat ask format

Load before the first frontier ask. Field order and bold labels are greppable - do not rename.

**Language:** Chat ask field bodies are **Thai content**; bold **English labels** stay greppable (`Plain explain`, `Context`, `Question`, `Trade-offs`, `Why it matters`, `Recommendation`, `If accepted`). Paths, commands, product names, and tech terms stay as-is. No dual `ไทย:`/`EN:` Chat blocks. `questions.md` / `decisions.md` stay English. Skill instructions on disk stay US English.

**Short vs rich:**
- **Rich** - question has explicit choices (A/B/…) **or** is a heavy blocking class (Verify / architecture fork / in-out scope). Requires Plain explain, Context, and per-choice Trade-offs.
- **Short** - simple yes/no or single-path questions. Skip Plain explain / Context / Trade-offs.

**Short template** (English labels; Thai bodies):
```
**Q1 - <short title>:** <คำถาม>

**Why it matters:** <หนึ่งประโยค>
**Recommendation:** <คำแนะนำ + เหตุผลสั้น>
**If accepted:** <ขั้นถัดไป>
```

**Rich template** (English labels; Thai bodies):
```
**Q1 - <short title>**

**Plain explain:** <ปัญหาคืออะไร ทำไมต้องตัดสินใจ เลือกแล้วเปลี่ยนอะไร - ภาษาง่าย สั้น>
**Context:** <รู้อะไรแล้วจาก spec/decisions/repo; อะไรยังคลุมเครือ; stake คืออะไร>
**Question:** <คำถาม>
- A) ...
- B) ...

**Trade-offs:**
- **A)** Pros: … | Cons: …
- **B)** Pros: … | Cons: …

**Why it matters:** <หนึ่งประโยค>
**Recommendation:** <คำแนะนำ + เหตุผลสั้น>
**If accepted:** <ขั้นถัดไป>
```

Keep trade-offs tight: 1-2 bullets each side per choice.

**Rich micro-example** (copy shape, not content):
```
**Q1 - Choose API framework**

**Plain explain:** ต้องเลือก HTTP library สำหรับเซิร์ฟเวอร์ การเลือกมีผลต่อความเร็วตอนสตาร์ท การซัพพอร์ต TypeScript และความคุ้นของทีมตอนดีบัก
**Context:** Spec ล็อก Node.js + TypeScript แล้ว แต่ยังไม่ล็อก framework ในรีโปยังไม่มี server package Stake คือ binding เส้นทาง/plugin ก่อนวางแผน
**Question:** ใช้ Fastify หรือ Express?
- A) Fastify
- B) Express

**Trade-offs:**
- **A)** Pros: เร็ว, TS-first | Cons: ทีมอาจคุ้น Express มากกว่า
- **B)** Pros: เอกสาร/ตัวอย่างเยอะ | Cons: ช้ากว่า พิมพ์ยากกว่า

**Why it matters:** เลือกผิดต้อง rewrite routing ตอน implement
**Recommendation:** A) Fastify - เข้ากับสแต็ก TS และเป้าสตาร์ท
**If accepted:** บันทึกใน decisions.md แล้วถาม frontier ข้อถัดไป (หรือ clear ถ้า frontier ว่าง)
```
