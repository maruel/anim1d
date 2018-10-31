// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package patterns

import (
	"reflect"
	"testing"

	"github.com/maruel/anim1d"
)

func TestInterpolationEmpty(t *testing.T) {
	b := make(anim1d.Frame, 1)
	for _, v := range []Interpolation{Interpolation(""), NearestSkip, Nearest, Linear} {
		v.Scale(nil, nil)
		v.Scale(nil, b)
		v.Scale(b, nil)
	}
}

func TestInterpolation(t *testing.T) {
	red := anim1d.Color{0xFF, 0x00, 0x00}
	green := anim1d.Color{0x00, 0xFF, 0x00}
	blue := anim1d.Color{0x00, 0x00, 0xFF}
	yellow := anim1d.Color{0xFF, 0xFF, 0x00}
	cyan := anim1d.Color{0x00, 0xFF, 0xFF}
	magenta := anim1d.Color{0xFF, 0x00, 0xFF}
	white := anim1d.Color{0xFF, 0xFF, 0xFF}
	black := anim1d.Color{}
	input := anim1d.Frame{red, green, blue, yellow, cyan, magenta, white}
	data := []struct {
		s        Interpolation
		input    anim1d.Frame
		expected anim1d.Frame
	}{
		{
			NearestSkip,
			input,
			anim1d.Frame{red, black, green, black, blue, black, yellow, black, cyan, black, magenta, black, white},
		},
		{NearestSkip, input, anim1d.Frame{yellow}},
		{NearestSkip, input, anim1d.Frame{green, magenta}},
		{NearestSkip, input, anim1d.Frame{green, yellow, magenta}},
		{
			Nearest,
			input,
			anim1d.Frame{red, red, green, green, blue, blue, yellow, yellow, cyan, cyan, magenta, magenta, white, white},
		},
		{Nearest, input, anim1d.Frame{yellow}},
		{Nearest, input, anim1d.Frame{green, magenta}},
		{Nearest, input, anim1d.Frame{green, yellow, magenta}},
		// TODO(maruel): This is broken.
		/*{
			Linear,
			input,
			anim1d.Frame{red, red, green, green, blue, blue, yellow, yellow, cyan, cyan, magenta, magenta, white, white},
		},*/
		{Linear, input, anim1d.Frame{anim1d.Color{0x80, 0xFF, 0x7F}}},
		{Linear, input, anim1d.Frame{anim1d.Color{0x00, 0x80, 0x7F}, anim1d.Color{0xFF, 0x7F, 0xFF}}},
		{Linear, input, anim1d.Frame{anim1d.Color{0x0, 0x80, 0x7F}, anim1d.Color{0x80, 0xFF, 0x7F}, anim1d.Color{0xFF, 0x7F, 0xFF}}},
	}
	for i, line := range data {
		out := make(anim1d.Frame, len(line.expected))
		line.s.Scale(line.input, out)
		if !reflect.DeepEqual(out, line.expected) {
			t.Fatalf("%d: %v != %v", i, out, line.expected)
		}
	}
}
