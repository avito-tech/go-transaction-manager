# Working in this repo

- Table-driven test cases are named in `snake_case`: `"bare_err_skip_dropped_wrapped_one_kept"`.
- A user-visible fix or change gets a line in `CHANGELOG.md` under the next unreleased version
  (`## [x.y.z]`, Keep a Changelog); a released section stays as it is.
- Before pushing: `make test` and `make lint` (golangci-lint, rules in `.golangci.yml`). CI runs the same
  for every Go version in `.github/workflows`.

## Changing `trm` and a driver together: release order

Every driver is its own Go module and requires `trm/v2` by version from the module proxy; the root
`go.work` applies only inside this checkout. So a driver cannot use a new `trm` symbol until that `trm`
is released. CI on Go 1.13 and 1.14 (no workspaces) fails with `undefined: drivers.X`, and that is a
true signal: anyone who `go get`s the driver fails the same way. A `replace` directive or a
`go.work.sum` entry hides it from CI without fixing it for consumers — do not add them.

Release in two steps, the way 2.0.2 went (`trm/v2.0.2` on 37eaee7; drivers bumped and tagged on 872ba58):

1. A PR with the `trm` change alone, plus its CHANGELOG line under the next version. Merge it.
2. Tag `trm` only, on that merge: `git tag trm/vX.Y.Z && git push origin trm/vX.Y.Z`. Wait until
   `https://proxy.golang.org/github.com/avito-tech/go-transaction-manager/trm/v2/@v/vX.Y.Z.info` answers.
   Not `make tag` yet: it tags every module, and the drivers would get the number before their change.
3. A PR with the driver change: `go get github.com/avito-tech/go-transaction-manager/trm/v2@vX.Y.Z` in the
   drivers that need it (bump the other drivers too, so they stay on one `trm`), then the code. Merge it.
4. `make tag version=Y.Z` on that merge: `trm/vX.Y.Z` already exists and is skipped, every driver gets
   `drivers/<name>/vX.Y.Z`, and `make tag.pkg` checks the proxy. One version, one CHANGELOG section, two commits.
