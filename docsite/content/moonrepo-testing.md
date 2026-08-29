---
title: Testing
slug: "moonrepo-testing"
---

How projects generated from these templates are tested: the test-first loop, which test to
write first for each layer, how fixtures are handled, and how the coverage floors are raised
over time. The templates ship the harnesses; this page is the workflow that goes with them.

## The loop

Every feature, fix or refactor is done test-first. Reviewers can tell when the tests were
written after the code, and the coverage gate will tell them when they were not written at all.

1. **Write the failing test first** — the smallest test that expresses the behaviour asked
   for, named after the behaviour (`TestCheckout_RejectsUnavailableItem`,
   `it('shows the confirmation once the order is placed')`), not the function. Run it. It must
   fail for the right reason: an assertion, not a compile error you then "fix" by stubbing.
2. **Make it pass with the least code** — no speculative branches, no config for cases the
   test does not cover.
3. **Refactor with the test green** — extract, rename, dedupe. Run the whole package/app
   suite, not just the new test.
4. **Repeat per behaviour.** One assertion cluster per test; a new edge case is a new test,
   not an extra `if` in an existing one.
5. Before opening the PR: `moon run :lint :typecheck :build`, the project's coverage task,
   and the definition-of-done checklist below.

Never delete or weaken a test to get green. If a test is wrong, say so explicitly and fix the
test in its own commit.

## Which test to write first (Go, `go-modular`)

| You are writing…            | Start with                                                                                                                                                                   | Harness                                  |
| --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------- |
| Domain rules (state machines, pricing, validation) | Pure table-driven unit test in `modules/<x>/domain`, every input pair enumerated                                                                                      | `testing` + testify, no DB               |
| A service                   | Unit test with a fake or mocked repository (`moon run go-modular:generate-mock`) covering success, each error branch, and idempotency                                        | testify + mockery                        |
| A repository                | Integration test against real Postgres/Redis via `pkg/testutils.NewTestEnv` + `RunAppMigrations`                                                                              | testcontainers                           |
| A handler / route           | `httptest` through the real Echo router with a fake service: status code, error envelope, auth required                                                                       | Echo `httptest`                          |
| Module wiring               | An `fx_test.go` next to `fx.go`: config validation, the route table, which routes sit behind the JWT guard (see `modules/auth/fx_test.go`)                                    | Echo `httptest`                          |
| A webhook receiver          | Implement `pkg/testutils/webhook.Receiver`; the first test is `AssertIdempotent` + `AssertLateEventIgnored` + `AssertRejectsBadSignature` over a **recorded** fixture in `testdata/webhooks/<provider>/` | replay harness              |
| A third-party client        | Contract test behind `//go:build contract` using `pkg/testutils/contract`; record once with `CONTRACT_RECORD=1` against the provider's test account, replay offline after       | contract recorder                        |
| A migration                 | Apply on a fresh DB in a test (`RunAppMigrations`), assert the table/constraint exists, then write the repository test that uses it                                            | testcontainers                           |
| The app as a whole          | Extend `internal/app/app_integration_test.go` — real Postgres, the fx graph exactly as `New()` builds it, one request per new module                                          | testcontainers + fxtest                  |

## Which test to write first (React, `react-app`)

| You are writing…                                   | Start with                                                                                                                                                       |
| -------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| State/derivation logic (totals, status → copy)     | Pure unit test of the function or hook in `tests/`, before any JSX                                                                                              |
| A component with behaviour (form, dialog, toggle)  | Testing Library test: render, act as the user (`userEvent`), assert on what the user sees — never on implementation details or class names                        |
| A route                                            | `renderAt('/path')` from `tests/helpers.tsx` (the real router): loader/action outcomes, error and empty states, the redirect when unauthenticated                |
| An API call                                        | Mock at the network boundary (`msw` or `vi.fn` on the client); assert the request shape and both success and error rendering                                     |
| A user journey                                     | One Playwright spec per journey in `tests-e2e/`; extend it rather than adding a second file for the same journey                                                |
| A `shared-ui` primitive                            | A Storybook story per variant is the test; add interaction tests only for behaviour Radix does not already guarantee                                            |

