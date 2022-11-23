package tdb_test

import (
	_ "embed"
	"fmt"
	"github.com/mark-summerfield/gong"
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
			{7844, "TURNER", "SALESMAN", 7698,
				date(1981, time.September, 8), 1500.0, 0.0, 30},
			{7876, "ADAMS", "CLERK", 7788, date(1983, time.January, 12),
				1100.0, tdb.RealSentinal, 20},
			{7839, "KING", "PRESIDENT", tdb.IntSentinal,
				date(1981, time.November, 17), 5000.0, tdb.RealSentinal, 10},
			{7902, "FORD", "ANALYST", 7566, date(1981, time.December, 3),
				3000.0, tdb.RealSentinal, 20},
		},
		Departments: []Department{
			{10, "ACCOUNTING", "NEW YORK"},
			{20, "RESEARCH", "DALLAS"},
			{30, "SALES", "CHICAGO"},
		},
	}
	raw, err := tdb.Marshal(db)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))
	// Output:
	// [emp empno int ename str job str mgr int hiredate date sal real comm real deptno int
	// %
	// 7844 <TURNER> <SALESMAN> 7698 1981-09-08 1500 0 30
	// 7876 <ADAMS> <CLERK> 7788 1983-01-12 1100 ! 20
	// 7839 <KING> <PRESIDENT> ! 1981-11-17 5000 ! 10
	// 7902 <FORD> <ANALYST> 7566 1981-12-03 3000 ! 20
	// ]
	// [dept deptno int dname str loc str
	// %
	// 10 <ACCOUNTING> <NEW YORK>
	// 20 <RESEARCH> <DALLAS>
	// 30 <SALES> <CHICAGO>
	// ]
}

func ExampleUnmarshal() {
	tdbText := `[emp empno int ename str job str mgr int hiredate date sal real comm real deptno int
%
7844 <TURNER> <SALESMAN> 7698 1981-09-08 1500 0 30
7876 <ADAMS> <CLERK> 7788 1983-01-12 1100 ! 20
7839 <KING> <PRESIDENT> ! 1981-11-17 5000 ! 10
7902 <FORD> <ANALYST> 7566 1981-12-03 3000 ! 20
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
	if managerID == tdb.IntSentinal {
		managerID = -1
	}
	commission := president.Commission
	if gong.IsRealClose(commission, tdb.RealSentinal) {
		commission = 0.0
	}
	fmt.Printf("%d %q %q %d %s %g %g %d\n", president.EID, president.Name,
		president.Job, managerID, president.HireDate.Format(tdb.DateFormat),
		president.Salary, commission, president.DeptID)
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
	ManagerID  int       `tdb:"mgr"`
	HireDate   time.Time `tdb:"hiredate:date"`
	Salary     float64   `tdb:"sal"`
	Commission float64   `tdb:"comm"`
	DeptID     int       `tdb:"deptno"`
}

type Department struct {
	DID      int    `tdb:"deptno"`
	Name     string `tdb:"dname"`
	Location string `tdb:"loc"`
}

func makeClassic(t *testing.T) classicDatabase {
	return classicDatabase{
		Employees: []Employee{
			{7369, "SMITH", "CLERK", 7902, date(1980, time.December, 17),
				800.0, tdb.RealSentinal, 20},
			{7499, "ALLEN", "SALESMAN", 7698, date(1981, time.February, 20),
				1600.0, 300.0, 30},
			{7521, "WARD", "SALESMAN", 7698, date(1981, time.February, 22),
				1250.0, 500.0, 30},
			{7566, "JONES", "MANAGER", 7839, date(1981, time.April, 2),
				2975.0, tdb.RealSentinal, 20},
			{7654, "MARTIN", "SALESMAN", 7698, date(1981,
				time.September, 28), 1250.0, 1400.0, 30},
			{7698, "BLAKE", "MANAGER", 7839, date(1981, time.May, 1),
				2850.0, tdb.RealSentinal, 30},
			{7782, "CLARK", "MANAGER", 7839, date(1981, time.June, 9),
				2450.0, tdb.RealSentinal, 10},
			{7788, "SCOTT", "ANALYST", 7566, date(1982, time.December, 9),
				3000.0, tdb.RealSentinal, 20},
			{7839, "KING", "PRESIDENT", tdb.IntSentinal,
				date(1981, time.November, 17), 5000.0, tdb.RealSentinal,
				10},
			{7844, "TURNER", "SALESMAN", 7698,
				date(1981, time.September, 8), 1500.0, 0.0, 30},
			{7876, "ADAMS", "CLERK", 7788, date(1983, time.January, 12),
				1100.0, tdb.RealSentinal, 20},
			{7900, "JAMES", "CLERK", 7698, date(1981, time.December, 3),
				950.0, tdb.RealSentinal, 30},
			{7902, "FORD", "ANALYST", 7566, date(1981, time.December, 3),
				3000.0, tdb.RealSentinal, 20},
			{7934, "MILLER", "CLERK", 7782, date(1982, time.January, 23),
				1300.0, tdb.RealSentinal, 10},
		},
		Departments: []Department{
			{10, "ACCOUNTING", "NEW YORK"},
			{20, "RESEARCH", "DALLAS"},
			{30, "SALES", "CHICAGO"},
			{40, "OPERATIONS", "BOSTON"},
		},
	}
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
