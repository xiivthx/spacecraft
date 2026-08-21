---
id: web-app/maintenance
---

# Maintenance

A screen shown when the application is temporarily unavailable due to scheduled maintenance or an unexpected outage.

- [ ] Clear status message - A plain-language explanation that the product is currently unavailable and why
  - 💡 Scheduled maintenance and unexpected outages call for different messaging, since users respond very differently to each, and conflating them feels evasive
- [ ] Estimated return time - When the product is expected to be back online, as specifically as possible
  - 💡 Vague messages like 'back soon' frustrate users, since even a rough estimate like 'within 2 hours' is far more reassuring
- [ ] Status page link (state) - A link to a live status page where users can monitor progress and see real-time updates
  - 💡 A status page hosted on a separate domain remains accessible even when your main infrastructure is down
- [ ] Contact or support link - A way to reach support for urgent issues that cannot wait for the maintenance window to end
- [ ] Brand consistency - A maintenance page styled consistently with the product, even if the content is minimal
  - 💡 A well-designed maintenance page communicates professionalism and reduces user anxiety

Related

[Settings Web app](../web-app/settings.md)
[Saving changes Flows](../flow/saving-changes.md)
[404 Website](../website/404.md)
