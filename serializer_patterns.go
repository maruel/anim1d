// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/color"
	"image/png"
	"reflect"
	"time"
)

const rainbowKey = "Rainbow"

// patternsLookup lists all known patterns that can be instantiated.
var patternsLookup map[string]reflect.Type

var knownPatterns = []Pattern{
	// Patterns
	&Color{},
	&Frame{},
	&Rainbow{},
	&Repeated{},
	&Aurore{},
	&NightStars{},
	&Lightning{},
	&WishingStar{},
	// Mixers
	&Gradient{},
	&Split{},
	&Transition{},
	&Loop{},
	&Chronometer{},
	&Rotate{},
	&PingPong{},
	&Crop{},
	&Subset{},
	&Dim{},
	&Add{},
	&Scale{},
}

func init() {
	patternsLookup = make(map[string]reflect.Type, len(knownPatterns))
	for _, i := range knownPatterns {
		r := reflect.TypeOf(i).Elem()
		patternsLookup[r.Name()] = r
	}
}

// SPattern

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

// UnmarshalJSON decodes a Pattern.
//
// It knows how to decode Color, Frame or other arbitrary Pattern.
//
// If unmarshalling fails, 's' is not touched.
func (s *SPattern) UnmarshalJSON(b []byte) error {
	// Try to decode first as a string, then as a dict. Not super efficient but
	// it works.
	if p2, err := parsePatternString(b); err != nil {
		return err
	} else if p2 != nil {
		s.Pattern = p2
		return nil
	}
	o, err := jsonUnmarshalWithType(b, patternsLookup, nil)
	if err != nil {
		return err
	}
	if o == nil {
		s.Pattern = nil
	} else {
		s.Pattern = o.(Pattern)
	}
	return nil
}

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

// UnmarshalJSON decodes the string "Rainbow" to the rainbow.
func (r *Rainbow) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	if s != rainbowKey {
		return errors.New("invalid color string")
	}
	return err
}

// MarshalJSON encodes the rainbow as a string "Rainbow".
func (r *Rainbow) MarshalJSON() ([]byte, error) {
	return json.Marshal(rainbowKey)
}

// MarshalJSON includes the additional key "_type" to help with unmarshalling.
func (s *SPattern) MarshalJSON() ([]byte, error) {
	if s.Pattern == nil {
		return []byte("{}"), nil
	}
	return jsonMarshalWithType(s.Pattern)
}

// LoadPNG loads a PNG file and creates a Loop out of the lines.
//
// If vertical is true, rotate the image by 90°.
func LoadPNG(content []byte, frameDuration time.Duration, vertical bool) *Loop {
	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	bounds := img.Bounds()
	maxY := bounds.Max.Y
	maxX := bounds.Max.X
	if vertical {
		// Invert axes.
		maxY, maxX = maxX, maxY
	}
	buf := make([]Frame, maxY)
	for y := 0; y < maxY; y++ {
		buf[y] = make(Frame, maxX)
	}
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			c := Color{c1.R, c1.G, c1.B}
			if vertical {
				buf[x][y] = c
			} else {
				buf[y][x] = c
			}
		}
	}
	children := make([]SPattern, maxY)
	for i, p := range buf {
		children[i].Pattern = p
	}
	return &Loop{
		Patterns: children,
		ShowMS:   uint32(frameDuration / time.Millisecond),
	}
}

//

// parsePatternString returns a Pattern object out of the serialized JSON
// string.
func parsePatternString(b []byte) (Pattern, error) {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return nil, nil
	}
	// Could try to do one after the other? It's kind of a hack at the moment.
	if len(s) != 0 {
		switch s[0] {
		case '#':
			// "#RRGGBB"
			c := &Color{}
			err := json.Unmarshal(b, c)
			return c, err
		case 'L':
			// "LRRGGBBRRGGBB..."
			var f Frame
			err := json.Unmarshal(b, &f)
			return f, err
		case rainbowKey[0]:
			// "Rainbow"
			r := &Rainbow{}
			err := json.Unmarshal(b, r)
			return r, err
		}
	}
	return nil, errors.New("unrecognized pattern string, should start with '#', 'L' or be a known constant")
}
