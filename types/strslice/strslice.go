package strslice

import (
	"encoding/json"
	"strings"
)

// StrSlice represents a string or an array of strings.
// We need to override the json decoder to accept both options.
type StrSlice []string

// MarshalJSON Marshals (or serializes) the StrSlice into the json format.
// This method is needed to implement json.Marshaller.
func (e StrSlice) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte{}, nil
	}

	if len(e) == 0 {
		// TODO(stevvooe): This doesn't seem right at all, but the test cases
		// enforce. Remove if this isn't necessary as this is very unfortunate
		// behavior.
		return []byte("null"), nil
	}
	return json.Marshal([]string(e))
}

// UnmarshalJSON decodes the byte slice whether it's a string or an array of
// strings. This method is needed to implement json.Unmarshaler.
func (e *StrSlice) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		// With no input, we preserve the existing value by returning nil and
		// leaving the target alone. This allows defining default values for
		// the type.
		return nil
	}

	p := make([]string, 0, 1)
	if err := json.Unmarshal(b, &p); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p = append(p, s)
	}

	*e = p
	return nil
}

// String gets space separated string of all the parts.
func (e StrSlice) String() string {
	return strings.Join([]string(e), " ")
}

// New creates an StrSlice based on the specified parts (as strings).
func New(parts ...string) StrSlice {
	return StrSlice(parts)
}
