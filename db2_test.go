package tdb_test

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"github.com/mark-summerfield/tdb"
	"testing"
	"time"
	// "github.com/mark-summerfield/gong"
	// "golang.org/x/exp/maps"
	// "golang.org/x/exp/slices"
)

// maps.Equal() & maps.EqualFunc() & slices.Equal() & slices.EqualFunc()
// https://pkg.go.dev/golang.org/x/exp/maps
// https://pkg.go.dev/golang.org/x/exp/slices
// gong.IsRealClose() & gong.IsRealZero()

//go:embed eg/db2.tdb
var DbEg2Text string

func TestDb2(t *testing.T) {
	db := makeDb(t)
	fmt.Println("======= Tdb example data ======")
	fmt.Println(db)
	check(&db)
	raw, err := tdb.Marshal(db)
	if err != nil {
		t.Error(err)
	}
	if string(raw) != DbEg2Text {
		t.Error("raw != text")
	}
	fmt.Println("======= Tdb example text ======")
	raw = []byte(DbEg2Text)
	fmt.Println(string(raw))
	var database Database
	err = tdb.Unmarshal(raw, &database)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("========== Database ===========")
	fmt.Println(database)
	check(&database)
	raw, err = tdb.Marshal(database)
	if err != nil {
		t.Error(err)
	}
	text := string(raw)
	fmt.Println("======= Tdb result text =======")
	fmt.Println(text)
	fmt.Println("===============================")
}

func check(db *Database) {
	fmt.Println("========= Tdb warnings ========")
	warnings := tdb.Check(db)
	if len(warnings) > 0 {
		for _, warning := range warnings {
			fmt.Println(warning)
		}
	} else {
		fmt.Println("No warnings")
	}
}

type Database struct {
	tdb.MetaData `tdb:"MetaData"`
	Customers    []Customer `tdb:"Customer"`
	Invoices     []Invoice  `tdb:"Invoice"`
	LineItems    []LineItem `tdb:"LineItem"`
}

type Customer struct {
	Cid     int     `tdb:"CID"`
	Company string  `tdb:"Company"`
	Address *string `tdb:"Address"` // Pointer since nullable
	Contact string  `tdb:"Contact"`
	Email   string  `tdb:"Email"`
	Icon    []byte  `tdb:"Icon"` // PNG; null indicated by nil
}

type Invoice struct {
	Inum   int       `tdb:"INUM"`
	Cid    int       `tdb:"CID"`
	Raised time.Time `tdb:"Raised_Date"`
	Due    time.Time `tdb:"Due_Date"`
	Paid   bool      `tdb:"Paid"`
	Desc   *string   `tdb:"Desc"`
}

type LineItem struct {
	Liid      int       `tdb:"LIID"`
	Inum      int       `tdb:"INUM"`
	Delivered time.Time `tdb:"Delivery_Date"`
	UnitPrice float64   `tdb:"Unit_Price"`
	Quantity  int       `tdb:"Quantity"`
	Desc      *string   `tdb:"Desc"`
}

func makeDb(t *testing.T) Database {
	db := makeDataTables(t)
	db.MetaData = tdb.NewMetaData()
	db.MetaData.Add(makeCustomerMeta())
	db.MetaData.Add(makeInvoiceMeta())
	db.MetaData.Add(makeLineItemMeta())
	return db
}

func makeDataTables(t *testing.T) Database {
	icon, err := hex.DecodeString("89504E470D0A1A0A0000000D494844520000000C0000000C080600000056755CE7000000097048597300000EC400000EC401952B0E1B0000002849444154289163646068F8CF4002602245317D34B0600A3530621183FB7310FA81640D8C832FE2002C7F051786CBFA670000000049454E44AE426082")
	if err != nil {
		t.Error(err)
	}
	custAddress := "123 Somewhere"
	invDesc := "COD"
	lineDesc := []string{"Bales of <hay>", "Straps & Things",
		"Washers (1\")"}
	db := Database{
		Customers: []Customer{
			{50, "Best People", &custAddress, "John Doe", "j@doe.com", nil},
			{19, "Supersuppliers", nil, "Jane Doe", "jane@super.com", nil},
		},
		Invoices: []Invoice{
			{152, 50, time.Date(2022,
				time.January, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2022, time.February, 17, 0, 0, 0, 0, time.UTC),
				false, &invDesc},
			{153, 19,
				time.Date(2022, time.January, 19, 0, 0, 0, 0, time.UTC),
				time.Date(2022, time.February, 19, 0, 0, 0, 0, time.UTC),
				true, nil},
		},
		LineItems: []LineItem{
			{1839, 152,
				time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
				29.99, 2, &lineDesc[0]},
			{1840, 152,
				time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
				5.98, 3, &lineDesc[1]},
			{1620, 153,
				time.Date(2022, time.January, 19, 0, 0, 0, 0, time.UTC),
				11.5, 1, &lineDesc[2]},
		}}
	db.Customers[0].Icon = icon
	return db
}

func makeCustomerMeta() tdb.MetaTable {
	customer := tdb.NewMetaTable("Customers")
	cid := tdb.IntField("CID")
	cid.SetUnique()
	_ = cid.SetMin(1)
	customer.Add(cid)
	customer.Add(tdb.StrField("Company"))
	address := tdb.StrField("Address")
	address.SetNullable()
	customer.Add(address)
	customer.Add(tdb.StrField("Contact"))
	customer.Add(tdb.StrField("Email"))
	iconField := tdb.BytesField("Icon")
	iconField.SetNullable()
	customer.Add(iconField)
	return customer
}

func makeInvoiceMeta() tdb.MetaTable {
	invoice := tdb.NewMetaTable("Invoices")
	inum := tdb.IntField("INUM")
	inum.SetUnique()
	_ = inum.SetMin(100)
	invoice.Add(inum)
	cid := tdb.IntField("CID")
	_ = cid.SetRef("Customers.CID")
	invoice.Add(cid)
	invoice.Add(tdb.DateField("Raised_Date"))
	invoice.Add(tdb.DateField("Due_Date"))
	invoice.Add(tdb.BoolField("Paid"))
	desc := tdb.StrField("Description")
	desc.SetNullable()
	invoice.Add(desc)
	return invoice
}

func makeLineItemMeta() tdb.MetaTable {
	lineItem := tdb.NewMetaTable("LineItems")
	liid := tdb.IntField("LIID")
	liid.SetUnique()
	_ = liid.SetMin(1)
	lineItem.Add(liid)
	inum := tdb.IntField("INUM")
	_ = inum.SetRef("Invoices.INUM")
	lineItem.Add(inum)
	lineItem.Add(tdb.DateField("Delivery_Date"))
	price := tdb.RealField("Unit_Price")
	_ = price.SetMin(0.0)
	lineItem.Add(price)
	quantity := tdb.IntField("Quantity")
	_ = quantity.SetMin(0)
	lineItem.Add(quantity)
	desc := tdb.StrField("Description")
	desc.SetNullable()
	lineItem.Add(desc)
	return lineItem
}
