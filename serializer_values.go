// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const randKey = "rand"

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
// It knows how to decode Const or other arbitrary Value.
func (v *SValue) UnmarshalJSON(b []byte) error {
	// Try to decode first as a int, then as a string, then as a dict. Not super
	// efficient but it works.
	if c, err := jsonUnmarshalInt32(b); err == nil {
		c2 := Const(c)
		v.Value = &c2
		return nil
	}
	/*
		if v, err := jsonUnmarshalString(b); err == nil {
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
			return fmt.Errorf("unknown value %q", v)
		}
	*/
	if s, err := jsonUnmarshalString(b); err == nil {
		v2, err := parseValueString(s)
		if err != nil {
			return err
		}
		v.Value = v2
		return nil
	}
	o, err := jsonUnmarshalWithType(b, valuesLookup, nil)
	if err != nil {
		return err
	}
	v.Value = o.(Value)
	return nil
}

// MarshalJSON includes the additional key "_type" to help with unmarshalling.
func (v *SValue) MarshalJSON() ([]byte, error) {
	if v.Value == nil {
		// nil value marshals to the constant 0.
		return []byte("0"), nil
	}
	return jsonMarshalWithType(v.Value)
}

// UnmarshalJSON decodes the "t".
func (t *TimeMS) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return t.FromString(s)
}

func (t *TimeMS) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON decodes the "t".
func (l *Length) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return l.FromString(s)
	// TODO(maruel): Do a jsonUnmarshalWithType(b, valuesLookup, nil) for each type?
}

func (l *Length) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

// UnmarshalJSON decodes the int to the const.
//
// If unmarshalling fails, 'c' is not touched.
func (c *Const) UnmarshalJSON(b []byte) error {
	i, err := jsonUnmarshalInt32(b)
	if err == nil {
		*c = Const(i)
		return nil
	}
	if s, err := jsonUnmarshalString(b); err == nil {
		return c.FromString(s)
	}
	return err
}

// MarshalJSON encodes the const as a int.
func (c *Const) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*c))
}

// UnmarshalJSON decodes the percent in the form of a string.
//
// If unmarshalling fails, 'p' is not touched.
func (p *Percent) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return p.FromString(s)
}

// MarshalJSON encodes the percent as a string.
func (p *Percent) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// UnmarshalJSON decodes the add in the form of a string.
//
// If unmarshalling fails, 'o' is not touched.
func (o *OpAdd) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return o.FromString(s)
}

// MarshalJSON encodes the add as a string.
func (o *OpAdd) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *OpGroup) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return o.FromString(s)
}

func (o *OpGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

// UnmarshalJSON decodes the mod in the form of a string.
//
// If unmarshalling fails, 'o' is not touched.
func (o *OpMod) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return o.FromString(s)
}

// MarshalJSON encodes the mod as a string.
func (o *OpMod) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *OpMul) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return o.FromString(s)
}

func (o *OpMul) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *OpSub) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return o.FromString(s)
}

func (o *OpSub) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

// UnmarshalJSON decodes the string to the rand.
//
// If unmarshalling fails, 'r' is not touched.
func (r *Rand) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err == nil {
		return r.FromString(s)
	}
	v, err := jsonUnmarshalWithType(b, valuesLookup, nil)
	if err == nil {
		if r2, ok := v.(*Rand); ok {
			*r = *r2
			return nil
		}
		return errors.New("rand: internal error")
	}
	return err
}

// MarshalJSON encodes the rand as a string.
func (r *Rand) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON is because MovePerHour is a superset of SValue.
func (m *MovePerHour) UnmarshalJSON(b []byte) error {
	var s SValue
	if err := s.UnmarshalJSON(b); err != nil {
		return err
	}
	*m = MovePerHour(s)
	return nil
}

// MarshalJSON is because MovePerHour is a superset of SValue.
func (m *MovePerHour) MarshalJSON() ([]byte, error) {
	s := SValue{m.Value}
	return s.MarshalJSON()
}

//

// valuesLookup lists all the known values that can be instantiated.
var valuesLookup map[string]reflect.Type

var knownValues = []Value{
	&TimeMS{},
	&Length{},
	new(Const),
	new(Percent),
	&OpAdd{},
	&OpGroup{},
	&OpMod{},
	&OpMul{},
	&OpStep{},
	&OpSub{},
	&Rand{},
}

func init() {
	valuesLookup = make(map[string]reflect.Type, len(knownValues))
	for _, i := range knownValues {
		r := reflect.TypeOf(i).Elem()
		valuesLookup[r.Name()] = r
	}
}

// parseValueString returns a Value object out of the string.
func parseValueString(s string) (Value, error) {
	// TODO(maruel): Smarter. Use brute force in the meantime.
	for _, t := range knownValues {
		v := reflect.New(reflect.TypeOf(t).Elem()).Interface().(Value)
		if err := v.FromString(s); err == nil {
			return v, nil
		}
	}
	return nil, fmt.Errorf("value: failed to parse %q", s)
}
