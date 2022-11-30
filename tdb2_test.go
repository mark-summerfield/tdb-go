package tdb_test

import (
	_ "embed"
	"fmt"
	tdb "github.com/mark-summerfield/tdb-go"
	"os"
	"testing"
	"time"
)

//go:embed eg/csv.tdb
var Csv string

func TestCsv(t *testing.T) {
	db := makeCSV(t)
	raw, err := tdb.Marshal(db)
	if err != nil {
		t.Error(err)
	}
	if string(raw) != Csv {
		fmt.Println("======= Tdb Marshal Database ======")
		fmt.Print(string(raw))
		_ = os.WriteFile("/tmp/1", raw, 0666)
		_ = os.WriteFile("/tmp/2", []byte(Csv), 0666)
		fmt.Println("wrote /tmp/[12] (actual/expected)")
		fmt.Println("===================================")
		t.Error("Database: raw != text")
	}
	var database csvDatabase
	err = tdb.Unmarshal(raw, &database)
	if err != nil {
		t.Error(err)
	}
	raw2, err := tdb.Marshal(database)
	if err != nil {
		t.Error(err)
	}
	if string(raw2) != Csv {
		fmt.Println("====== Tdb Unmarshal Database =====")
		fmt.Print(string(raw2))
		_ = os.WriteFile("/tmp/3", raw2, 0666)
		_ = os.WriteFile("/tmp/4", []byte(Csv), 0666)
		fmt.Println("wrote /tmp/[34] (actual/expected)")
		fmt.Println("===================================")
		t.Error("Database: raw2 != text")
	}
}

type csvDatabase struct {
	PriceList []Price
}

type Price struct {
	Date     time.Time `tdb:"date"`
	Price    float64
	Quantity *int
	ID       string
	Desc     *string `tdb:"Description"`
}

func makeCSV(t *testing.T) csvDatabase {
	db := csvDatabase{
		PriceList: []Price{
			{time.Date(2022, time.September, 21, 0, 0, 0, 0, time.UTC),
				3.99, nil, "CH1-A2", nil},
			{time.Date(2022, time.October, 2, 0, 0, 0, 0, time.UTC),
				4.49, nil, "HV2-K9", nil},
			{time.Date(2022, time.October, 2, 0, 0, 0, 0, time.UTC),
				5.89, nil, "SX4-D1", nil},
			{time.Date(2022, time.November, 13, 0, 0, 0, 0, time.UTC),
				8.49, nil, "PV7-X2", nil},
		},
	}
	q0 := 2
	d0 := "Chisels (pair), 1in & 1¼in"
	db.PriceList[0].Quantity = &q0
	db.PriceList[0].Desc = &d0
	q1 := 1
	d1 := "Hammer, 2lb"
	db.PriceList[1].Quantity = &q1
	db.PriceList[1].Desc = &d1
	q2 := 1
	d2 := "Eversure Sealant, 13-floz"
	db.PriceList[2].Quantity = &q2
	db.PriceList[2].Desc = &d2
	return db
}
