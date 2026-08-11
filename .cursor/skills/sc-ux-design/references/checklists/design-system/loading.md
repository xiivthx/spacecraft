---
id: design-system/loading
---

# Loading

A loading indicator is a visual element that communicates to users that content or an action is being processed. It provides feedback through animations like spinning wheels, progress bars, or skeleton screens to maintain user engagement during wait times.

- [ ] Visual indicator (state) - A clear representation that content is loading or in progress of change
- [ ] Text (state) - Explaining the loading state
  - 💡 Let the standard 'loading' text be a fallback, and try to be more specific e.g. 'adding your contacts, syncing your emails'
- [ ] Time (state) - Determine how long the time between two actions must be to require a loading component
  - 💡 Don't interrupt two actions that instantly transition, but if there's a known 5 second period between two actions then the component is appropriate
- [ ] Accessibility (state) - Ensure your loading state can be clearly seen
  - 💡 Place the loading indicator in a container to ensure accessibility, regardless of background If the loading indicator is on its own, create a version appropriate for light/dark backgrounds
- [ ] Visuals (state) - Entertain the user with an illustration during the loading state

Related

[Skeleton Design system](../design-system/skeleton.md)
[Icon Design system](../design-system/icon.md)
[Color System Design system](../design-system/color-system.md)
