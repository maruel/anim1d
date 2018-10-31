// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package patterns

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/color"
	"image/png"
	"time"

	"github.com/maruel/anim1d"
)

const rainbowKey = "Rainbow"

func init() {
	var knownPatterns = []anim1d.Pattern{
		// Patterns
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
	for _, i := range knownPatterns {
		if err := anim1d.RegisterPattern(i); err != nil {
			panic(err)
		}
	}
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

func (r *Rainbow) Type() anim1d.Type {
	return anim1d.Str
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
	buf := make([]anim1d.Frame, maxY)
	for y := 0; y < maxY; y++ {
		buf[y] = make(anim1d.Frame, maxX)
	}
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			c := anim1d.Color{c1.R, c1.G, c1.B}
			if vertical {
				buf[x][y] = c
			} else {
				buf[y][x] = c
			}
		}
	}
	children := make([]anim1d.SPattern, maxY)
	for i, p := range buf {
		children[i].Pattern = p
	}
	return &Loop{
		Patterns: children,
		ShowMS:   uint32(frameDuration / time.Millisecond),
	}
}

//

// parsePatternString returns a anim1d.Pattern object out of the serialized JSON
// string.
func parsePatternString(b []byte) (anim1d.Pattern, error) {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return nil, err
	}
	// Could try to do one after the other? It's kind of a hack at the moment.
	if len(s) != 0 {
		switch s[0] {
		case '#':
			// "#RRGGBB"
			c := &anim1d.Color{}
			err := json.Unmarshal(b, c)
			return c, err
		case 'L':
			// "LRRGGBBRRGGBB..."
			var f anim1d.Frame
			err := json.Unmarshal(b, &f)
			return f, err
		case rainbowKey[0]:
			// "Rainbow"
			r := &Rainbow{}
			err := json.Unmarshal(b, r)
			return r, err
		}
	}
	return nil, errors.New("unrecognized pattern string")
}

func jsonUnmarshalString(b []byte) (string, error) {
	var s string
	err := json.Unmarshal(b, &s)
	return s, err
}
