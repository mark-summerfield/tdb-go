// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"
	"github.com/mark-summerfield/tdb"
	"time"
)

func main() {
	fmt.Println("======= Tdb example text ======")
	raw := []byte(Text)
	fmt.Println(string(raw))
	var database Database
	err := tdb.Unmarshal(raw, &database)
	if err != nil {
		panic(err)
	}
	fmt.Println("========== Database ===========")
	fmt.Println(database)
	raw, err = tdb.Marshal(database)
	if err != nil {
		panic(err)
	}
	text := string(raw)
	fmt.Println("======= Tdb result text =======")
	fmt.Println(text)
	fmt.Println("===============================")
}

const Text = `[Customers CID int Company str Address str Contact str Email str
%
50 <Best People> <123 Somewhere> <John Doe> <j@doe.com> 
19 <Supersuppliers> # <Jane Doe> <jane@super.com> 
]
[Invoices INUM int CID int Raised_Date date Due_Date date Paid bool Description str
%
152 50 2022-01-17 2022-02-17 no <COD> 
153 19 2022-01-19 2022-02-19 yes # 
]
[Items IID int INUM int Delivery_Date date Unit_Price real Quantity int Description str
%
1839 152 2022-01-16 29.99 2 <Bales of hay> 
1840 152 2022-01-16 5.98 3 <Straps> 
1620 153 2022-01-19 11.5 1 <Washers (1-in)> 
]
`

type Database struct {
	Customers []Customer
	Invoices  []Invoice
	Items     []Item
}

type Customer struct {
	Cid     int    `tdb:"CID"`
	Company string `tdb:"Company"`
	Address string `tdb:"Address"`
	Contact string `tdb:"Contact"`
	Email   string `tdb:"Email"`
}

type Invoice struct {
	Inum   int       `tdb:"INUM"`
	Cid    int       `tdb:"CID"`
	Raised time.Time `tdb:"Raised_Date"`
	Due    time.Time `tdb:"Due_Date"`
	Paid   bool      `tdb:"Paid"`
	Desc   string    `tdb:"Desc"`
}

type Item struct {
	Iid       int       `tdb:"IID"`
	Inum      int       `tdb:"INUM"`
	Delivered time.Time `tdb:"Delivered_Date"`
	UnitPrice float64   `tdb:"Unit_Price"`
	Quantity  int       `tdb:"Quantity"`
	Desc      string    `tdb:"Desc"`
}
