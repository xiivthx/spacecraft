---
id: web-app/2fa
---

# 2FA

A screen that guides users through setting up or completing two-factor authentication to add a second layer of security to their account

- [ ] Method selection - Authenticator app, SMS, or email code
  - 💡 Offer at least 2 methods as some users potentially don't have access to one or the other
- [ ] Setup instructions - Clear step-by-step guidance for completing setup, particularly for authenticator flows that require scanning a QR code
- [ ] QR code or setup key - A scannable QR code or copyable key for linking an authenticator app to the account
  - 💡 Ensure QR code size large enough to be easy to scan
- [ ] Verification step - A code entry step confirming the setup was successful before 2FA is enabled on the account
  - 💡 Enabling 2FA without this step risks a silent setup failure that locks the user out  -  a serious support burden
- [ ] Recovery codes - A set of one-time backup codes the user can use if they lose access to their 2FA method
  - 💡 Make download or copying the code a required step before completing setup for the sake of the user
- [ ] Setup confirmation (state) - A clear success state confirming that 2FA is now active on the account
- [ ] Disable or reset option - A way for users to turn off or reconfigure 2FA, accessible from account security settings
  - 💡 Re-authentication should be required before disabling 2FA

Related

[Login Web app](../web-app/login.md)
[Verifying account Flows](../flow/verifying-account.md)
[Account Mobile app](../mobile-app/account.md)
