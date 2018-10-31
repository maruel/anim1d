// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package patterns

import (
	"encoding/json"
	"testing"

	"github.com/maruel/anim1d"
)

func TestGradient(t *testing.T) {
	a := &anim1d.Color{0x10, 0x10, 0x10}
	b := &anim1d.Color{0x20, 0x20, 0x20}
	testFrame(t, &Gradient{Left: anim1d.SPattern{a}, Right: anim1d.SPattern{b}, Curve: anim1d.Direct}, expectation{0, anim1d.Frame{{0x18, 0x18, 0x18}}})
	testFrame(t, &Gradient{Left: anim1d.SPattern{a}, Right: anim1d.SPattern{b}, Curve: anim1d.Direct}, expectation{0, anim1d.Frame{{0x10, 0x10, 0x10}, {0x20, 0x20, 0x20}}})
	testFrame(t, &Gradient{Left: anim1d.SPattern{a}, Right: anim1d.SPattern{b}, Curve: anim1d.Direct}, expectation{0, anim1d.Frame{{0x10, 0x10, 0x10}, {0x18, 0x18, 0x18}, {0x20, 0x20, 0x20}}})
}

func TestTransition(t *testing.T) {
	// TODO(maruel): Add.
}

func TestLoop(t *testing.T) {
	// TODO(maruel): Add.
}

func TestRotate(t *testing.T) {
	a := anim1d.Color{10, 10, 10}
	b := anim1d.Color{20, 20, 20}
	c := anim1d.Color{30, 30, 30}
	p := &Rotate{Child: anim1d.SPattern{anim1d.Frame{a, b, c}}, MovePerHour: anim1d.MovePerHour{anim1d.Const(360000)}}
	e := []expectation{
		{0, anim1d.Frame{a, b, c}},
		{5, anim1d.Frame{a, b, c}},
		{10, anim1d.Frame{c, a, b}},
		{20, anim1d.Frame{b, c, a}},
		{30, anim1d.Frame{a, b, c}},
		{40, anim1d.Frame{c, a, b}},
		{50, anim1d.Frame{b, c, a}},
		{60, anim1d.Frame{a, b, c}},
	}
	testFrames(t, p, e)
}

func TestRotateRev(t *testing.T) {
	// Works in reverse too.
	a := anim1d.Color{10, 10, 10}
	b := anim1d.Color{20, 20, 20}
	c := anim1d.Color{30, 30, 30}
	p := &Rotate{Child: anim1d.SPattern{anim1d.Frame{a, b, c}}, MovePerHour: anim1d.MovePerHour{anim1d.Const(-360000)}}
	e := []expectation{
		{0, anim1d.Frame{a, b, c}},
		{5, anim1d.Frame{a, b, c}},
		{10, anim1d.Frame{b, c, a}},
		{20, anim1d.Frame{c, a, b}},
		{30, anim1d.Frame{a, b, c}},
		{40, anim1d.Frame{b, c, a}},
		{50, anim1d.Frame{c, a, b}},
		{60, anim1d.Frame{a, b, c}},
	}
	testFrames(t, p, e)
}

func TestChronometer(t *testing.T) {
	r := anim1d.Color{0xff, 0x00, 0x00}
	g := anim1d.Color{0x00, 0xff, 0x00}
	b := anim1d.Color{0x00, 0x00, 0xff}
	p := &Chronometer{Child: anim1d.SPattern{anim1d.Frame{{}, r, g, b}}}
	exp := []expectation{
		{0, anim1d.Frame{r, {}, {}, {}, {}, {}}},                        // 0:00:00
		{1000 * 10, anim1d.Frame{g, r, {}, {}, {}, {}}},                 // 0:00:10
		{1000 * 20, anim1d.Frame{g, {}, r, {}, {}, {}}},                 // 0:00:20
		{1000 * 60, anim1d.Frame{r, {}, {}, {}, {}, {}}},                // 0:01:00
		{1000 * 600, anim1d.Frame{r, g, {}, {}, {}, {}}},                // 0:10:00
		{1000 * 3600, anim1d.Frame{r, b, {}, {}, {}, {}}},               // 1:00:00
		{1000 * (3600 + 20*60 + 30), anim1d.Frame{{}, b, g, r, {}, {}}}, // 1:20:30
	}
	testFrames(t, p, exp)
}