Do not test template filler (`-welcome.tsx`), generated files (`routeTree.gen.ts`) or styling
tokens — they are excluded from the coverage denominator for that reason.

## Fixtures and test data

- Third-party payloads are **recorded, never invented**. A hand-written provider payload
  proves nothing about the provider. See `testdata/webhooks/README.md` and
  `testdata/contract/README.md` in `go-modular`.
- Seeders (`database/seeders`) are the only source of "realistic" local data; add a factory
  there rather than inline SQL in tests.
- Redact names, phone numbers and emails in any captured payload before committing.

## Coverage: report, then ratchet

Both templates measure coverage against an honest denominator (every source file, not only
the ones a test happened to import) and ship with **floors at zero** — a template cannot know
your project's baseline.

- `moon run go-modular:coverage` — `scripts/coverage-gate.sh` prints three tiers: overall,
  `modules/...`, and the packages named in `COVERAGE_CRITICAL_PACKAGES` (payments,
  webhooks — whatever must never regress). Floors: `COVERAGE_MIN`, `COVERAGE_MIN_MODULES`,
  `COVERAGE_MIN_CRITICAL` in `moon.yml`.
- `moon run react-app:test-coverage` — vitest with `coverage.include: ['app/**']`; floors
  go in `coverage.thresholds` in `vitest.config.ts`.

**The ratchet rule:** once you have a baseline, set each floor to *measured − 5* and only ever
raise it. Record the history next to the value (`YYYY-MM-DD what landed: measured → floor`) so
the trend is visible. Lowering a floor is a decision that belongs in the PR description, not
a quiet edit.

## Definition of done

The `gitlab-cicd` template ships a merge-request template with this checklist
(`.gitlab/merge_request_templates/Default.md`). For GitHub, put the same list in
`.github/pull_request_template.md`.

- `moon run :lint :typecheck :build` passes locally.
- Tests added or updated at the unit layer for the logic in the change.
- Anything touching persistence has an integration test against the real store (testcontainers).
- Webhook receivers pass the replay harness (idempotent · late-event-ignored ·
  bad-signature-rejected) on **recorded** fixtures.
- Migrations are additive-first and apply cleanly from scratch (`migrate:reset --up`).
- `.env.example` regenerated if config changed.
- The coverage task passes; if a floor was raised, the history line was added.

## CI

`templates/gitlab-cicd/test/all.yml` runs exactly the local gates so CI and `moon run …` agree:
`lint-typecheck-build` (only affected projects, unless a workspace-level file changed), a
`test-api` job with Docker-in-Docker for testcontainers and the gate's `overall` line as the
job coverage, and a `.test-frontend` job template per SPA that publishes the JUnit report.

Before pushing a CI change, reproduce the clean-checkout conditions locally:

```sh
MOON_CACHE=off moon run :lint :typecheck :build   # no cached task outputs
git clean -fdX                                    # drop every ignored file (generated code, node_modules, build/)
pnpm install && moon run go-modular:coverage react-app:test-coverage
```

Two moon behaviours bite here: tasks inferred from `package.json` have no declared outputs,
so moon can report them as cached with nothing on disk; and embedded/generated inputs
(`docs/swagger.json`, `routeTree.gen.ts`) must be produced by a task the test task depends on.
The templates declare both explicitly — keep it that way when adding tasks.

## Template smoke tests

Building and running each template end to end after changing it:

```sh
# go-modular / go-clean
moon go-modular:build && moon go-modular:start
moon go-modular:docker-build && moon go-modular:docker-run

# react-app / nextjs-app / strapi-cms
moon react-app:build && moon react-app:start
moon react-app:docker-build && moon react-app:docker-run

# shared-ui
moon shared-ui:build && moon shared-ui:storybook-build
```

The acceptance test for a template change is to generate a throwaway project from the modified
template and run `moon run :lint :typecheck :build` plus its coverage task there — remember that
moon caches downloaded templates (see [Troubleshooting]({{< ref "troubleshooting.md" >}})).
