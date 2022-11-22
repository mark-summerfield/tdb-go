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
var CSV string

func TestCSV(t *testing.T) {
	db := makeCSV(t)
	raw, err := tdb.Marshal(db)
	if err != nil {
		t.Error(err)
	}
	if string(raw) != CSV {
		fmt.Println("======= Tdb Marshal Database ======")
		fmt.Print(string(raw))
		_ = os.WriteFile("/tmp/1", raw, 0666)
		_ = os.WriteFile("/tmp/2", []byte(CSV), 0666)
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
	if string(raw2) != CSV {
		fmt.Println("====== Tdb Unmarshal Database =====")
		fmt.Print(string(raw2))
		_ = os.WriteFile("/tmp/3", raw2, 0666)
		_ = os.WriteFile("/tmp/4", []byte(CSV), 0666)
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
	Quantity int
	ID       string
	Desc     string `tdb:"Description"`
}

func makeCSV(t *testing.T) csvDatabase {
	return csvDatabase{
		PriceList: []Price{
			{time.Date(2022, time.September, 21, 0, 0, 0, 0, time.UTC),
				3.99, 2, "CH1-A2", "Chisels (pair), 1in & 1¼in"},
			{time.Date(2022, time.October, 2, 0, 0, 0, 0, time.UTC),
				4.49, 1, "HV2-K9", "Hammer, 2lb"},
			{time.Date(2022, time.October, 2, 0, 0, 0, 0, time.UTC),
				5.89, 1, "SX4-D1", "Eversure Sealant, 13-floz"},
			{time.Date(2022, time.November, 13, 0, 0, 0, 0, time.UTC),
				8.49, tdb.IntSentinal, "PV7-X2", tdb.StrSentinal},
		},
	}
}
