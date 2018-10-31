// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// static is for patterns that do not change over time.

package patterns

import (
	"github.com/maruel/anim1d"
	"github.com/maruel/anim1d/math32"
)

// Rainbow renders rainbow colors.
type Rainbow struct {
	buf anim1d.Frame
}

// Render implements Pattern.
func (r *Rainbow) Render(pixels anim1d.Frame, timeMS uint32) {
	if len(r.buf) != len(pixels) {
		r.buf.Reset(len(pixels))
		const start = 380
		const end = 781
		const delta = end - start
		// TODO(maruel): Use integer arithmetic.
		scale := math32.Logn(2)
		step := 1. / float32(len(pixels))
		for i := range pixels {
			j := math32.Log1p(float32(len(pixels)-i-1)*step) / scale
			r.buf[i] = waveLength2RGB(int(start + delta*(1-j)))
		}
	}
	copy(pixels, r.buf)
}

func (r *Rainbow) String() string {
	return rainbowKey
}

// waveLengthToRGB returns a color over a rainbow.
//
// This code was inspired by public domain code on the internet.
func waveLength2RGB(w int) (c anim1d.Color) {
	switch {
	case w < 380:
	case w < 420:
		// Red peaks at 1/3 at 420.
		c.R = uint8(196 - (170*(440-w))/(440-380))
		c.B = uint8(26 + (229*(w-380))/(420-380))
	case w < 440:
		c.R = uint8((0x89 * (440 - w)) / (440 - 420))
		c.B = 255
	case w < 490:
		c.G = uint8((255 * (w - 440)) / (490 - 440))
		c.B = 255
	case w < 510:
		c.G = 255
		c.B = uint8((255 * (510 - w)) / (510 - 490))
	case w < 580:
		c.R = uint8((255 * (w - 510)) / (580 - 510))
		c.G = 255
	case w < 645:
		c.R = 255
		c.G = uint8((255 * (645 - w)) / (645 - 580))
	case w < 700:
		c.R = 255
	case w < 781:
		c.R = uint8(26 + (229*(780-w))/(780-700))
	default:
	}
	return
}

// Repeated repeats a Frame to fill the pixels.
type Repeated struct {
	Frame anim1d.Frame
}

// Render implements Pattern.
func (r *Repeated) Render(pixels anim1d.Frame, timeMS uint32) {
	if len(pixels) == 0 || len(r.Frame) == 0 {
		return
	}
	for i := 0; i < len(pixels); i += len(r.Frame) {
		copy(pixels[i:], r.Frame)
	}
}
