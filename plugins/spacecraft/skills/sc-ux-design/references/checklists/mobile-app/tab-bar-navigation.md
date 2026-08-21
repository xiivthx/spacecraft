---
id: mobile-app/tab-bar-navigation
---

# Tab Bar Navigation

The persistent bottom navigation bar that gives users access to the top-level sections of the app

- [ ] Tab count - Limited to the most important destinations - 3 to 5 items is the typical range
  - 💡 More than 5 sections dilutes each option and may not fit
- [ ] Icon and label - Each tab paired with both an icon and a text label
  - 💡 Labels may not be necessary but only if icons are explicitly clear at what lies within that section
- [ ] Active and default states - Active should be visually distinct whether it's colour, border or icon weight
- [ ] Badge counts (state) - Useful for tabs with unread counts like messages and notifications, updating in real time
  - 💡 Instead of a number, a dot is suitable if numbers feels less relevant
- [ ] Fixed presence - For the sections the tab bar applies to, tab bar should remain visible, and can be hidden on any page a level deeper
- [ ] Tap target size - Each tab at least 44×44pt to be reliably tappable with a thumb in any grip position
- [ ] Haptic feedback - A subtle tap haptic on tab selection confirming the action

Related

[gesture navigation Mobile app](../mobile-app/gesture-navigation.md)
