// Package time provides time.Time wrappers for the skylark embedded language.
package time

import (
	"errors"
	"time"

	"github.com/google/skylark"
	"github.com/google/skylark/syntax"
)

// Time is the type of a Skylark time.Time.
type Time time.Time

// Zero is the zero time.
var Zero = Time(time.Time{})

// String returns the string representation of this Time.
func (t Time) String() string {
	return time.Time(t).String()
}

// Type returns the typename of the value.
func (t Time) Type() string { return "time" }

// Freeze makes the value immutable.
func (t Time) Freeze() {} // immutable

// Truth reports the python-truthiness of the value.
func (t Time) Truth() skylark.Bool { return skylark.Bool(t != Zero) }

// Hash returns a hash of the value for used in sorting.
func (t Time) Hash() (uint32, error) { return 0, errors.New("nope") }

// CompareSameType returns whether the two values are equal.
func (t Time) CompareSameType(op syntax.Token, y skylark.Value, depth int) (bool, error) {
	t2 := y.(Time)
	return time.Time(t).Equal(time.Time(t2)), nil
}
