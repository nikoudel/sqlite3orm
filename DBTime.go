package sqlite3orm

import (
    "fmt"
    "errors"
    "time"
    "database/sql/driver"
)

type DBTime struct {
	Time time.Time
}

func (t *DBTime) Scan(src interface{}) error {
	
	if src != nil {

		switch v := src.(type) {

		    case time.Time:
		        t.Time = src.(time.Time)

		    case []uint8:

		    	var err error

		        t.Time, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(src.([]uint8)))

		        if err != nil {
		        	return err
		        }

	        default:
		        return errors.New(fmt.Sprintf("failed parsing a time.Time; unexpected type %T", v))
	    }
	}

	return nil
}

func (t DBTime) String() string {
    return fmt.Sprintf("%v", t.Time)
}

func (t DBTime) Value() (driver.Value, error) {
	return t.Time.String(), nil
}
