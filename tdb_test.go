package tdb_test

import (
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

func Test001(t *testing.T) {
	icon, err := hex.DecodeString("89504E470D0A1A0A0000000D494844520000000C0000000C080600000056755CE7000000097048597300000EC400000EC401952B0E1B0000002849444154289163646068F8CF4002602245317D34B0600A3530621183FB7310FA81640D8C832FE2002C7F051786CBFA670000000049454E44AE426082")
	if err != nil {
		t.Error(err)
	}
	TdbDatabase.Customers[0].Icon = icon
	fmt.Println("======= Tdb example data ======")
	fmt.Println(TdbDatabase)
	raw, err := tdb.Marshal(TdbDatabase)
	if err != nil {
		t.Error(err)
	}
	if string(raw) != TdbText {
		t.Error("raw != text")
	}
	fmt.Println("======= Tdb example text ======")
	raw = []byte(TdbText)
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
	text := string(raw)
	fmt.Println("======= Tdb result text =======")
	fmt.Println(text)
	fmt.Println("===============================")
}

var TdbDatabase = Database{
	Customers: []Customer{
		{50, "Best People", "123 Somewhere", "John Doe", "j@doe.com", nil},
		{19, "Supersuppliers", "", "Jane Doe", "jane@super.com", nil},
	},
	Invoices: []Invoice{
		{152, 50, time.Date(2022, time.January, 17, 0, 0, 0, 0, time.UTC),
			time.Date(2022, time.February, 17, 0, 0, 0, 0, time.UTC),
			false, "COD"},
		{153, 19, time.Date(2022, time.January, 19, 0, 0, 0, 0, time.UTC),
			time.Date(2022, time.February, 19, 0, 0, 0, 0, time.UTC), true,
			""},
	},
	LineItems: []LineItem{
		{1839, 152, time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
			29.99, 2, "Bales of <hay>"},
		{1840, 152, time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
			5.98, 3, "Straps & Things"},
		{1620, 153, time.Date(2022, time.January, 19, 0, 0, 0, 0, time.UTC),
			11.5, 1, "Washers (1\")"},
	},
}

type Database struct {
	Customers []Customer `tdb:"Customer"`
	Invoices  []Invoice  `tdb:"Invoice"`
	LineItems []LineItem `tdb:"LineItem"`
}

type Customer struct {
	Cid     int    `tdb:"CID"`
	Company string `tdb:"Company"`
	Address string `tdb:"Address"`
	Contact string `tdb:"Contact"`
	Email   string `tdb:"Email"`
	Icon    []byte `tdb:"Icon"` // PNG
}

type Invoice struct {
	Inum   int       `tdb:"INUM"`
	Cid    int       `tdb:"CID"`
	Raised time.Time `tdb:"Raised_Date"`
	Due    time.Time `tdb:"Due_Date"`
	Paid   bool      `tdb:"Paid"`
	Desc   string    `tdb:"Desc"`
}

type LineItem struct {
	Liid      int       `tdb:"LIID"`
	Inum      int       `tdb:"INUM"`
	Delivered time.Time `tdb:"Delivered_Date"`
	UnitPrice float64   `tdb:"Unit_Price"`
	Quantity  int       `tdb:"Quantity"`
	Desc      string    `tdb:"Desc"`
}

const TdbText = `[Customers CID int Company str Address str Contact str Email str Icon bytes
%
50 <Best People> <123 Somewhere> <John Doe> <j@doe.com> (89504E470D0A1A0A0000000D494844520000000C0000000C080600000056755CE7000000097048597300000EC400000EC401952B0E1B0000002849444154289163646068F8CF4002602245317D34B0600A3530621183FB7310FA81640D8C832FE2002C7F051786CBFA670000000049454E44AE426082)
19 <Supersuppliers> # <Jane Doe> <jane@super.com> #
]
[Invoices INUM int CID int Raised_Date date Due_Date date Paid bool Description str
%
152 50 2022-01-17 2022-02-17 no <COD> 
153 19 2022-01-19 2022-02-19 yes # 
]
[LineItems LIID int INUM int Delivery_Date date Unit_Price real Quantity int Description str
%
1839 152 2022-01-16 29.99 2 <Bales of &lt;hay&gt;> 
1840 152 2022-01-16 5.98 3 <Straps &amp; Things> 
1620 153 2022-01-19 11.5 1 <Washers (1")> 
]
`
