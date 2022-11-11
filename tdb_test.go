package tdb

import (
	"testing"
)

/*
func realEqual(x, y float64) bool {
    return math.Abs(x-y) < 0.0001
}
*/

/*
func expectEqualSlice(expected, actuals []string, what string) string {
    if !reflect.DeepEqual(actuals, expected) {
        return fmt.Sprintf("expected %s=%s, got %s", what, expected,
            actuals)
    }
    return ""
}
*/

/*
func expectEmptySlice(slice []string, what string) string {
    if slice != nil {
        return fmt.Sprintf("expected %s=nil, got %s", what, slice)
    }
    return ""
}
*/

type T1 struct {
	ft int
}

/*
func Test001(t *testing.T) {
	expected := "Tdb1\n[T1 f1 int]"
	var
	actual := tdb.Unmarshal(
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
*/
