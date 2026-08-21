---
id: design-system/banner
---

# Banner

A banner is a prominent notification item that displays important messages to users. It communicates things like errors, success confirmations, warnings, or general information using distinct colors and placement to stand out from regular page content and grab attention.

- [ ] Style - How the banner types look (fill, text colour, border radius etc)
  - 💡 Ensure the banner stands out enough on your default UI as it's usage typically means it should not be missed
- [ ] Content - A title is a must but you can consider a description if you feel a title does not convey enough information
- [ ] Types (state) - Banners typically appear as 5 options: information, success, warning, error, neutral. This helps visually establish context and level of importance.
- [ ] Call to action button (state) - A banner may call for the user to do something, like fix the error it is calling out, or viewing the successful outcome of an event.
  - 💡 Only have a button if it feels needed, banners typically work best to inform
- [ ] Placement - Where the banner sits on the page, dependent on context and importance
  - 💡 Consider the impact the banner contents has on the user experience to determine where it's placed best e.g. information affecting the entire application should have a banner at the top on every page. Information only affecting a small feature may be within a relevant section.
- [ ] Dismissable (state) - Whether the banner can be dismissed. Positive/neutral banners (info and success) would suit this, but not warning/error banners as they are too critical.

Related

[Button Design system](../design-system/button.md)
[Icon Design system](../design-system/icon.md)
