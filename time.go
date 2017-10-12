// Package time provides time.Time wrappers for the skylark embedded language.
package time

import (
	"errors"
	"time"

	"github.com/google/skylark"
	"github.com/google/skylark/syntax"
)

// Now returns the current Time.
func Now(_ *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("too many args")
	}
	if len(kwargs) != 0 {
		return nil, errors.New("too many kwargs")
	}
	return Time(time.Now()), nil
}

// Delta returns a duration created from kwargs.  Expected "hours", "minutes",
// "seconds", "milliseconds", or "nanoseconds", and assigned an int.
func Delta(_ *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("too many args")
	}
	var d time.Duration
	for _, t := range kwargs {
		if len(t) != 2 {
			panic("invalid kwarg")
		}
		s, ok := t[0].(skylark.String)
		if !ok {
			panic("invalid kwarg name")
		}
		i, ok := t[1].(skylark.Int)
		if !ok {
			return nil, errors.New("invalid value for timedelta arg, must be int")
		}
		v, ok := i.Int64()
		if !ok {
			return nil, errors.New("numeric value overflows int64")
		}
		switch s {
		case "hours":
			d += time.Hour * time.Duration(v)
		case "minutes":
			d += time.Minute * time.Duration(v)
		case "seconds":
			d += time.Second * time.Duration(v)
		case "milliseconds":
			d += time.Millisecond * time.Duration(v)
		case "nanoseconds":
			d += time.Nanosecond * time.Duration(v)
		default:
			return nil, errors.New("invalid duration unit: " + string(s))
		}
	}
	return Duration(d), nil
}

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
func (t Time) Truth() skylark.Bool { return skylark.Bool(!time.Time(t).IsZero()) }

// Hash returns a hash of the value for used in sorting.
func (t Time) Hash() (uint32, error) { return hashString(t.String()), nil }

// CompareSameType compares this Time to another.
func (t Time) CompareSameType(op syntax.Token, y skylark.Value, depth int) (bool, error) {
	t1 := time.Time(t)
	t2 := time.Time(y.(Time))

	switch op {
	case syntax.EQL:
		return t1.Equal(t2), nil
	case syntax.NEQ:
		return !t1.Equal(t2), nil
	case syntax.LE:
		return t1.Equal(t2) || t1.Before(t2), nil
	case syntax.LT:
		return t1.Before(t2), nil
	case syntax.GE:
		return t1.Equal(t2) || t1.After(t2), nil
	case syntax.GT:
		return t1.After(t2), nil
	}
	panic(op)
}

// Binary implements binary operators.
func (t Time) Binary(op syntax.Token, y skylark.Value, side skylark.Side) (skylark.Value, error) {
	if side == skylark.Right {
		return nil, nil
	}
	if time2, ok := y.(Time); ok {
		return t.binaryTime(op, time2)
	}
	if d, ok := y.(Duration); ok {
		return t.binaryDuration(op, d)
	}

	return nil, nil
}

func (t Time) binaryTime(op syntax.Token, y Time) (skylark.Value, error) {
	switch op {
	case syntax.MINUS:
		return Duration(time.Time(t).Sub(time.Time(y))), nil
	}
	return nil, nil
}

func (t Time) binaryDuration(op syntax.Token, y Duration) (skylark.Value, error) {
	switch op {
	case syntax.MINUS:
		return Time(time.Time(t).Add(-time.Duration(y))), nil
	case syntax.PLUS:
		return Time(time.Time(t).Add(time.Duration(y))), nil
	}
	return nil, nil
}

// hashString computes the FNV hash of s.  Copied from github.com/google/skylark/hashtable.go
func hashString(s string) uint32 {
	var h uint32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// Duration represents a span of time.
type Duration time.Duration

// String returns the string representation of this Time.
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Type returns the typename of the value.
func (d Duration) Type() string { return "time" }

// Freeze makes the value immutable.
func (d Duration) Freeze() {} // immutable

// Truth reports the python-truthiness of the value.
func (d Duration) Truth() skylark.Bool { return skylark.Bool(d != 0) }

// Hash returns a hash of the value for used in sorting.
func (d Duration) Hash() (uint32, error) { return hashString(d.String()), nil }

// CompareSameType compares this Time to another.
func (d Duration) CompareSameType(op syntax.Token, y skylark.Value, depth int) (bool, error) {
	d2 := y.(Duration)

	switch op {
	case syntax.EQL:
		return d == d2, nil
	case syntax.NEQ:
		return d != d2, nil
	case syntax.LE:
		return d <= d2, nil
	case syntax.LT:
		return d < d2, nil
	case syntax.GE:
		return d >= d2, nil
	case syntax.GT:
		return d > d2, nil
	}
	panic(op)
}

// Binary implements binary operators.
func (d Duration) Binary(op syntax.Token, y skylark.Value, side skylark.Side) (skylark.Value, error) {
	if side == skylark.Right {
		return nil, nil
	}
	d2, ok := y.(Duration)
	if !ok {
		return nil, nil
	}
	switch op {
	case syntax.MINUS:
		return Duration(d - d2), nil
	case syntax.PLUS:
		return Duration(d + d2), nil
	}
	return nil, nil
}
