---
id: mobile-app/login
---

# Login

Everything a returning user needs to authenticate quickly and securely.

- [ ] Social sign-in - Sign-in options that connect to an existing Apple or Google account, bypassing manual credential entry.
  - 💡 On iOS, Apple Sign In is required by App Store guidelines if any other social provider is offered.
- [ ] Email field - The input where users enter the email address associated with their account.
- [ ] Password field - A masked text input for the account password, with the option to reveal what has been typed.
- [ ] Biometric authentication - Face ID or fingerprint sign-in for returning users who have already authenticated once with a password
  - 💡 Encouraged to suggest after initial login so it's a faster experience in the future
- [ ] Credential autofill (state) - System-level support for pre-filling saved email and password from the user's password manager.
  - 💡 textContentType on iOS and autoComplete on Android are the attributes that trigger native autofill.
- [ ] Forgot password link - The link users reach for when they can't recall their password, leading into the reset flow.
- [ ] Error states (state) - Feedback shown when authentication fails, distinguishing between an unrecognised email address and an incorrect password.
  - 💡 Generic 'incorrect credentials' gives users no useful signal  -  knowing whether the email or password is wrong helps them recover without guessing.
- [ ] Passwordless sign-in (magic link) - An alternative sign-in method that sends a one-time link to the user's email, requiring no password
  - 💡 Useful for infrequent-use apps where remembering a password between sessions is difficult

On other platforms

[Login Website](../website/login.md)
[Login Web app](../web-app/login.md)

Related

[Resetting password Flows](../flow/resetting-password.md)
[Input Field Design system](../design-system/input-field.md)
[Verifying account Flows](../flow/verifying-account.md)
[Button Design system](../design-system/button.md)
[Toast Design system](../design-system/toast.md)
[2FA Web app](../web-app/2fa.md)
[Onboarding Web app](../web-app/onboarding.md)
