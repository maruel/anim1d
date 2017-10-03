// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/maruel/anim1d"
)

const randKey = "rand"

func init() {
	knownValues := []anim1d.Value{
		new(Const),
		new(Percent),
		&OpAdd{},
		&OpMod{},
		&OpStep{},
		&Rand{},
	}
	for _, v := range knownValues {
		anim1d.RegisterValue(v)
	}
}

/*
// UnmarshalJSON decodes a Value.
//
// It knows how to decode Const or other arbitrary Value.
//
// If unmarshalling fails, 'f' is not touched.
func (s *SValue) UnmarshalJSON(b []byte) error {
	// Try to decode first as a int, then as a string, then as a dict. Not super
	// efficient but it works.
	if c, err := jsonUnmarshalInt32(b); err == nil {
		s.Value = Const(c)
		return nil
	}
	if v, err := jsonUnmarshalString(b); err == nil {
		// It could be either a Percent or a Rand.
		if v == randKey {
			s.Value = &Rand{}
			return nil
		}
		if strings.HasPrefix(v, "+") {
			var o OpAdd
			if err := o.UnmarshalJSON(b); err == nil {
				s.Value = &o
			}
			return err
		}
		if strings.HasPrefix(v, "-") {
			var o OpAdd
			if err := o.UnmarshalJSON(b); err == nil {
				o.AddMS = -o.AddMS
				s.Value = &o
			}
			return err
		}
		if strings.HasPrefix(v, "%") {
			var o OpMod
			if err := o.UnmarshalJSON(b); err == nil {
				s.Value = &o
			}
			return err
		}
		if strings.HasSuffix(v, "%") {
			var p Percent
			if err := p.UnmarshalJSON(b); err == nil {
				s.Value = &p
			}
			return err
		}
		return fmt.Errorf("unknown value %q", v)
	}
	o, err := jsonUnmarshalWithType(b, valuesLookup, nil)
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
*/

func (c *Const) Type() anim1d.Type {
	return anim1d.Int
}

// UnmarshalJSON decodes the int to the const.
//
// If unmarshalling fails, 'c' is not touched.
func (c *Const) UnmarshalJSON(b []byte) error {
	i, err := jsonUnmarshalInt32(b)
	if err != nil {
		return err
	}
	*c = Const(i)
	return err
}

// MarshalJSON encodes the const as a int.
func (c *Const) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*c))
}

func (p *Percent) Type() anim1d.Type {
	return anim1d.Str
}

// UnmarshalJSON decodes the percent in the form of a string.
//
// If unmarshalling fails, 'p' is not touched.
func (p *Percent) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(s, "%") {
		return errors.New("percent must end with %")
	}
	f, err := strconv.ParseFloat(s[:len(s)-1], 32)
	if err == nil {
		// Convert back to fixed point.
		*p = Percent(int32(f * 655.36))
	}
	return err
}

// MarshalJSON encodes the percent as a string.
func (p *Percent) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatFloat(float64(*p)/655.36, 'g', 4, 32) + "%")
}

func (o *OpAdd) Type() anim1d.Type {
	return anim1d.Str
}

// UnmarshalJSON decodes the add in the form of a string.
//
// If unmarshalling fails, 'o' is not touched.
func (o *OpAdd) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	i := int64(0)
	if strings.HasPrefix(s, "+") {
		i, err = strconv.ParseInt(s[1:], 10, 32)
	} else if strings.HasPrefix(s, "-") {
		i, err = strconv.ParseInt(s, 10, 32)
	} else {
		return errors.New("add: must start with + or -")
	}
	if err == nil {
		o.AddMS = int32(i)
	}
	if i < 0 {
		return errors.New("add: value must be positive")
	}
	return err
}

// MarshalJSON encodes the add as a string.
func (o *OpAdd) MarshalJSON() ([]byte, error) {
	if o.AddMS >= 0 {
		return json.Marshal("+" + strconv.FormatInt(int64(o.AddMS), 10))
	}
	return json.Marshal(strconv.FormatInt(int64(o.AddMS), 10))
}

func (o *OpMod) Type() anim1d.Type {
	return anim1d.Str
}

// UnmarshalJSON decodes the mod in the form of a string.
//
// If unmarshalling fails, 'o' is not touched.
func (o *OpMod) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(s, "%") {
		return errors.New("mod: must start with %")
	}
	i, err := strconv.ParseInt(s[1:], 10, 32)
	if err == nil {
		o.TickMS = int32(i)
	}
	if i < 0 {
		return errors.New("mod: value must be positive")
	}
	return err
}

// MarshalJSON encodes the mod as a string.
func (o *OpMod) MarshalJSON() ([]byte, error) {
	return json.Marshal("%" + strconv.FormatInt(int64(o.TickMS), 10))
}

func (r *Rand) Type() anim1d.Type {
	return anim1d.Str
}

// UnmarshalJSON decodes the string to the rand.
//
// If unmarshalling fails, 'r' is not touched.
func (r *Rand) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err == nil {
		// Shortcut.
		if s != randKey {
			return errors.New("invalid format")
		}
		r.TickMS = 0
		return nil
	}
	// SValue.UnmarshalJSON would handle it but implement it here so calling
	// UnmarshalJSON on a concrete instance still work. The issue is that we do
	// not want to recursively call ourselves so create a temporary type.
	type tmpRand Rand
	var r2 tmpRand
	if err := json.Unmarshal(b, &r2); err != nil {
		return err
	}
	*r = Rand(r2)
	return nil
}

// MarshalJSON encodes the rand as a string.
func (r *Rand) MarshalJSON() ([]byte, error) {
	if r.TickMS == 0 {
		// Shortcut.
		return json.Marshal(randKey)
	}
	type tmpRand Rand
	r2 := tmpRand(*r)
	return jsonMarshalWithTypeName(r2, "Rand")
}

func (m *MovePerHour) Type() anim1d.Type {
	return anim1d.Str
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
