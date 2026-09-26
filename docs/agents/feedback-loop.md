# Feedback loop

This is a Go workspace (`go.work`): `trm` and each driver under `drivers/` is its own module.

## While working: the modules the change touches

Run inside each module the diff touches (`cd trm`, `cd drivers/sqlx`, ...):

```bash
go test -race -count=1 ./...
golangci-lint run -c ../../.golangci.yml --new-from-rev origin/main   # from trm: -c ../.golangci.yml
```

`--new-from-rev` reports only lines the change added or edited; a pre-existing finding on an untouched
line is not yours to fix.

## Before finishing: the whole project, once

From the repo root, the same commands CI runs:

```bash
make test
make lint
```

Both must be green. `make lint` runs golangci-lint in every module with the root `.golangci.yml`;
a finding it reports on a line the change did not touch is still worth naming, but it is not
blocking.

## Known to fail regardless of the change

- `TestTransaction_awaitDone_byContext` and `TestTransaction_awaitDone_byRollback` in `drivers/gorm`
  race with the driver's background goroutine and fail now and then. A failure there with no gorm
  change in the diff is a flake, not a finding. Everything else red counts against the change.
