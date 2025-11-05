package types

import (
	"encoding/json"
	"fmt"
)

type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}

	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexString(fmt.Sprintf("%f", num))
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexString", data)
}
