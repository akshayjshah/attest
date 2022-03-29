attest
======

[![Build](https://github.com/akshayjshah/attest/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/akshayjshah/attest/actions/workflows/ci.yaml)
[![Report Card](https://goreportcard.com/badge/github.com/akshayjshah/attest)](https://goreportcard.com/report/github.com/akshayjshah/attest)
[![GoDoc](https://pkg.go.dev/badge/github.com/akshayjshah/attest.svg)](https://pkg.go.dev/github.com/akshayjshah/attest)

`attest` is a small package of type-safe assertion helpers. Under the hood,
it uses [cmp] for equality testing and diffing. You may enjoy `attest` if you
prefer:

- Brevity: assertions usually print diffs rather than full values.
- Natural ordering: every assertion uses `got == want` order.
- Interoperability: assertions work with any `cmp.Option`.
- Minimalism: just a few assertions, not a whole DSL.

## An example

```go
package main

import (
  "testing"
  "time"

  "github.com/akshayjshah/attest"
)

func TestExample(t *testing.T) {
  attest.Equal(t, 1, 1)
  attest.Equal(t, int64(1), int(1)) // compiler error!

  attest.Approximately(
    t,
    time.Minute - 1, // got
    time.Minute,     // want
    time.Second,     // tolerance
  )
  attest.Approximately(t, 0.99, 1.0, 0.2)
  attest.Approximately(t, 9, 10, 0.5) // compiler error!

  err := fmt.Errorf("read from HTTP body: %w", io.EOF)
  attest.Error(t, err)
  attest.ErrorIs(t, err, io.EOF)

  attest.Zero(t, "")
  attest.Contains(t, []int{0, 1, 2}, 2)
}
```

Here's some example output from a failed assertion using `attest.Equal`:

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

## Installation, status, and support

Install `attest` with `go get github.com/akshayjshah/attest@v0.1.1`. As you can
see, the package is currently _unstable_. I hope to cut a stable 1.0 soon after
the Go 1.19 release.

`attest` supports the [two most recent major releases][go-versions] of Go, with
a minimum of Go 1.18.

## Legal

Offered under the [Apache 2 license][license].

[cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
[go-versions]: https://golang.org/doc/devel/release#policy
[license]: https://github.com/akshayjshah/attest/blob/main/LICENSE
