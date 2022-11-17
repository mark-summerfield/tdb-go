package tdb_test

import (
	"fmt"
	"github.com/mark-summerfield/tdb"
	"regexp"
	"testing"
	// "github.com/mark-summerfield/gong"
	// "golang.org/x/exp/maps"
	// "golang.org/x/exp/slices"
)

// maps.Equal() & maps.EqualFunc() & slices.Equal() & slices.EqualFunc()
// https://pkg.go.dev/golang.org/x/exp/maps
// https://pkg.go.dev/golang.org/x/exp/slices
// gong.IsRealClose() & gong.IsRealZero()

func compare(n int, raw []byte, expected string, t *testing.T) {
	actual := string(raw)
	if actual != expected {
		t.Errorf("\nTest%03d:\nexpected: %q !=\nactual:   %q", n, expected,
			actual)
	}
}

func expectError(n int, code int, err error, t *testing.T) {
	if err == nil {
		t.Errorf("Test%03d: expected e%d#…", n, code)
	} else {
		e := err.Error()
		found, _ := regexp.MatchString(fmt.Sprintf("^e%d#", code), e)
		if !found {
			t.Errorf("Test%03d: expected e%d#…, got %s", n, code, e)
		}
	}
}

type OneRecord struct {
	OneField int
}
type OneDatabase struct {
	OneTable []OneRecord
}

func Test001(t *testing.T) {
	d := OneDatabase{OneTable: []OneRecord{{2}, {3}, {5}}}
	raw, err := tdb.Marshal(d)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	compare(1, raw, "[OneTable OneField int\n%\n2\n3\n5\n]\n", t)
}

func Test002(t *testing.T) {
	d := OneDatabase{}
	_, err := tdb.Marshal(d)
	expectError(2, tdb.CannotMarshalEmpty, err, t)
}

func Test003(t *testing.T) {
	d := OneDatabase{OneTable: []OneRecord{}}
	_, err := tdb.Marshal(d)
	expectError(3, tdb.CannotMarshalEmpty, err, t)
}
