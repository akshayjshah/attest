package attest

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
)

type mockTB struct {
	fatal bool
	out   string
}

func (m *mockTB) Helper() {}

func (m *mockTB) Errorf(tmpl string, args ...any) {
	m.fatal = false
	m.out = fmt.Sprintf(tmpl, args...)
}

func (m *mockTB) Fatalf(tmpl string, args ...any) {
	m.fatal = true
	m.out = fmt.Sprintf(tmpl, args...)
}

func (m *mockTB) AssertError(t testing.TB) {
	if m.out == "" {
		t.Fatalf("expected failure")
	}
	if m.fatal {
		t.Fatalf("expected error, got fatal")
	}
	m.clear()
}

func (m *mockTB) AssertFatal(t testing.TB) {
	if m.out == "" {
		t.Fatalf("expected failure")
	}
	if !m.fatal {
		t.Fatalf("expected fatal, got error")
	}
	m.clear()
}

func (m *mockTB) clear() {
	m.fatal = false
	m.out = ""
}

func TestEqual(t *testing.T) {
	Equal(t, 1, 1)

	var mock mockTB
	Equal(&mock, 1, 2)
	mock.AssertFatal(t)
}

func TestError(t *testing.T) {
	Ok(t, nil)
	var err error
	Ok(t, err)
	Error(t, errors.New("foo"))
	ErrorIs(t, fmt.Errorf("something: %w", io.EOF), io.EOF)

	var mock mockTB
	Ok(&mock, errors.New("foo"))
	mock.AssertFatal(t)
	Error(&mock, nil)
	mock.AssertFatal(t)
	ErrorIs(&mock, fmt.Errorf("something: %v", io.EOF), io.EOF)
	mock.AssertFatal(t)
}

func TestZero(t *testing.T) {
	var n *int
	Zero(t, n)
	var s []int
	Zero(t, s)
	var m map[string]string
	Zero(t, m)
	NotZero(t, 3)

	var mock mockTB
	Zero(&mock, 3)
	mock.AssertFatal(t)
	NotZero(&mock, 0)
	mock.AssertFatal(t)
}

func TestBool(t *testing.T) {
	True(t, true)
	False(t, false)

	var mock mockTB
	True(&mock, false)
	mock.AssertFatal(t)
	False(&mock, true)
	mock.AssertFatal(t)
}

func TestPanics(t *testing.T) {
	Panics(t, func() { panic("oh no") })

	var mock mockTB
	Panics(&mock, func() {})
	mock.AssertFatal(t)
}

func TestApproximately(t *testing.T) {
	Approximately(t, 3.0, 3.05, 0.1)
	Approximately(t, 11, 10, -3)

	var mock mockTB
	Approximately(&mock, 3.0, 3.05, 0.01)
	mock.AssertFatal(t)
	Approximately(&mock, 3.0, 3.0, math.NaN())
	mock.AssertFatal(t)
}

func TestContains(t *testing.T) {
	Contains(t, []int{0, 1, 2}, 1)

	var mock mockTB
	Contains(&mock, []int{0, 1}, 2)
	mock.AssertFatal(t)
}

func TestAllow(t *testing.T) {
	type pair struct {
		first, second int
	}

	var p pair
	Zero(t, p, Allow(p))
	var null *pair
	Zero(t, null, Allow(pair{}))
}

func TestComparer(t *testing.T) {
	type mod3 int
	Equal(t, mod3(3), mod3(6), Comparer(func(x, y mod3) bool {
		return x%3 == y%3
	}))
}

func TestSprint(t *testing.T) {
	var mock mockTB
	True(&mock, false, Sprint("a", " message"))
	if !strings.HasSuffix(mock.out, ": a message") {
		t.Errorf("no user-supplied message in output")
	}
	mock.AssertFatal(t)
}

func TestSprintf(t *testing.T) {
	var mock mockTB
	True(&mock, false, Sprintf("%s %s", "a", "message"))
	if !strings.HasSuffix(mock.out, ": a message") {
		t.Errorf("no user-supplied message in output")
	}
	mock.AssertFatal(t)
}

func TestContinue(t *testing.T) {
	var mock mockTB
	True(&mock, false, Continue())
	mock.AssertError(t)
}

func TestFatal(t *testing.T) {
	var mock mockTB
	True(&mock, false, Continue(), Fatal())
	mock.AssertFatal(t)
}

func TestOptions(t *testing.T) {
	var mock mockTB
	True(&mock, false, Options(Continue(), Fatal()))
	mock.AssertFatal(t)
}
