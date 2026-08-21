---
id: mobile-app/account
---

# Account

Private account settings like credentials, linked accounts, notifications, and destructive actions.

- [ ] Email - The current email address displayed with an option to update it
  - 💡 Verifying the new address before completing the change is standard practice, where notifying the old address too adds a useful security signal.
- [ ] Password change - A way for users to update their account password.
  - 💡 Requiring the current password before accepting a new one prevents unauthorised changes on an unlocked device
- [ ] Linked accounts - A view of which third-party accounts are connected for sign-in or data access, with ability to disconnect
- [ ] Save confirmation (state) - Clear feedback that changes have been saved, either inline or as a toast
  - 💡 Auto-save with a subtle confirmation is more pleasant than explicit save, but if you want an explicit save button, it should remain disabled until there are changes
- [ ] Delete or deactivate account - Options to deactivate or permanently delete the account, clearly separated from other settings

On other platforms

[Account Web app](../web-app/account.md)

Related

[Login Mobile app](../mobile-app/login.md)
[Settings Web app](../web-app/settings.md)
[Billing Website](../website/billing.md)
[Account Web app](../web-app/account.md)
[Deleting account Flows](../flow/deleting-account.md)
[Verifying account Flows](../flow/verifying-account.md)
