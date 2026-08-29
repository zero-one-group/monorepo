# Contract fixtures

Recorded responses from third-party **test** environments, replayed offline by
`pkg/testutils/contract`. Contract tests are behind the `contract` build tag.

```sh
moon run go-modular:test-contract                        # replay (offline, fast)
CONTRACT_RECORD=1 <PROVIDER>_API_TOKEN=… moon run go-modular:test-contract   # re-record
```

Layout: `<provider>/<name>.json`, one directory per provider. Only ever record against a
provider's test/sandbox account — never a live one.

Re-record deliberately (the provider announced a change, or a status code you did not know
about), commit the diff, and read it: the fixture diff *is* the contract change.
