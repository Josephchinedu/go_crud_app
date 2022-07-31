package models

import (
	"encoding/json"
	"strconv"
)

type StringInt int

// UnmarshalJSON create a custom unmarshal for the StringInt
/// this helps us check the type of our value before unmarshalling it

func (st *StringInt) UnmarshalJSON(b []byte) error {

	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int64(v))
	case string:

		i, err := strconv.Atoi(v)
		if err != nil {

			return err

		}
		*st = StringInt(i)

	}
	return nil
}

// User schemaof the user table
type User struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Age      StringInt `json:"age"`
}
