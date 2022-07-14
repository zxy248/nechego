package fishing

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

func (f *Fish) Scan(src any) error {
	switch v := src.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), f); err != nil {
			return err
		}
	default:
		return errors.New("unknown type")
	}
	return nil
}

func (f Fish) Value() (driver.Value, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
