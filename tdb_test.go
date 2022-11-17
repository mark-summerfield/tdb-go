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

func Test002(t *testing.T) {
	d := ADatabase{}
	_, err := Marshal(d)
	expectError(2, CannotMarshalEmpty, err, t)
}

func Test003(t *testing.T) {
	d := ADatabase{Records: []Record{}}
	_, err := Marshal(d)
	expectError(3, CannotMarshalEmpty, err, t)
}

func Test004(t *testing.T) {
	type ADatabase struct {
		ATable string
	}
	d := ADatabase{"one"}
	_, err := Marshal(d)
	expectError(3, CannotMarshalOuter, err, t)
}

func Test005(t *testing.T) {
	d := "duh"
	_, err := Marshal(d)
	expectError(3, CannotMarshal, err, t)
}

func Test006(t *testing.T) {
	type ARecord struct {
		Names []string
	}
	type ADatabase struct {
		ATable []ARecord
	}
	d := ADatabase{
		ATable: []ARecord{
			{Names: []string{"one", "two"}},
		},
	}
	_, err := Marshal(d)
	expectError(3, InvalidSliceType, err, t)
}

func Test007(t *testing.T) {
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
	expectError(3, InvalidFieldType, err, t)
}

func Test008(t *testing.T) {
	type ARecord struct {
		Items []complex64
	}
	type ADatabase struct {
		ATable []ARecord
	}
	d := ADatabase{
		ATable: []ARecord{
			{Items: []complex64{2 + 0.5i}},
		},
	}
	_, err := Marshal(d)
	expectError(3, InvalidSliceType, err, t)
}
