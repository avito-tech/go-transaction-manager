# Working in this repo

- Table-driven test cases are named in `snake_case`: `"bare_err_skip_dropped_wrapped_one_kept"`.
- A user-visible fix or change gets a line in `CHANGELOG.md` under the next unreleased version
  (`## [x.y.z]`, Keep a Changelog); a released section stays as it is.
- Before pushing: `make test` and `make lint` (golangci-lint, rules in `.golangci.yml`). CI runs the same
  for every Go version in `.github/workflows`.
