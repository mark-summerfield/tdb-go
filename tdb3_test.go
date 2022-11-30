package tdb_test

import (
	_ "embed"
	"fmt"
	tdb "github.com/mark-summerfield/tdb-go"
	"os"
	"testing"
	"time"
)

//go:embed eg/classic.tdb
var Classic string

func ExampleMarshal() {
	db := classicDatabase{
		Employees: []Employee{
			{7844, "TURNER", "SALESMAN", nil,
				date(1981, time.September, 8), 1500.0, nil, 30},
			{7876, "ADAMS", "CLERK", nil, date(1983, time.January, 12),
				1100.0, nil, 20},
			{7839, "KING", "PRESIDENT", nil,
				date(1981, time.November, 17), 5000.0, nil, 10},
			{7902, "FORD", "ANALYST", nil, date(1981, time.December, 3),
				3000.0, nil, 20},
		},
		Departments: []Department{
			{10, "ACCOUNTING", "NEW YORK"},
			{20, "RESEARCH", "DALLAS"},
			{30, "SALES", "CHICAGO"},
		},
	}
	c0 := 0.0
	db.Employees[0].Commission = &c0
	m0 := 7698
	db.Employees[0].ManagerID = &m0
	m1 := 7788
	db.Employees[1].ManagerID = &m1
	m3 := 7566
	db.Employees[3].ManagerID = &m3
	raw, err := tdb.Marshal(db)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))
	// Output:
	// [emp empno int ename str job str mgr int? hiredate date sal real comm real? deptno int
	// %
	// 7844 <TURNER> <SALESMAN> 7698 1981-09-08 1500 0 30
	// 7876 <ADAMS> <CLERK> 7788 1983-01-12 1100 ? 20
	// 7839 <KING> <PRESIDENT> ? 1981-11-17 5000 ? 10
	// 7902 <FORD> <ANALYST> 7566 1981-12-03 3000 ? 20
	// ]
	// [dept deptno int dname str loc str
	// %
	// 10 <ACCOUNTING> <NEW YORK>
	// 20 <RESEARCH> <DALLAS>
	// 30 <SALES> <CHICAGO>
	// ]
}

func ExampleUnmarshal() {
	tdbText := `[emp empno int ename str job str mgr int? hiredate date sal real comm real? deptno int
%
7844 <TURNER> <SALESMAN> 7698 1981-09-08 1500 0 30
7876 <ADAMS> <CLERK> 7788 1983-01-12 1100 ? 20
7839 <KING> <PRESIDENT> ? 1981-11-17 5000 ? 10
7902 <FORD> <ANALYST> 7566 1981-12-03 3000 ? 20
]
[dept deptno int dname str loc str
%
10 <ACCOUNTING> <NEW YORK>
20 <RESEARCH> <DALLAS>
30 <SALES> <CHICAGO>
]`
	db := classicDatabase{}
	if err := tdb.Unmarshal([]byte(tdbText), &db); err != nil {
		panic(err)
	}
	fmt.Printf("%d Employees\n", len(db.Employees))
	fmt.Printf("%d Departments\n", len(db.Departments))
	president := db.Employees[2]
	managerID := president.ManagerID
	commission := president.Commission
	fmt.Printf("%d %q %q %d %s %g %g %d\n", president.EID, president.Name,
		president.Job, *managerID,
		president.HireDate.Format(tdb.DateFormat),
		president.Salary, *commission, president.DeptID)
	research := db.Departments[1]
	fmt.Printf("%d %q %q\n", research.DID, research.Name, research.Location)
	// Output:
	// 4 Employees
	// 3 Departments
	// 7839 "KING" "PRESIDENT" -1 1981-11-17 5000 0 10
	// 20 "RESEARCH" "DALLAS"
}

func TestClassic(t *testing.T) {
	db := makeClassic(t)
	raw, err := tdb.MarshalDecimals(db, 1)
	if err != nil {
		t.Error(err)
	}
	if string(raw) != Classic {
		fmt.Println("======= Tdb Marshal Database ======")
		fmt.Print(string(raw))
		_ = os.WriteFile("/tmp/1", raw, 0666)
		_ = os.WriteFile("/tmp/2", []byte(Classic), 0666)
		fmt.Println("wrote /tmp/[12] (actual/expected)")
		fmt.Println("===================================")
		t.Error("Database: raw != text")
	}
	var database classicDatabase
	err = tdb.Unmarshal(raw, &database)
	if err != nil {
		t.Error(err)
	}
	raw2, err := tdb.MarshalDecimals(database, 1)
	if err != nil {
		t.Error(err)
	}
	if string(raw2) != Classic {
		fmt.Println("====== Tdb Unmarshal Database =====")
		fmt.Print(string(raw2))
		_ = os.WriteFile("/tmp/3", raw2, 0666)
		_ = os.WriteFile("/tmp/4", []byte(Classic), 0666)
		fmt.Println("wrote /tmp/[34] (actual/expected)")
		fmt.Println("===================================")
		t.Error("Database: raw2 != text")
	}
}

