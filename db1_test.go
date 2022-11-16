package tdb_test

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"github.com/mark-summerfield/tdb"
	"testing"
	"time"
)

//go:embed eg/db1.tdb
var DbEg1Text string

func TestDb1(t *testing.T) {
	db := makeDb(t)
	fmt.Println("======= Tdb Go Data ======")
	fmt.Println(db)
	fmt.Println("======= Tdb Marshal Database ======")
	raw, err := tdb.Marshal(db)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(raw))
	if string(raw) != DbEg1Text {
		t.Error("Database: raw != text")
	}
	/*
		fmt.Println("======= Tdb example data ======")
		raw, err := tdb.Marshal(db)
		if err != nil {
			t.Error(err)
		}
		if string(raw) != DbEg1Text {
			t.Error("raw != text")
		}
		fmt.Println(string(raw))
		fmt.Println("======= Tdb example text ======")
		raw = []byte(DbEg1Text)
		fmt.Println(string(raw))
		var database Database
		err = tdb.Unmarshal(raw, &database)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("========== Database ===========")
		fmt.Println(database)
		raw, err = tdb.Marshal(database)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("======= Tdb result text =======")
		fmt.Println(string(raw))
		fmt.Println("===============================")
	*/
}

type Database struct {
	Customers []Customer
	Invoices  []Invoice
	LineItems []LineItem
}

type Customer struct {
	Cid     int
	Company string
	Address string
	Contact string
	Email   string
	Icon    []byte
}

type Invoice struct {
	Inum   int
	Cid    int
	Raised time.Time `tdb:"date"`
	Due    time.Time `tdb:"date"`
	Paid   bool
	Desc   string
}

type LineItem struct {
	Liid      int
	Inum      int
	Delivered time.Time `tdb:"date"`
	UnitPrice float64
	Quantity  int
	Desc      string
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
