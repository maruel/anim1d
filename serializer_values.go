// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"errors"
	"fmt"
	"reflect"
)

// RegisterValue registers a known value.
func RegisterValue(v Value) error {
	t, ok := v.(JsonEncoding)
	if !ok {
		return errors.New("implementing JsonEncoding is required")
	}
	ty := t.Type()
	r := reflect.TypeOf(t).Elem()
	valuesTable[ty] = append(valuesTable[ty], r)
	if ty == Dict {
		valuesDictLookup[r.Name()] = r
	}
	return nil
}

// SValue is the serializable version of Value.
type SValue struct {
	Value
}

// Eval implements Value.
func (s *SValue) Eval(timeMS uint32, l int) int32 {
	if s.Value == nil {
		return 0
	}
	return s.Value.Eval(timeMS, l)
}

// UnmarshalJSON decodes a Value.
//
// It knows how to decode an arbitrary Value.
//
// If unmarshalling fails, the value is not touched.
func (s *SValue) UnmarshalJSON(b []byte) error {
	// Try to decode first as a int, then as a string, then as a dict. Not super
	// efficient but it works.
	if c, err := jsonUnmarshalInt32(b); err == nil {
		/*
			for _, v := range valuesTable[Int] {
				s.Value = Const(c)
			}
		*/
		return fmt.Errorf("unknown value %q", c)
	}
	if v, err := jsonUnmarshalString(b); err == nil {
		/*
			// It could be either a Percent or a Rand.
			if v == randKey {
				s.Value = &Rand{}
				return nil
			}
			if strings.HasSuffix(v, "%") {
				var p Percent
				if err = p.UnmarshalJSON(b); err == nil {
					s.Value = &p
				}
				return err
			}

			// Operations:
			if strings.HasPrefix(v, "+") {
				var o OpAdd
				if err = o.UnmarshalJSON(b); err == nil {
					s.Value = &o
				}
				return err
			}
			if strings.HasPrefix(v, "-") {
				var o OpAdd
				if err = o.UnmarshalJSON(b); err == nil {
					o.AddMS = -o.AddMS
					s.Value = &o
				}
				return err
			}
			if strings.HasPrefix(v, "%") {
				var o OpMod
				if err = o.UnmarshalJSON(b); err == nil {
					s.Value = &o
				}
				return err
			}
		*/
		/*
			for _, v := range valuesTable[Str] {
				s.Value = Const(c)
			}
		*/
		return fmt.Errorf("unknown value %q", v)
	}
	o, err := jsonUnmarshalWithType(b, valuesDictLookup, nil)
	if err != nil {
		return err
	}
	s.Value = o.(Value)
	return nil
}

// MarshalJSON includes the additional key "_type" to help with unmarshalling.
func (s *SValue) MarshalJSON() ([]byte, error) {
	if s.Value == nil {
		// nil value marshals to the constant 0.
		return []byte("0"), nil
	}
	return jsonMarshalWithType(s.Value)
}

//

var valuesDictLookup map[string]reflect.Type
var valuesTable map[Type][]reflect.Type
