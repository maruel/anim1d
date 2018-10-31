// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package patterns

import "github.com/maruel/anim1d"

// Interpolation specifies a way to scales a pixel strip.
type Interpolation string

// All the kinds of interpolations.
const (
	NearestSkip Interpolation = "nearestskip" // Selects the nearest pixel but when upscaling, skips on missing pixels.
	Nearest     Interpolation = "nearest"     // Selects the nearest pixel, gives a blocky view.
	Linear      Interpolation = "linear"      // Linear interpolation, recommended and default value.
)

// Scale interpolates a frame into another using integers as much as possible
// for reasonable performance.
func (i Interpolation) Scale(in, out anim1d.Frame) {
	li := len(in)
	lo := len(out)
	if li == 0 || lo == 0 {
		return
	}
	switch i {
	case NearestSkip:
		if li < lo {
			// Do not touch skipped pixels.
			for i, p := range in {
				out[(i*lo+lo/2)/li] = p
			}
			return
		}
		// When the destination is smaller than the source, Nearest and NearestSkip
		// have the same behavior.
		fallthrough
	case Nearest, "":
		fallthrough
	default:
		for i := range out {
			out[i] = in[(i*li+li/2)/lo]
		}
	case Linear:
		for i := range out {
			x := (i*li + li/2) / lo
			c := in[x]
			if x < li-1 {
				gradient := uint8(127)
				c.Mix(in[x+1], gradient)
			}
			out[i] = c
			//a := in[(i*li+li/2)/lo]
			//b := in[(i*li+li/2)/lo]
			//out[i] = (a + b) / 2
		}
	}
}
