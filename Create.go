package sqlite3orm

import (
    "fmt"
    "reflect"
    "errors"
    "bytes"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type DBWrapper struct {
	SqlDB *sql.DB
	isDebug bool
}

func (w DBWrapper) CreateTable(instance interface{}) error {

	var buffer bytes.Buffer

	indexMap := make(map[string][]string)

	v := reflect.ValueOf(instance)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
        t = v.Elem().Type()
    }

    if t.Kind() == reflect.Slice {
    	return errors.New("not expecting a slice")
    }

	buffer.WriteString(fmt.Sprintf("CREATE TABLE %s (", t.Name()))		

	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)

		sqlType, err := mapSqliteType(field.Type.Name())

		if err != nil {
			return err
		}

		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(field.Name)
		buffer.WriteString(" ")
		buffer.WriteString(sqlType)

		index := field.Tag.Get("index")

		if index != "" {

			_, ok := indexMap[index]

			if ok {
				indexMap[index] = append(indexMap[index], field.Name)
			} else {
				indexMap[index] = []string{field.Name}
			}
		}
	}

	buffer.WriteString(")")	

	if w.isDebug {
		fmt.Println(buffer.String())		
	}

    _, err := w.SqlDB.Exec(buffer.String())
    
    if err != nil {
        return err
    }

	return createIndex(indexMap, t.Name(), w)
}

func mapSqliteType(instance interface{}) (string, error) {
	
	switch t := instance.(type) {

	case string, DBTime:
		return "TEXT", nil

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool:
		return "INTEGER", nil

	case float32, float64:
		return "REAL", nil
		
	default:
		return "", errors.New(fmt.Sprintf("unknown type: %T", t))
	}
}

func createIndex(indexMap map[string][]string, typeName string, w DBWrapper) error {

	var buffer bytes.Buffer

	for index, fields := range indexMap {

		buffer.WriteString(fmt.Sprintf("CREATE INDEX %s ON %s (", index, typeName))

    	for i, field := range fields {

    		if i > 0 {
    			buffer.WriteString(", ")
    		}

    		buffer.WriteString(field)
    	}

    	buffer.WriteString(")")

		if w.isDebug {
			fmt.Println(buffer.String())		
		}

    	_, err := w.SqlDB.Exec(buffer.String())
    
	    if err != nil {
	        return err
	    }

		buffer.Reset()
	}

	return nil
}