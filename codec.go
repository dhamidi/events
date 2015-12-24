package events

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Codec struct {
	Decode func(src []byte, dest interface{}) error
	Encode func(src interface{}) ([]byte, error)
}

var (
	JSONCodec = &Codec{
		Decode: json.Unmarshal,
		Encode: json.Marshal,
	}

	FormCodec = &Codec{
		Decode: decodeForm,
		Encode: encodeForm,
	}
)

type FormUnmarshaler interface {
	UnmarshalForm(values url.Values) error
}

type FormMarshaler interface {
	MarshalForm() (url.Values, error)
}

func decodeForm(content []byte, dest interface{}) error {
	values, err := url.ParseQuery(string(content))
	if err != nil {
		return err
	}

	unmarshaler, ok := dest.(FormUnmarshaler)
	if ok {
		return unmarshaler.UnmarshalForm(values)
	}

	unmarshalOne := func(from string, dest interface{}) error {
		unmarshaler, ok := dest.(encoding.TextUnmarshaler)
		if ok {
			return unmarshaler.UnmarshalText([]byte(from))
		}
		destString, ok := dest.(*string)
		if ok {
			*destString = from
			return nil
		}

		return nil
	}

	targetValue := reflect.Indirect(reflect.ValueOf(dest))
	targetType := targetValue.Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		value := values.Get(strings.ToLower(field.Name))
		dest := targetValue.Field(i)
		if dest.Kind() != reflect.Ptr {
			dest = dest.Addr()
		}
		if err := unmarshalOne(value, dest.Interface()); err != nil {
			return err
		}
	}

	return nil
}

func encodeForm(src interface{}) ([]byte, error) {
	marshaler, ok := src.(FormMarshaler)
	if ok {
		values, err := marshaler.MarshalForm()
		if err != nil {
			return nil, err
		} else {
			return []byte(values.Encode()), nil
		}
	}

	values := new(url.Values)
	marshalOne := func(src interface{}) (string, error) {
		marshaler, ok := src.(encoding.TextMarshaler)
		if ok {
			value, err := marshaler.MarshalText()
			if err != nil {
				return "", err
			}
			return string(value), nil
		}

		return "", fmt.Errorf("cannot marshal %T", src)
	}

	targetType := reflect.TypeOf(src)
	targetValue := reflect.ValueOf(src)
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		dest := targetValue.Field(i)
		if value, err := marshalOne(dest.Interface()); err != nil {
			return nil, err
		} else {
			values.Set(strings.ToLower(field.Name), value)
		}
	}

	return []byte(values.Encode()), nil
}
