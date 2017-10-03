// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// values contains all kind of non pattern types usable as values.

package anim1d

import (
	"math/rand"
)

// MinMax limits the value between a min and a max
func MinMax(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// MinMax32 limits the value between a min and a max
func MinMax32(v, min, max int32) int32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Values

// Value defines a value that may be constant or that may evolve over time.
type Value interface {
	Eval(timeMS uint32, l int) int32
}

// Const is a constant value.
type Const int32

// Eval implements Value.
func (c Const) Eval(timeMS uint32, l int) int32 {
	return int32(c)
}

// Percent is a percentage of the length. It is stored as a 16.16 fixed point.
type Percent int32

// Eval implements Value.
func (p Percent) Eval(timeMS uint32, l int) int32 {
	return int32(int64(l) * int64(p) / 65536)
}

// OpAdd adds a constant to timeMS.
type OpAdd struct {
	AddMS int32
}

// Eval implements Value.
func (o *OpAdd) Eval(timeMS uint32, l int) int32 {
	return int32(timeMS) + o.AddMS
}

// OpMod is a value that is cycling downward.
type OpMod struct {
	TickMS int32 // The cycling time. Maximum is ~25 days.
}

// Eval implements Value.
func (o *OpMod) Eval(timeMS uint32, l int) int32 {
	return int32(timeMS % uint32(o.TickMS))
}

// OpStep is a value that is cycling upward.
//
// It is useful for offsets that are increasing as a stepping function.
type OpStep struct {
	TickMS int32 // The cycling time. Maximum is ~25 days.
}

// Eval implements Value.
func (o *OpStep) Eval(timeMS uint32, l int) int32 {
	return int32(timeMS / uint32(o.TickMS) * uint32(o.TickMS))
}

// Rand is a value that pseudo-randomly changes every TickMS millisecond. If
// unspecified, changes every 60fps.
type Rand struct {
	TickMS int32 // The resolution at which the random value changes.
}

// Eval implements Value.
func (r *Rand) Eval(timeMS uint32, l int) int32 {
	m := uint32(r.TickMS)
	if m == 0 {
		m = 16
	}
	return int32(rand.NewSource(int64(timeMS / m)).Int63())
}

/*
// Equation evaluate an equation at every call.
type Equation struct {
	V string
	f func(timeMS uint32) int32
}

// Eval implements Value.
func (e *Equation) Eval(timeMS uint32) int32 {
	// Compiles the equation to an actual value and precompile it.
	if e.f == nil {
		e.f = func(timeMS uint32) int32 {
			return 0
		}
	}
	return e.f(timeMS)
}
*/

//

const epsilon = 1e-7
