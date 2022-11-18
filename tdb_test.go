package tdb

import (
	"fmt"
	"regexp"
	"testing"
)

func compare(n int, raw []byte, expected string, t *testing.T) {
	actual := string(raw)
	if actual != expected {
		t.Errorf("\nTest%03d:\nexpected: %q !=\nactual:   %q", n, expected,
			actual)
	}
}

func expectError(code int, err error, t *testing.T) {
	if err == nil {
		t.Errorf("TestE%03d: expected e%d#…", code, code)
	} else {
		e := err.Error()
		found, _ := regexp.MatchString(fmt.Sprintf("^e%d#", code), e)
		if !found {
			t.Errorf("TestE%03d: expected e%d#…, got %s", code, code, e)
		}
	}
}

type Record struct {
	AField int
}
type ADatabase struct {
	Records []Record
}

func Test001(t *testing.T) {
	d := ADatabase{Records: []Record{{2}, {3}, {5}}}
	raw, err := Marshal(d)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	compare(1, raw, "[Records AField int\n%\n2\n3\n5\n]\n", t)
}

func TestE100(t *testing.T) {
	type ADatabase struct {
		ATable string
	}
	d := ADatabase{"one"}
	_, err := Marshal(d)
	expectError(e100, err, t)
}

func TestE101(t *testing.T) {
	d := "duh"
	_, err := Marshal(d)
	expectError(e101, err, t)
}

func TestE102(t *testing.T) {
	d := ADatabase{}
	_, err := Marshal(d)
	expectError(e102, err, t)
	d = ADatabase{Records: []Record{}}
	_, err = Marshal(d)
	expectError(e102, err, t)
}

func TestE103(t *testing.T) {
	type ARecord struct {
		Names []string
	}
	type ADatabase struct {
		ATable []ARecord
	}
	a := ADatabase{
		ATable: []ARecord{
			{Names: []string{"one", "two"}},
		},
	}
	_, err := Marshal(a)
	expectError(e103, err, t)
	type BRecord struct {
		Items []complex64
	}
	type BDatabase struct {
		BTable []BRecord
	}
	b := BDatabase{
		BTable: []BRecord{
			{Items: []complex64{2 + 0.5i}},
		},
	}
	_, err = Marshal(b)
	expectError(e103, err, t)
}

func TestE104(t *testing.T) {
	type ARecord struct {
		Items complex64
	}
	type ADatabase struct {
		ATable []ARecord
	}
	d := ADatabase{
		ATable: []ARecord{
			{Items: 2 + 0.5i},
		},
	}
	_, err := Marshal(d)
	expectError(e104, err, t)
}
