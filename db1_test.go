package tdb_test

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"github.com/mark-summerfield/tdb"
	"os"
	"testing"
	"time"
)

//go:embed eg/db1.tdb
var DbEg1Text string

func TestDb1(t *testing.T) {
	db := makeDb(t)
	raw, err := tdb.Marshal(db)
	if err != nil {
		t.Error(err)
	}
	if string(raw) != DbEg1Text {
		fmt.Println("======= Tdb Marshal Database ======")
		fmt.Print(string(raw))
		_ = os.WriteFile("/tmp/1", raw, 0666)
		_ = os.WriteFile("/tmp/2", []byte(DbEg1Text), 0666)
		fmt.Println("wrote /tmp/[12] (actual/expected)")
		fmt.Println("===================================")
		t.Error("Database: raw != text")
	}
	var database Database
	err = tdb.Unmarshal(raw, &database)
	if err != nil {
		t.Error(err)
	}
	raw2, err := tdb.Marshal(database)
	if err != nil {
		t.Error(err)
	}
	if string(raw2) != DbEg1Text {
		fmt.Println("====== Tdb Unmarshal Database =====")
		fmt.Print(string(raw2))
		_ = os.WriteFile("/tmp/3", raw2, 0666)
		_ = os.WriteFile("/tmp/4", []byte(DbEg1Text), 0666)
		fmt.Println("wrote /tmp/[34] (actual/expected)")
		fmt.Println("===================================")
		t.Error("Database: raw2 != text")
	}
}

type Database struct {
	Customers []Customer
	Invoices  []Invoice
	LineItems []LineItem `tdb:"Items"`
}

type Customer struct {
	Cid     int `tdb:"CID"`
	Company string
	Address string
	Contact string
	Email   string
	Icon    []byte
}

type Invoice struct {
	Inum   int       `tdb:"INUM"`
	Cid    int       `tdb:"CID"`
	Raised time.Time `tdb:"Raised_Date:date"`
	Due    time.Time `tdb:"Due_Date:date"`
	Paid   bool
	Desc   string `tdb:"Description"`
}

type LineItem struct {
	Liid      int       `tdb:"LIID"`
	Inum      int       `tdb:"INUM"`
	Delivered time.Time `tdb:"Delivery_Date:date"`
	UnitPrice float64   `tdb:"Unit_Price"`
	Quantity  int
	Desc      string `tdb:"Description"`
}

func makeDb(t *testing.T) Database {
	icon, err := hex.DecodeString("89504E470D0A1A0A0000000D494844520000000C0000000C080600000056755CE7000000097048597300000EC400000EC401952B0E1B0000002849444154289163646068F8CF4002602245317D34B0600A3530621183FB7310FA81640D8C832FE2002C7F051786CBFA670000000049454E44AE426082")
	if err != nil {
		t.Error(err)
	}
	db := Database{
		Customers: []Customer{
			{50, "Best People", "123 Somewhere", "John Doe", "j@doe.com",
				icon},
			{19, "Supersuppliers", tdb.StrSentinal, "Jane Doe",
				"jane@super.com", tdb.BytesSentinal},
		},
		Invoices: []Invoice{
			{152, 50, time.Date(2022,
				time.January, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2022, time.February, 17, 0, 0, 0, 0, time.UTC),
				false, "COD"},
			{153, 19,
				time.Date(2022, time.January, 19, 0, 0, 0, 0, time.UTC),
				time.Date(2022, time.February, 19, 0, 0, 0, 0, time.UTC),
				true, tdb.StrSentinal},
		},
		LineItems: []LineItem{
			{1839, 152,
				time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
				29.99, 2, "Bales of <hay>"},
			{1840, 152,
				time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
				5.98, 3, "Straps & Things"},
			{1620, 153,
				time.Date(2022, time.January, 19, 0, 0, 0, 0, time.UTC),
				11.5, 1, "Washers (1\")"},
		}}
	return db
}