func TestPingPong(t *testing.T) {
	a := anim1d.Color{0x10, 0x10, 0x10}
	b := anim1d.Color{0x20, 0x20, 0x20}
	c := anim1d.Color{0x30, 0x30, 0x30}
	d := anim1d.Color{0x40, 0x40, 0x40}
	e := anim1d.Color{0x50, 0x50, 0x50}
	f := anim1d.Color{0x60, 0x60, 0x60}

	p := &PingPong{Child: anim1d.SPattern{anim1d.Frame{a, b}}, MovePerHour: anim1d.MovePerHour{anim1d.Const(360000)}}
	exp := []expectation{
		{0, anim1d.Frame{a, b, {}}},
		{5, anim1d.Frame{a, b, {}}},
		{10, anim1d.Frame{b, a, {}}},
		{20, anim1d.Frame{{}, b, a}},
		{30, anim1d.Frame{{}, a, b}},
		{40, anim1d.Frame{a, b, {}}},
		{50, anim1d.Frame{b, a, {}}},
		{60, anim1d.Frame{{}, b, a}},
	}
	testFrames(t, p, exp)

	p = &PingPong{Child: anim1d.SPattern{anim1d.Frame{a, b, c, d, e, f}}, MovePerHour: anim1d.MovePerHour{anim1d.Const(3600)}}
	exp = []expectation{
		{0, anim1d.Frame{a, b, c, d}},
		{500, anim1d.Frame{a, b, c, d}},
		{1000, anim1d.Frame{b, a, d, e}},
		{2000, anim1d.Frame{c, b, a, f}},
		{3000, anim1d.Frame{d, c, b, a}},
		{4000, anim1d.Frame{e, d, a, b}},
		{5000, anim1d.Frame{f, a, b, c}},
		{6000, anim1d.Frame{a, b, c, d}},
	}
	testFrames(t, p, exp)
}

func TestCrop(t *testing.T) {
	// Crop skips the beginning and the end of the source.
	f := anim1d.Frame{
		{0x10, 0x10, 0x10},
		{0x20, 0x20, 0x20},
		{0x30, 0x30, 0x30},
	}
	p := &Crop{Child: anim1d.SPattern{f}, Before: anim1d.SValue{anim1d.Const(1)}, After: anim1d.SValue{anim1d.Const(2)}}
	testFrame(t, p, expectation{0, f[1:3]})
}

func TestSubset(t *testing.T) {
	// Subset skips the beginning and the end of the destination.
	f := anim1d.Frame{
		{0x10, 0x10, 0x10},
		{0x20, 0x20, 0x20},
		{0x30, 0x30, 0x30},
	}
	p := &Subset{Child: anim1d.SPattern{f}, Offset: anim1d.SValue{anim1d.Const(1)}, Length: anim1d.SValue{anim1d.Const(2)}}
	// Skip the beginning and the end of the destination.
	expected := anim1d.Frame{
		{},
		{0x10, 0x10, 0x10},
		{0x20, 0x20, 0x20},
		{},
	}
	testFrame(t, p, expectation{0, expected})
}

func TestDim(t *testing.T) {
	p := &Dim{Child: anim1d.SPattern{&anim1d.Color{0x60, 0x60, 0x60}}, Intensity: anim1d.SValue{anim1d.Const(127)}}
	testFrame(t, p, expectation{0, anim1d.Frame{{0x2f, 0x2f, 0x2f}}})
}

func TestAdd(t *testing.T) {
	a := anim1d.Color{0x60, 0x60, 0x60}
	b := anim1d.Color{0x10, 0x20, 0x30}
	p := &Add{Patterns: []anim1d.SPattern{{&a}, {&b}}}
	testFrame(t, p, expectation{0, anim1d.Frame{{0x70, 0x80, 0x90}}})
}

func TestScale(t *testing.T) {
	f := anim1d.Frame{{0x60, 0x60, 0x60}, {0x10, 0x20, 0x30}}
	p := &Scale{Child: anim1d.SPattern{f}, Interpolation: NearestSkip, RatioMilli: anim1d.SValue{anim1d.Const(667)}}
	expected := anim1d.Frame{{0x60, 0x60, 0x60}, {}, {0x10, 0x20, 0x30}}
	testFrame(t, p, expectation{0, expected})
}

//

type expectation struct {
	offsetMS uint32
	colors   anim1d.Frame
}

func testFrames(t *testing.T, p anim1d.Pattern, expectations []expectation) {
	var pixels anim1d.Frame
	for frame, e := range expectations {
		pixels.Reset(len(e.colors))
		p.Render(pixels, e.offsetMS)
		if !e.colors.IsEqual(pixels) {
			x := marshalPattern(e.colors)
			t.Fatalf("frame=%d bad expectation:\n%s\n%s", frame, x, marshalPattern(pixels))
		}
	}
}

func testFrame(t *testing.T, p anim1d.Pattern, e expectation) {
	pixels := make(anim1d.Frame, len(e.colors))
	p.Render(pixels, e.offsetMS)
	if !e.colors.IsEqual(pixels) {
		t.Fatalf("%s != %s", marshalPattern(e.colors), marshalPattern(pixels))
	}
}

// marshalPattern is a shorthand to JSON encode a pattern.
func marshalPattern(p anim1d.Pattern) []byte {
	b, err := json.Marshal(&anim1d.SPattern{p})
	if err != nil {
		panic(err)
	}
	return b
}
