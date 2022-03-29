package attest

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// An Option configures assertions.
type Option interface {
	apply(*attester)
}

type options struct {
	options []Option
}

// Options composes multiple Options into one. This may be useful if you're
// writing a helper package that bundles several options together, or if most
// assertions in your tests use a common set of options.
func Options(opts ...Option) Option {
	return &options{opts}
}

func (o *options) apply(t *attester) {
	for _, opt := range o.options {
		opt.apply(t)
	}
}

type msgOption struct {
	msg string
}

// Sprintf adds an explanation to the default failure message. If your tests
// make many similar assertions, the additional explanation may clarify the
// test output.
//
// Arguments are passed to fmt.Sprintf for formatting.
func Sprintf(template string, args ...any) Option {
	return &msgOption{fmt.Sprintf(template, args...)}
}

// Sprint adds an explanation to the default failure message. If your tests
// make many similar assertions, the additional explanation may clarify the
// test output.
//
// Arguments are passed to fmt.Sprint for formatting.
func Sprint(args ...any) Option {
	return &msgOption{fmt.Sprint(args...)}
}

func (o *msgOption) apply(t *attester) {
	t.msg = o.msg
}

type fatalOption struct {
	fatal bool
}

// Fatal stops the test immediately when an assertion fails. This is the
// default behavior, but Fatal may still be useful to reverse the effect of
// Continue.
func Fatal() Option {
	return &fatalOption{true}
}

// Continue allows the test to continue executing when an assertion fails.
// By default, failed assertions stop the test immediately.
func Continue() Option {
	return &fatalOption{false}
}

func (o *fatalOption) apply(t *attester) {
	t.fatal = o.fatal
}

type cmpOption struct {
	cmp []cmp.Option
}

// Cmp configures the underlying equality assertion, if any. See the
// github.com/google/go-cmp/cmp package documentation for an explanation of the
// default logic.
//
// In particular, note that the cmp package's default behavior is to panic when
// comparing structs with unexported fields. If you control the type in
// question, implement an Equal method and cmp will use it by default. If you
// don't control the type, use Allow or Comparer. If none of those approaches
// fit your needs, cmp and its cmpopts subpackage offer many other ways to
// relax this safety check.
//
// If you're comparing types generated from a Protocol Buffer schema,
// google.golang.org/protobuf/testing/protocmp's Transform() option
// safely transforms messages to a comparable, diffable type.
func Cmp(opts ...cmp.Option) Option {
	return &cmpOption{opts}
}

// Allow configures the underlying cmp package to forcibly introspect
// unexported fields of the specified struct types. By default, cmp panics when
// comparing structs with unexported fields. Allow panics if called with
// anything other than a struct type (for example, a pointer or slice).
//
// It's useful as a quick hack but is usually a bad idea: changes in the
// internals of some other package may break your tests. If you control the
// type in question, implement an Equal method instead. If you don't control
// the type, Comparer is usually safer.
func Allow(types ...any) Option {
	return Cmp(cmp.AllowUnexported(types...))
}

// Comparer configures the underlying cmp package to compare values of type T
// using the provided function. This is especially useful when comparing
// third-party types with unexported fields.
//
// The equality function must be symmetric (the order of the two arguments
// doesn't matter), deterministic (it always returns the same result), and pure
// (it may not mutate its arguments).
func Comparer[T any](equal func(T, T) bool) Option {
	return Cmp(cmp.Comparer(equal))
}

func (o *cmpOption) apply(t *attester) {
	t.cmp = append(t.cmp, o.cmp...)
}
