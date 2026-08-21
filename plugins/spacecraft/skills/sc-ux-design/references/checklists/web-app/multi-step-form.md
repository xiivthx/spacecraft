---
id: web-app/multi-step-form
---

# Multi-step form

A form split across multiple steps or screens to reduce cognitive load when collecting a large amount of information from the user.

- [ ] Progress indicator (state) - A clear visual showing how many steps exist and which one the user is currently on
  - 💡 Labels for each step help give context hint to future steps
- [ ] Step heading and context - A clear title and brief context for each step, so users know what they are being asked to provide.
- [ ] Field grouping - Each step containing only fields that belong together logically e.g. phone and email for contact details
- [ ] Step-level validation - Validation happening at each step before the user proceeds, not surfaced all at once at the end
- [ ] Back navigation (state) - The ability to return to a previous step to review or change answers without losing subsequent progress
- [ ] Save and resume (state) - The ability to save progress and return to the form later
  - 💡 Particularly useful for forms that require information that connect to be instantly provided or sourced
- [ ] Final review step - A summary of all entered information before final submission, giving the user a chance to review and edit
  - 💡 Offer a link to each section at this step for editing

Related

[Submitting a form Web app](../flow/submitting-a-form.md)
[Input Field Web app](../design-system/input-field.md)
[Button Design system](../design-system/button.md)
[Dropdown Menu Design system](../design-system/dropdown-menu.md)
