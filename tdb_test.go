package tdb

import (
	_ "embed"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func compare(name string, raw []byte, expected string, t *testing.T) {
	actual := string(raw)
	if actual != expected {
		actual = strings.TrimSpace(actual)
		expected = strings.TrimSpace(expected)
		if actual != expected {
			t.Errorf("\nTest%s:\nexpected: %q !=\nactual:   %q", name,
				expected, actual)
		}
	}
}

func expectError(code int, err error, t *testing.T) {
	if err == nil {
		t.Errorf("TestE%03d: expected e%d#…, got nil", code, code)
	} else {
		e := err.Error()
		found, _ := regexp.MatchString(fmt.Sprintf("^e%d#", code), e)
		if !found {
			t.Errorf("TestE%03d: expected e%d#…, got %s", code, code, e)
		}
	}
}

type Database struct {
	Places    []Place
	LineItems []Item `tdb:"Items"`
}

type Place struct {
	Pid  int `tdb:"PID"`
	Name string
}

type Item struct {
	Iid         int       `tdb:"IID"`
	Description string    `tdb:"Desc:str"`
	When        time.Time `tdb:"date"`
}

type Record struct {
	AField int
}

type ADatabase struct {
	Records []Record
}

func Test001(t *testing.T) {
	d1 := ADatabase{Records: []Record{{2}, {3}, {5}}}
	raw, err := Marshal(d1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	compare("001", raw, "[Records AField int\n%\n2\n3\n5\n]\n", t)
	d2 := ADatabase{}
	err = Unmarshal(raw, &d2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("unexpectedly unequal:\nONE: %v\nTWO: %v", d1, d2)
	}
}

func Test002(t *testing.T) {
	type DBA struct {
		Places []Place
	}
	data := `[Places PID int Name str
%
801 <One>
802 <Two>
]`
	db := DBA{}
	if err := Unmarshal([]byte(data), &db); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	raw, err := Marshal(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	compare("002", raw, data, t)
}

func Test003(t *testing.T) {
	type Rec struct {
		Ok bool
	}
	type DBA struct {
		Recs []Rec
	}
	data := "[Recs Ok bool\n%\nT f y N 1 0]"
	expected := "[Recs Ok bool\n%\nT\nF\nT\nF\nT\nF\n]"
	db := DBA{}
	if err := Unmarshal([]byte(data), &db); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	raw, err := Marshal(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	compare("003", raw, expected, t)
}

//go:embed eg/incidents.tdb
var Incidents string

func TestIncidents(t *testing.T) {
	type Incident struct {
		Report_ID                   string
		Date                        time.Time `tdb:"date"`
		Aircraft_ID                 string
		Aircraft_Type               string
		Pilot_Percent_Hours_on_Type float32
		Pilot_Total_Hours           int
		MidAir                      bool
		Airport                     string
		Narrative                   string
	}

	type IncidentsDb struct {
		Aircraft_Incidents []Incident
	}

	incidents := IncidentsDb{}
	if err := Unmarshal([]byte(Incidents), &incidents); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	raw, err := MarshalDecimals(incidents, 4)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	compare("Incidents", raw, Incidents, t)
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

func TestE107(t *testing.T) {
	type Record struct {
		F int
	}
	type Database struct {
		Records []Record
	}
	db := Database{}
	raw := []byte("[T F int%")
	err := Unmarshal(raw, &db)
	expectError(e107, err, t)
}

func TestE108(t *testing.T) {
	type Record struct {
		F int
	}
	type Database struct {
		Records []Record
	}
	db := Database{}
	raw := []byte("[T F int%]")
	err := Unmarshal(raw, db)
	expectError(e108, err, t)
}

func TestE109(t *testing.T) {
	type Record struct {
		F int
	}
	db := make([]Record, 0)
	raw := []byte("[T F int%]")
	err := Unmarshal(raw, &db)
	expectError(e109, err, t)
}

func TestE110(t *testing.T) {
	type Record struct {
		F []byte
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bytes\n%\n(20AC\n]")
	err := Unmarshal(raw, &db)
	expectError(e110, err, t)
}

func TestE112(t *testing.T) {
	type Record struct {
		F []byte
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F uint\n%\n(20AC)\n]")
	err := Unmarshal(raw, &db)
	expectError(e112, err, t)
}

func TestE114(t *testing.T) {
	type Record struct {
		F int
		G int
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F int G int\n%\n1 2\n3 F\n]")
	err := Unmarshal(raw, &db)
	expectError(e114, err, t)
}

func TestE115(t *testing.T) {
	type Record struct {
		F bool
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bool\n%\nT !\n]")
	err := Unmarshal(raw, &db)
	expectError(e115, err, t)
}

func TestE116(t *testing.T) {
	type Record struct {
		F bool
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bool\n%\nT ()\n]")
	err := Unmarshal(raw, &db)
	expectError(e116, err, t)
}

func TestE117(t *testing.T) {
	type Record struct {
		F bool
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bool\n%\nT <>\n]")
	err := Unmarshal(raw, &db)
	expectError(e117, err, t)
}

func TestE118(t *testing.T) {
	type Record struct {
		F bool
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bool\n%\nT -1\n]")
	err := Unmarshal(raw, &db)
	expectError(e118, err, t)
}

func TestE120(t *testing.T) {
	type Record struct {
		F int
		G int
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F int G int\n%\n1 2\n3\n]")
	err := Unmarshal(raw, &db)
	expectError(e120, err, t)
}

func TestE121(t *testing.T) {
	type Record struct {
		F bool
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bool\n%\nT x\n]")
	err := Unmarshal(raw, &db)
	expectError(e121, err, t)
}

func TestE122(t *testing.T) {
	type Record struct {
		f int
		G int
	}
	type Database struct {
		T []Record
	}
	db := Database{[]Record{{f: 1, G: 2}}} // filled in purely for linting
	raw := []byte("[T f int G int\n%\n1 2\n3 4\n]")
	err := Unmarshal(raw, &db)
	expectError(e122, err, t)
}

func TestE123(t *testing.T) {
	type Record struct {
		F []byte
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bytes\n%\n(20AC) (EF1G)\n]")
	err := Unmarshal(raw, &db)
	expectError(e123, err, t)
}

func TestE124(t *testing.T) {
	type Record struct {
		F int
		G int
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F int G int\n%\n1 2 3 4")
	err := Unmarshal(raw, &db)
	expectError(e124, err, t)
}

func TestE125(t *testing.T) {
	type Record struct {
		F int
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F int\n%\n1 1-0\n]")
	err := Unmarshal(raw, &db)
	expectError(e125, err, t)
}

func TestE126(t *testing.T) {
	type Record struct {
		F float32
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F real\n%\n1 1-0\n]")
	err := Unmarshal(raw, &db)
	expectError(e126, err, t)
}

func TestE127(t *testing.T) {
	type Record struct {
		F time.Time
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F date\n%\n2020-1-9-3\n]")
	err := Unmarshal(raw, &db)
	expectError(e127, err, t)
}

func TestE129(t *testing.T) {
	type Record struct {
		F []byte
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F\n%\n(20AC)\n]")
	err := Unmarshal(raw, &db)
	expectError(e129, err, t)
}

func TestE130(t *testing.T) {
	type Record struct {
		F bool
	}
	type Database struct {
		T []Record
	}
	db := Database{}
	raw := []byte("[T F bool\n%\nT 2\n]")
	err := Unmarshal(raw, &db)
	expectError(e130, err, t)
}
