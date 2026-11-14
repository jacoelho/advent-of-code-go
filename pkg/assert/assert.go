// Package assert provides minimal assertions for tests.
package assert

import "testing"

// Equal reports a test error when actual differs from expected.
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// NoError fails the test immediately when err is non-nil.
func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
