package attest_test

import (
	"fmt"
	"time"

	"github.com/akshayjshah/attest"
	"github.com/stretchr/testify/assert"
)

type logTB struct{}

func (_ *logTB) Helper()                         {}
func (_ *logTB) Errorf(tmpl string, args ...any) { fmt.Printf("ERROR: "+tmpl, args...) }
func (_ *logTB) Fatalf(tmpl string, args ...any) { fmt.Printf("FATAL: "+tmpl, args...) }

func ExampleAllow() {
	type point struct {
		x, y float64
	}
	attest.Equal(
		&logTB{}, // no testing.T in examples :(
		point{1.0, 1.0},
		point{1.0, 1.0},
		// Without Allow, the underlying cmp library panics because point has
		// unexported fields. We could also use Comparer, or we could implement an
		// Equal method on point.
		attest.Allow(point{}),
	)
	// Output:
}

func ExampleComparer() {
	type point struct {
		x, y float64
	}
	attest.Equal(
		&logTB{}, // no testing.T in examples :(
		point{1.0, 1.0},
		point{1.0, 1.0},
		// Without Comparer, the underlying cmp library panics because point has
		// unexported fields. We could also use Allow, or we could implement an
		// Equal method on point.
		attest.Comparer(func(left, right point) bool {
			return left.x == right.x && left.y == right.y
		}),
	)
	// Output:
}

func ExampleSprint() {
	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)
	attest.False(
		&logTB{}, // no testing.T in examples :(
		today.Before(tomorrow),
		attest.Sprint("alas, time", " marches on"),
	)
	// Output:
	// FATAL: got true, want false: alas, time marches on
}

func ExampleSprintf() {
	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)
	attest.False(
		&logTB{}, // no testing.T in examples :(
		today.Before(tomorrow),
		attest.Sprintf("%s, time marches on", "alas"),
	)
	// Output:
	// FATAL: got true, want false: alas, time marches on
}

func ExampleOptions() {
	type point struct {
		x, y float64
	}
	// If all our tests have some options in common, it's nice to extract them
	// into a named bundle.
	defaults := attest.Options(
		attest.Continue(),
		attest.Comparer(func(left, right point) bool {
			return left.x == right.x && left.y == right.y
		}),
	)
	// We can reuse our default options in each test. We can also specify more
	// options without an ugly cascade of appends.
	assert.Zero(
		&logTB{}, // no testing.T in examples :(
		point{},
		defaults,       // our defaults
		attest.Fatal(), // override Continue() from defaults
	)
}
