---
id: design-system/accordion
---

# Accordion

An accordion is a vertically stacked list of items that reveal or hide associated content sections when clicked. They help organize information hierarchically and saves screen space by showing only relevant content.

- [ ] Header - The clickable area that triggers expanding and collapsing and has a title
  - 💡 The header text should clearly convey the content that is inside
- [ ] Expand/collapse icon - Visual indicator of the current state - typically a plus/minus, or a caret/chevron that rotates between states
- [ ] Content area - The content that shows or hides when toggled, containing detailed information associated with the header
  - 💡 While headers should be text only, content can be anything that is suited to help answer or elaborate on the header (text, image, video, etc)
- [ ] States - Default (collapsed), expanded, hover, focused, and disabled
- [ ] Expanding logic - Decide if users can open multiple sections at once, or one at a time
  - 💡 Can be determined by whether you think information from other items are related and should be seen together

Related

[Tabs Design system](../design-system/tabs.md)
[Button Design system](../design-system/button.md)
