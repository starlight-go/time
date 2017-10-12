package time

import (
	"testing"
	"time"

	"github.com/google/skylark"
)

func TestTimeAddDuration(t *testing.T) {
	tm := time.Now()
	fakeNow := func(_ *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
		return Time(tm), nil
	}

	globals := skylark.StringDict{
		"now":       skylark.NewBuiltin("now", fakeNow),
		"timedelta": skylark.NewBuiltin("timedelta", Delta),
	}
	err := skylark.ExecFile(new(skylark.Thread), "foo.sky", `output = now() + timedelta(seconds=5)`, globals)

	if err != nil {
		t.Fatal(err)
	}
	v, ok := globals["output"]
	if !ok {
		t.Fatal("missing output var in globals")
	}
	actual, ok := v.(Time)
	if !ok {
		t.Fatalf("expected type Time but was %T", v)
	}
	actT := time.Time(actual)
	exp := tm.Add(time.Second * 5)
	if !actT.Equal(tm.Add(time.Second * 5)) {
		t.Fatalf("Expected %v but got %v", exp, actT)
	}
}
