package sqlite3orm

import (
    "fmt"
    "reflect"
    "bytes"
    "errors"
)

func (w DBWrapper) Insert(instance interface{}) error {
	
	var buffer bytes.Buffer
	var values bytes.Buffer
	var params []interface{}

	v := reflect.ValueOf(instance)

	if v.Kind() != reflect.Ptr {
        return errors.New("expecting a pointer")
    }

	el := v.Elem()

	t := el.Type()

	buffer.WriteString(fmt.Sprintf("INSERT INTO %s (", t.Name()))

	for i := 0; i < t.NumField(); i++ {

		tField := t.Field(i)
		vField := el.Field(i)

		if i > 0 {
			buffer.WriteString(", ")
			values.WriteString(", ")
		}

		buffer.WriteString(tField.Name)
		values.WriteString("?")

		params = append(params, vField.Interface())
	}

	buffer.WriteString(") VALUES (")
	buffer.WriteString(values.String())
	buffer.WriteString(")")

	if w.isDebug {
		fmt.Println(buffer.String())		
	}

    _, err := w.DB.Exec(buffer.String(), params...)
    
    if err != nil {
        return err
    }

	return nil
}