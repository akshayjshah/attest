attest
======

[![Build](https://github.com/akshayjshah/attest/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/akshayjshah/attest/actions/workflows/ci.yaml)
[![Report Card](https://goreportcard.com/badge/go.akshayshah.org/attest)](https://goreportcard.com/report/go.akshayshah.org/attest)
[![GoDoc](https://pkg.go.dev/badge/go.akshayshah.org/attest.svg)](https://pkg.go.dev/go.akshayshah.org/attest)

`attest` is a small package of type-safe assertion helpers. Under the hood,
it uses [cmp] for equality testing and diffing. You may enjoy `attest` if you
prefer:

- Type safety: it's impossible to compare values with different types.
- Brevity: assertions usually print diffs rather than full values.
- Minimalism: just a few assertions, not a whole DSL.
- Natural ordering: every assertion uses `got == want` order.
- Interoperability: assertions work with any `cmp.Option`.

## Installation

```
go get go.akshayshah.org/attest
```

## Usage

```go
package main

import (
  "testing"
  "time"

  "github.com/akshayjshah/attest"
)

func TestExample(t *testing.T) {
  attest.Equal(t, 1, 1)
  attest.NotEqual(t, 2, 1)
  attest.Approximately(
    t,
    time.Minute - 1, // got
    time.Minute,     // want
    time.Second,     // tolerance
  )
  attest.Zero(t, "")
  attest.Contains(t, []int{0, 1, 2}, 2)

  var err error
  attest.Ok(t, err)
  err = fmt.Errorf("read config: %w", io.EOF)
  attest.Error(t, err)
  attest.ErrorIs(t, err, io.EOF)

  // You can enrich the default failure message.
  attest.Equal(t, 1, 2, attest.Sprintf("integer %s", "addition"))

  // The next two assertions won't compile.
  attest.Equal(t, int64(1), int(1))
  attest.Approximately(t, 9, 10, 0.5)
}
```

Failed assertions usually print a diff. Here's an example using `attest.Equal`:

```
--- FAIL: TestEqual (0.00s)
    attest_test.go:58: got != want
        diff (+got, -want):
          attest.Point{
                X: 1,
        -       Y: 4.2,
        +       Y: 3.5,
          }
```

## Status: Stable

This module is stable. It supports the [two most recent major
releases][go-support-policy] of Go.

Within those parameters, `attest` follows semantic versioning. No
breaking changes will be made without incrementing the major version.

## Legal

Offered under the [MIT license][license].

[cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
[go-support-policy]: https://golang.org/doc/devel/release#policy
[license]: https://github.com/akshayjshah/attest/blob/main/LICENSE
