// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// jsonUnmarshalDict unmarshals data into a map of interface{} without mangling
// int64.
func jsonUnmarshalDict(b []byte) (map[string]interface{}, error) {
	var tmp map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&tmp)
	return tmp, err
}

func jsonUnmarshalInt32(b []byte) (int32, error) {
	var i int32
	err := json.Unmarshal(b, &i)
	return i, err
}

func jsonUnmarshalString(b []byte) (string, error) {
	var s string
	err := json.Unmarshal(b, &s)
	return s, err
}

func jsonMarshalWithType(v interface{}) ([]byte, error) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return jsonMarshalWithTypeName(v, t.Elem().Name())
	default:
		return jsonMarshalWithTypeName(v, t.Name())
	}
}

func jsonMarshalWithTypeName(v interface{}, name string) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil || (len(b) != 0 && b[0] != '{') {
		// Special case check for custom marshallers that do not encode as a dict.
		return b, err
	}
	// Inject "_type".
	tmp, err := jsonUnmarshalDict(b)
	if err != nil {
		return nil, err
	}
	tmp["_type"] = name
	return json.Marshal(tmp)
}

func jsonUnmarshalWithType(b []byte, lookup map[string]reflect.Type, null interface{}) (interface{}, error) {
	tmp, err := jsonUnmarshalDict(b)
	if err != nil {
		return nil, err
	}
	if len(tmp) == 0 {
		// No error but nothing was present. Treat "{}" as equivalent encoding for
		// null.
		return null, nil
	}
	n, ok := tmp["_type"]
	if !ok {
		return nil, errors.New("missing value type")
	}
	name, ok := n.(string)
	if !ok {
		return nil, errors.New("invalid value type")
	}
	// "_type" will be ignored, no need to reencode the dict to json.
	return parseDictToType(name, b, lookup)
}

// parseDictToType decodes an object out of the serialized JSON dict.
func parseDictToType(name string, b []byte, lookup map[string]reflect.Type) (interface{}, error) {
	t, ok := lookup[name]
	if !ok {
		return nil, fmt.Errorf("type %#v not found", name)
	}
	v := reflect.New(t).Interface()
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}
