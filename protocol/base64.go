package protocol

import (
	"bytes"
	"encoding/base64"
	"reflect"
)

// URLEncodedBase64 represents a byte slice holding URL-encoded base64 data.
// When fields of this type are unmarshaled from JSON, the data is base64
// decoded into a byte slice.
type URLEncodedBase64 []byte

// UnmarshalJSON base64 decodes a URL-encoded value, storing the result in the
// provided byte slice.
func (e *URLEncodedBase64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	// TODO: Investigate this line. It is commented as trimming the leading spaces but appears to trim the leading and trailing double quotes instead.
	// Trim the leading spaces
	data = bytes.Trim(data, "\"")

	// Trim the trailing equal characters.
	data = bytes.TrimRight(data, "=")

	out := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	n, err := base64.RawURLEncoding.Decode(out, data)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(e).Elem()
	v.SetBytes(out[:n])
	return nil
}

// MarshalJSON base64 encodes a non URL-encoded value, storing the result in the
// provided byte slice.
func (e URLEncodedBase64) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + base64.RawURLEncoding.EncodeToString(e) + `"`), nil
}
