package types

import (
	"encoding/json"
	"fmt"
)

// FlexString handles both string and number types from JSON
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}

	// If that fails, try as number
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexString(fmt.Sprintf("%f", num))
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexString", data)
}
