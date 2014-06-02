package sqlite3orm

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

func (w DBWrapper) Select(instance interface{}, where string) error {

	slice, err := getElem(instance, true)

	if err != nil {
		return err
	}

	tmp := reflect.Zero(reflect.SliceOf(slice.Type().Elem()))

	err = w.selectRows(&tmp, where, false)

	if err != nil {
		return err
	}

	slice.Set(tmp)

	return nil
}

func (w DBWrapper) SelectFirst(instance interface{}, where string) error {

	el, err := getElem(instance, false)

	if err != nil {
		return err
	}

	slice := reflect.Zero(reflect.SliceOf(el.Type()))

	if err := w.selectRows(&slice, where, true); err != nil {
		return err
	}

	if slice.Len() > 0 {
		el.Set(slice.Index(0))
	} else {
		return errors.New("empty result")
	}

	return nil
}

func (w DBWrapper) selectRows(slice *reflect.Value, where string, getOnlyFirst bool) error {

	sliceItemType := slice.Type().Elem()

	var buffer bytes.Buffer

	buffer.WriteString("SELECT ")

	for i := 0; i < sliceItemType.NumField(); i++ {

		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(sliceItemType.Field(i).Name)
	}

	buffer.WriteString(fmt.Sprintf(" FROM %s", sliceItemType.Name()))

	if where != "" {
		buffer.WriteString(" WHERE ")
		buffer.WriteString(where)
	}

	fmt.Println(buffer.String())

	rows, err := w.SqlDB.Query(buffer.String())

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {

		if err = addItem(rows, slice); err != nil {
			return err
		}

		if getOnlyFirst {
			return nil
		}
	}

	return nil
}

func getElem(instance interface{}, isSliceExpected bool) (reflect.Value, error) {

	p := reflect.ValueOf(instance)

	if p.Kind() != reflect.Ptr {
		return p, errors.New("expecting a pointer")
	}

	elem := p.Elem()

	if isSliceExpected && elem.Kind() != reflect.Slice {
		return elem, errors.New("expecting a pointer to a slice")
	}

	if !isSliceExpected && elem.Kind() == reflect.Slice {
		return elem, errors.New("not expecting a slice")
	}

	return elem, nil
}

func addItem(rows *sql.Rows, slice *reflect.Value) error {

	sliceItemType := slice.Type().Elem()

	pItem := reflect.New(sliceItemType)

	item := pItem.Elem()

	dest := make([]interface{}, item.NumField())

	for i := 0; i < item.NumField(); i++ {

		dest[i] = item.Field(i).Addr().Interface()
	}

	if err := rows.Scan(dest...); err != nil {
		return err
	}

	*slice = reflect.Append(*slice, item)

	return nil
}
