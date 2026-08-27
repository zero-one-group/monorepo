# Webhook fixtures

Real payloads captured from each provider's **test / sandbox** environment, replayed by
`pkg/testutils/webhook` against your receivers. Never hand-write a fixture — a payload you
invented proves nothing about the provider.

Layout: `<provider>/<name>.json` (raw body, byte-exact) and optional
`<provider>/<name>.headers.json` (the signature and event-id headers the provider sent).

How to capture: trigger the event in the provider's test mode and copy the body from its
webhook log or a request inspector (e.g. ngrok). One fixture per event/status you handle.

Every receiver test should include at minimum: `AssertIdempotent` on the happy event,
`AssertLateEventIgnored` with the provider's terminal event followed by an earlier one,
and `AssertRejectsBadSignature`.

Redact personal data (names, phone numbers, emails) in captured payloads before committing.
