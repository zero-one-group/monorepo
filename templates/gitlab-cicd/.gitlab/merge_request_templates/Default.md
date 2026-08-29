## What & why

<!-- One paragraph. Link the issue or decision record this depends on, if any. -->

## Definition of done

Tick what applies; explain anything unticked. Reviewers: do not approve a change to persistence,
queues or integrations with unticked boxes.

**Everyone**

- [ ] `moon run :lint :typecheck :build` passes locally
- [ ] Tests added or updated at the unit layer for the logic in this MR
- [ ] The project's coverage task passes; if a floor was raised, the history line was added
- [ ] Migrations (if any) are additive-first and apply cleanly from scratch (`migrate:reset --up`)
- [ ] `.env.example` regenerated if config changed

**If this touches persistence, a queue or an integration**

- [ ] Integration test against the real store (testcontainers), not a mock
- [ ] Webhook receivers pass the replay harness: idempotent · late-event-ignored ·
      bad-signature-rejected, using **recorded** fixtures under `testdata/webhooks`
- [ ] Correlation ID is on every new log line, span and audit row
- [ ] Amounts / statuses from a third party are looked up against our own records, never trusted
      from the payload

**If this touches a third-party API contract**

- [ ] Contract fixtures re-recorded from the provider's test account (`CONTRACT_RECORD=1`) and the
      diff read

## How to verify

<!-- Commands, URLs, fixtures a reviewer can use. -->
