---
id: web-app/account
---

# Account

Where users view and manage their personal information, preferences, and account-level details

- [ ] Profile photo - A way for users to upload or change their profile image
  - 💡 Consider a sensible fallback like initials or a placeholder icon for users who haven't uploaded a photo so a profile image still exists
- [ ] Display name - The name shown to other users or across the product interface e.g. username, email address, first and last name
  - 💡 Worth clarifying whether this is a public-facing name or an internal account label
- [ ] Account details - Fields for email address, phone number, job title, or other relevant identifying information based on the product and information collected
  - 💡 If the information becomes extensive, it's suitable to group e.g. contact information, job history
- [ ] Linked accounts (if applicable) - A view of which third-party accounts are connected for sign-in or data access, with ability to disconnect
- [ ] Save confirmation (state) - Clear feedback that changes have been saved, either inline or as a toast
  - 💡 Auto-save with a subtle confirmation is more pleasant than explicit save, but if you want an explicit save button, it should remain disabled until there are changes
- [ ] Delete or deactivate account - Options to deactivate or permanently delete the account, clearly separated from other settings

On other platforms

[Account Mobile App](../mobile-app/account.md)

Related

[Verifying account Flows](../flow/verifying-account.md)
[Deleting account Flows](../flow/deleting-account.md)
[Account Mobile app](../mobile-app/account.md)
