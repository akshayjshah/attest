// Package attest is a small, type-safe library of assertion helpers.
//
// Under the hood, attest uses github.com/google/go-cmp/cmp to test equality
// and diff values. All of attest's assertions work with any cmp.Option.
package attest

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Equal asserts that two values are equal.
func Equal[T any](tb TB, got, want T, opts ...Option) bool {
	tb.Helper()
	t := newAttester(tb, opts...)
	diff := t.Diff(got, want)
	if diff == "" {
		return true
	}
	t.Printf("got != want")
	t.Printf("diff (+got, -want):")
	t.Printf(diff)
	return t.Attest()
}

// Ok asserts that the error is nil.
func Ok(tb TB, err error, opts ...Option) bool {
	tb.Helper()
	if err == nil {
		return true
	}
	t := newAttester(tb, opts...)
	t.Printf("unexpected error")
	t.Printf("error: %v", err)
	t.Printf("type: %T", err)
	return t.Attest()
}

// Error asserts that the error is nil.
func Error(tb TB, err error, opts ...Option) bool {
	tb.Helper()
	if err != nil {
		return true
	}
	t := newAttester(tb, opts...)
	t.Printf("unexpected success")
	return t.Attest()
}

// ErrorIs asserts that got wraps want, using the same logic as the standard
// library's errors.Is.
func ErrorIs(tb TB, got, want error, opts ...Option) bool {
	tb.Helper()
	if errors.Is(got, want) {
		return true
	}
	t := newAttester(tb, Options(opts...), Cmp(cmpopts.EquateErrors()))
	t.Printf("got doesn't wrap want")
	t.Printf("diff (+got, -want):")
	t.Printf(t.Diff(got, want))
	return t.Attest()
}

// Zero asserts that the value is its type's zero value.
func Zero[T any](tb TB, got T, opts ...Option) bool {
	tb.Helper()
	var zero T
	t := newAttester(tb, opts...)
	diff := t.Diff(got, zero)
	if diff == "" {
		return true
	}
	t.Printf("got is non-zero")
	t.Printf("diff (+got, -zero):")
	t.Printf(diff)
	return t.Attest()
}

// NotZero asserts that the value is non-zero.
func NotZero[T any](tb TB, got T, opts ...Option) bool {
	tb.Helper()
	var zero T
	t := newAttester(tb, opts...)
	if !t.Equal(got, zero) {
		return true
	}
	t.Printf("got is zero %T", got)
	return t.Attest()
}

// True asserts that a boolean is true.
func True(tb TB, got bool, opts ...Option) bool {
	tb.Helper()
	if got {
		return true
	}
	t := newAttester(tb, opts...)
	t.Printf("got false, want true")
	return t.Attest()
}

// False asserts that a boolean is false.
func False(tb TB, got bool, opts ...Option) bool {
	tb.Helper()
	if !got {
		return true
	}
	t := newAttester(tb, opts...)
	t.Printf("got true, want false")
	return t.Attest()
}

// Panics asserts that the function panics.
func Panics(tb TB, f func(), opts ...Option) (ret bool) {
	tb.Helper()
	defer func() {
		tb.Helper()
		t := newAttester(tb, opts...)
		if r := recover(); r == nil {
			t.Printf("no panic")
		}
		ret = t.Attest()
	}()
	f()
	return
}

// Approximately asserts that got is within delta of want. For example,
//   pi := float64(22)/7
//   Approximately(t, pi, 3.14, 0.01)
// asserts that our estimate of pi is between 3.13 and 3.15, exclusive.
//
// Approximately works with any type whose underlying type is numeric, so it
// also works with time.Duration.
func Approximately[T Number](tb TB, got, want, delta T, opts ...Option) bool {
	tb.Helper()
	lower := want - delta
	upper := want + delta
	if lower > upper {
		lower, upper = upper, lower
	}
	if got > lower && got < upper {
		return true
	}
	t := newAttester(tb, opts...)
	t.Printf("%v not within %v of %v", got, delta, want)
	return t.Attest()
}

// Contains asserts that a slice contains a target element.
func Contains[T any](tb TB, got []T, want T, opts ...Option) bool {
	tb.Helper()
	t := newAttester(tb, opts...)
	for _, v := range got {
		if t.Equal(v, want) {
			return true
		}
	}
	t.Printf("%v not in slice", want)
	return t.Attest()
}

// A Number is any type whose underlying type is one of Go's built-in integer
// or float types.
type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// TB is the subset of testing.TB that attest depends on. The standard
// library's *testing.T, B, and F types all implement it.
type TB interface {
	Helper()
	Errorf(string, ...any)
	Fatalf(string, ...any)
}

type attester struct {
	tb    TB
	fatal bool
	msg   string
	cmp   []cmp.Option
	buf   bytes.Buffer
}

func newAttester(tb TB, opts ...Option) *attester {
	t := &attester{tb: tb, fatal: true}
	for _, opt := range opts {
		opt.apply(t)
	}
	return t
}

func (t *attester) Equal(got, want any) bool {
	return cmp.Equal(got, want, t.cmp...)
}

func (t *attester) Diff(got, want any) string {
	return cmp.Diff(want, got, t.cmp...)
}

func (t *attester) Printf(tmpl string, args ...any) {
	line := fmt.Sprintf(tmpl, args...)
	if t.buf.Len() > 0 {
		t.buf.WriteString("\n")
	} else if t.msg != "" {
		line += ": " + t.msg
	}
	t.buf.WriteString(line)
}

func (t *attester) Attest() bool {
	t.tb.Helper()
	if t.buf.Len() == 0 {
		return true
	}
	if t.fatal {
		t.tb.Fatalf(t.buf.String())
	} else {
		t.tb.Errorf(t.buf.String())
	}
	return false
}