type classicDatabase struct {
	Employees   []Employee   `tdb:"emp"`
	Departments []Department `tdb:"dept"`
}

type Employee struct {
	EID        int       `tdb:"empno"`
	Name       string    `tdb:"ename"`
	Job        string    `tdb:"job"`
	ManagerID  *int      `tdb:"mgr"`
	HireDate   time.Time `tdb:"hiredate:date"`
	Salary     float64   `tdb:"sal"`
	Commission *float64  `tdb:"comm"`
	DeptID     int       `tdb:"deptno"`
}

type Department struct {
	DID      int    `tdb:"deptno"`
	Name     string `tdb:"dname"`
	Location string `tdb:"loc"`
}

func makeClassic(t *testing.T) classicDatabase {
	db := classicDatabase{
		Employees: []Employee{
			{7369, "SMITH", "CLERK", nil, date(1980, time.December, 17),
				800.0, nil, 20},
			{7499, "ALLEN", "SALESMAN", nil, date(1981, time.February, 20),
				1600.0, nil, 30},
			{7521, "WARD", "SALESMAN", nil, date(1981, time.February, 22),
				1250.0, nil, 30},
			{7566, "JONES", "MANAGER", nil, date(1981, time.April, 2),
				2975.0, nil, 20},
			{7654, "MARTIN", "SALESMAN", nil, date(1981,
				time.September, 28), 1250.0, nil, 30},
			{7698, "BLAKE", "MANAGER", nil, date(1981, time.May, 1),
				2850.0, nil, 30},
			{7782, "CLARK", "MANAGER", nil, date(1981, time.June, 9),
				2450.0, nil, 10},
			{7788, "SCOTT", "ANALYST", nil, date(1982, time.December, 9),
				3000.0, nil, 20},
			{7839, "KING", "PRESIDENT", nil,
				date(1981, time.November, 17), 5000.0, nil, 10},
			{7844, "TURNER", "SALESMAN", nil,
				date(1981, time.September, 8), 1500.0, nil, 30},
			{7876, "ADAMS", "CLERK", nil, date(1983, time.January, 12),
				1100.0, nil, 20},
			{7900, "JAMES", "CLERK", nil, date(1981, time.December, 3),
				950.0, nil, 30},
			{7902, "FORD", "ANALYST", nil, date(1981, time.December, 3),
				3000.0, nil, 20},
			{7934, "MILLER", "CLERK", nil, date(1982, time.January, 23),
				1300.0, nil, 10},
		},
		Departments: []Department{
			{10, "ACCOUNTING", "NEW YORK"},
			{20, "RESEARCH", "DALLAS"},
			{30, "SALES", "CHICAGO"},
			{40, "OPERATIONS", "BOSTON"},
		},
	}
	m0 := 7902
	db.Employees[0].ManagerID = &m0
	m1 := 7698
	c1 := 300.0
	db.Employees[1].ManagerID = &m1
	db.Employees[1].Commission = &c1
	m2 := 7698
	c2 := 500.0
	db.Employees[2].ManagerID = &m2
	db.Employees[2].Commission = &c2
	m3 := 7839
	db.Employees[3].ManagerID = &m3
	m4 := 7698
	c4 := 1400.0
	db.Employees[4].ManagerID = &m4
	db.Employees[4].Commission = &c4
	m5 := 7839
	db.Employees[5].ManagerID = &m5
	m6 := 7839
	db.Employees[6].ManagerID = &m6
	m7 := 7566
	db.Employees[7].ManagerID = &m7
	m9 := 7698
	c9 := 0.0
	db.Employees[9].ManagerID = &m9
	db.Employees[9].Commission = &c9
	m10 := 7788
	db.Employees[10].ManagerID = &m10
	m11 := 7698
	db.Employees[11].ManagerID = &m11
	m12 := 7566
	db.Employees[12].ManagerID = &m12
	m13 := 7782
	db.Employees[13].ManagerID = &m13
	return db
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
