attest
======

[![Build](https://github.com/akshayjshah/attest/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/akshayjshah/attest/actions/workflows/ci.yaml)
[![Report Card](https://goreportcard.com/badge/github.com/akshayjshah/attest)](https://goreportcard.com/report/github.com/akshayjshah/attest)
[![GoDoc](https://pkg.go.dev/badge/github.com/akshayjshah/attest.svg)](https://pkg.go.dev/github.com/akshayjshah/attest)

`attest` is a small package of type-safe assertion helpers. Under the hood,
it uses [cmp] for equality testing and diffing. You may enjoy `attest` if you
prefer:

- Type safety: it's impossible to compare values with different types.
- Brevity: assertions usually print diffs rather than full values.
- Minimalism: just a few assertions, not a whole DSL.
- Natural ordering: every assertion uses `got == want` order.
- Interoperability: assertions work with any `cmp.Option`.

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
  err = fmt.Errorf("read from HTTP body: %w", io.EOF)
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

## Status and support

`attest` supports the [two most recent major releases][go-versions] of Go, with
a minimum of Go 1.18. Its currently _unstable_, but I hope to cut a stable 1.0
soon after the Go 1.19 release.

## Legal

Offered under the [Apache 2 license][license].

[cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
[go-versions]: https://golang.org/doc/devel/release#policy
[license]: https://github.com/akshayjshah/attest/blob/main/LICENSE
