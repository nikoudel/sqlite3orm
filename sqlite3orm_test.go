package sqlite3orm

import (
    "testing"
    "os"
    "time"
    "database/sql"
)

type TestStruct struct {
    String string `index:"ix_string_time"`
    Float float64
    Uint uint64
    Bool bool
    Time DBTime `index:"ix_string_time"`
}

func TestCreateInsertSelect(test *testing.T) {

	os.Remove("./test.db")

    db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		test.Errorf("failed opening database: %v\n", err)
		return
	}

	w := DBWrapper{SqlDB: db, isDebug: true}

    defer w.SqlDB.Close()

	instance := TestStruct{String: "abc", Float: 1.23, Uint: 123, Bool: true, Time: DBTime{Time: time.Now()}}

	err = w.CreateTable(instance)

	if err != nil {
		test.Errorf("failed creating table: %v\n", err)
		return
	}

	err = w.Insert(&instance)

	if err != nil {
		test.Errorf("insert failed: %v\n", err)
		return
	}

	dbItem := TestStruct{}

	err = w.SelectFirst(&dbItem, "String = 'abc'")

	if err != nil {
		test.Errorf("select failed: %v\n", err)
		return
	}

	if dbItem.String != instance.String || dbItem.Float != instance.Float || dbItem.Uint != instance.Uint || dbItem.Bool != instance.Bool || dbItem.Time.Time != instance.Time.Time {

		test.Errorf("expecting: %v\n", instance)
		test.Errorf("actual: %v\n", dbItem)
	}

	instance = TestStruct{String: "def", Float: 3.45, Uint: 345, Bool: false, Time: DBTime{Time: time.Now()}}

	err = w.Insert(&instance)

	if err != nil {
		test.Errorf("insert failed: %v\n", err)
		return
	}

	var items []TestStruct

	err = w.Select(&items, "String = 'abc' OR String = 'def'")

	if err != nil {
		test.Errorf("select failed: %v\n", err)
		return
	}

	if len(items) != 2 {
		test.Errorf("expecting: %d items\n", 2)
		test.Errorf("actual: %d items\n", len(items))
		return
	}

	if items[0].String != "abc" {
		test.Errorf("expecting: %s\n", "abc")
		test.Errorf("actual: %s\n", items[0].String)
	}

	if items[1].String != "def" {
		test.Errorf("expecting: %s\n", "def")
		test.Errorf("actual: %s\n", items[1].String)
	}
}