# FS Shim
[![License](https://img.shields.io/badge/license-mit-blue.svg?style=for-the-badge)](https://raw.githubusercontent.com/abursavich/fs-shim/main/LICENSE)
[![GoDev Reference](https://img.shields.io/static/v1?logo=go&logoColor=white&color=00ADD8&label=dev&message=reference&style=for-the-badge)](https://pkg.go.dev/bursavich.dev/fs-shim)
[![Go Report Card](https://goreportcard.com/badge/bursavich.dev/fs-shim?style=for-the-badge)](https://goreportcard.com/report/bursavich.dev/fs-shim)

FS Shim provides build-dependent implementations of the `io/fs` and `testing/fstest` packages, which were introduced in go 1.16.

- With go 1.16 or later, an implementation which aliases the standard library is used.
- With an earlier version of go, a forked version of the go 1.16 implementation is used.

This allows use features of `io/fs` without immediately requiring go 1.16.

As soon as it's reasonable to require go 1.16, use of this shim should be removed.
