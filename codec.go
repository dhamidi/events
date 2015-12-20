package events

import "encoding/json"

type Codec struct {
	Decode func(src []byte, dest interface{}) error
	Encode func(src interface{}) ([]byte, error)
}

var (
	JSONCodec = &Codec{
		Decode: json.Unmarshal,
		Encode: json.Marshal,
	}
)
