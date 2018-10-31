// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Pattern is a interface to draw an animated line.
type Pattern interface {
	// Render fills the buffer with the image at this time frame.
	//
	// The image should be derived from timeMS, which is the time since this
	// pattern was started.
	//
	// Calling Render() with a nil pattern is valid. Patterns should be callable
	// without crashing with an object initialized with default values.
	//
	// timeMS will cycle after 49.7 days. The reason it's not using time.Duration
	// is that int64 calculation on ARM is very slow and abysmal on xtensa, which
	// this code is transpiled to.
	Render(pixels Frame, timeMS uint32)
}

// RegisterPattern registers a known Pattern.
func RegisterPattern(p Pattern) error {
	typ := Dict
	if t, ok := p.(JsonEncoding); ok {
		typ = t.Type()
	}
	// TODO: test for json.Marshaler?
	switch typ {
	case Int:
		return errors.New("pattern of type Int is unsupported")
	case Str:
		patternsStr = append(patternsStr, reflect.TypeOf(p).Elem())
	case Dict:
		r := reflect.TypeOf(p).Elem()
		patternsDictLookup[r.Name()] = r
	}
	return nil
}

// SPattern is a Pattern that can be serialized.
//
// It is only meant to be used in mixers.
type SPattern struct {
	Pattern
}

// Render implements Pattern.
func (s *SPattern) Render(pixels Frame, timeMS uint32) {
	if s.Pattern == nil {
		return
	}
	s.Pattern.Render(pixels, timeMS)
}

// MarshalJSON includes the additional key "_type" to help with unmarshalling.
func (s *SPattern) MarshalJSON() ([]byte, error) {
	if s.Pattern == nil {
		return []byte("{}"), nil
	}
	if i, ok := s.Pattern.(json.Marshaler); ok {
		return i.MarshalJSON()
	}
	return nil, errors.New("pattern doesn't implement json.Marshaler")
	//return jsonMarshalWithType(s.Pattern)
}

// UnmarshalJSON decodes a Pattern.
//
// It knows how to decode Color, Frame or other arbitrary Pattern.
//
// If unmarshalling fails, 's' is not touched.
func (s *SPattern) UnmarshalJSON(b []byte) error {
	// Try to decode first as a string, then as a dict. Not super efficient but
	// it works.
	if str, err := jsonUnmarshalString(b); err == nil {
		for _, r := range patternsStr {
			obj := reflect.New(r).Interface().(json.Unmarshaler)
			if obj.UnmarshalJSON(b) == nil {
				s.Pattern = obj.(Pattern)
				return nil
			}
		}
		return fmt.Errorf("unknown pattern %q", str)
	}

	o, err := jsonUnmarshalWithType(b, patternsDictLookup, nil)
	if err == nil && o == nil {
		return fmt.Errorf("unknown pattern %q", string(b))
	}
	return err
}

//

// UnmarshalJSON decodes the string "#RRGGBB" to the color.
//
// If unmarshalling fails, 'c' is not touched.
func (c *Color) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return c.FromString(s)
}

// MarshalJSON encodes the color as a string "#RRGGBB".
func (c *Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Color) Type() Type {
	return Str
}

// UnmarshalJSON decodes the string "LRRGGBB..." to the colors.
//
// If unmarshalling fails, 'f' is not touched.
func (f *Frame) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return f.FromString(s)
}

// MarshalJSON encodes the frame as a string "LRRGGBB...".
func (f Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f Frame) Type() Type {
	return Str
}

//

var patternsDictLookup = map[string]reflect.Type{}
var patternsStr []reflect.Type

func init() {
	var knownPatterns = []Pattern{
		&Color{},
		&Frame{},
	}
	for _, i := range knownPatterns {
		if err := RegisterPattern(i); err != nil {
			panic(err)
		}
	}

}
