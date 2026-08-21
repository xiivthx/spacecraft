---
id: flow/resetting-password
---

# Resetting password

Every now and then, a user can't remember their password to log back in. Luckily, it's a straightforward process to change it. It's important to make the experience feel straightforward, so the user feels like they've gotten back into their account as smooth as possible.

- [ ] Place reset link close to password field - Style it as a link to show it is clickable
- [ ] Ask for account details to verify (state) - In this case it's usually the email address that's requested, because it can recognise your account and be the channel the link is securely sent to. Note: If the user already entered their email address on the previous login page, that can be prefill this field and speed up the flow!
- [ ] Show information has been sent (state) - Based on the account information provided in Step 2, explain how they can continue. If an email was provided, send an email for the next step. If a mobile number was provided, send a code or link to open.
- [ ] The message sent explains next steps (state) - This could be a link to a page that allows the user to reset their password. It could also be a code for the user to provide on a page to verify their account, to then reset their password.
- [ ] Reset the password! - Whether it's a verification code or a link to click behind an email address, the next page should be a clear text field to enter a new password. You can provide guidelines if you have requirements for the password to pass a threshold of strength to be accepted.
- [ ] Password successfully reset - After the password has been reset, indicate the successful and push their momentum to their initial intent: logging in.

Related

[Login Web app](../web-app/login.md)
[Input Field Design system](../design-system/input-field.md)
[Button Design system](../design-system/button.md)
[Login Website](../website/login.md)
[Login Mobile app](../mobile-app/login.md)
[Toast Design system](../design-system/toast.md)
[2FA Web app](../web-app/2fa.md)
[Showing input error Flows](showing-input-error.md)
