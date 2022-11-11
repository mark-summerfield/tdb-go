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

func Test001(t *testing.T) {
	expected := "Hello tdb v0.1.0\n"
	actual := Hello()
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
