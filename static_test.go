// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"testing"
)

func TestColor(t *testing.T) {
	p := &Color{255, 255, 255}
	e := []expectation{{3000, Frame{{255, 255, 255}}}}
	testFrames(t, p, e)
}

func TestColor_Add(t *testing.T) {
	a := Color{255, 255, 255}
	b := Color{1, 1, 1}
	a.Add(b)
	expected := Color{255, 255, 255}
	if a != expected {
		t.Fatal(a)
	}
}

func TestColor_Dim(t *testing.T) {
	a := Color{200, 150, 100}
	a.Dim(100)
	expected := Color{78, 58, 39}
	if a != expected {
		t.Fatal(a)
	}
}

func TestColor_Mix(t *testing.T) {
	white := Color{255, 255, 255}
	black := Color{0, 0, 0}

	data := []struct {
		start    Color
		new      Color
		mix      uint8
		expected Color
	}{
		{white, black, 0, white},
		{white, black, 255, black},
		// Rounding is hard.
		{white, black, 128, Color{127, 127, 127}},
		{black, white, 128, Color{128, 128, 128}},
		// Test for overflow.
		{white, white, 0, white},
		{white, white, 128, white},
		{white, white, 255, white},
		// Verify channels.
		{Color{0x10, 0x20, 0x30}, black, 0, Color{0x10, 0x20, 0x30}},
		{black, Color{0x10, 0x20, 0x30}, 255, Color{0x10, 0x20, 0x30}},
	}
	for i, line := range data {
		c := line.start
		c.Mix(line.new, line.mix)
		if c != line.expected {
			t.Fatalf("%d: %v.Mix(%v, %v) = %v; expected %v", i, line.start, line.new, line.mix, c, line.expected)
		}
	}
}

func TestColor_FromString(t *testing.T) {
	c := Color{}
	if err := c.FromString("123456"); err == nil {
		t.Fail()
	}
	if err := c.FromString("#123456"); err != nil {
		t.Fail()
	}
	if c.R != 0x12 || c.G != 0x34 || c.B != 0x56 {
		t.Fail()
	}
}

func TestColor_FromRGBString(t *testing.T) {
	c := Color{}
	if err := c.FromRGBString("12345"); err == nil {
		t.Fail()
	}
	if err := c.FromRGBString("1g3456"); err == nil {
		t.Fail()
	}
	if err := c.FromRGBString("123g56"); err == nil {
		t.Fail()
	}
	if err := c.FromRGBString("12345g"); err == nil {
		t.Fail()
	}
}

func TestFrame_Render(t *testing.T) {
	f := Frame{Color{0x10, 0x20, 0x30}, Color{0x40, 0x50, 0x60}}
	d := make(Frame, 2)
	f.Render(d, 1000)
	if !d.IsEqual(f) {
		t.Fatal(d)
	}
}

func TestFrame_Dim(t *testing.T) {
	f := Frame{Color{0x10, 0x20, 0x30}, Color{0x40, 0x50, 0x60}}
	f.Dim(128)
	expected := Frame{Color{0x10 / 2, 0x20 / 2, 0x30 / 2}, Color{0x40 / 2, 0x50 / 2, 0x60 / 2}}
	if !f.IsEqual(expected) {
		t.Fatal(f)
	}
}

func TestFrame_Add(t *testing.T) {
	f := Frame{Color{0x10, 0x20, 0x30}, Color{0x40, 0x50, 0x60}}
	d := Frame{Color{3, 2, 1}, Color{6, 5, 4}}
	f.Add(d)
	expected := Frame{Color{0x13, 0x22, 0x31}, Color{0x46, 0x55, 0x64}}
	if !f.IsEqual(expected) {
		t.Fatal(f)
	}
}

func TestFrame_Mix(t *testing.T) {
	f := Frame{Color{0x10, 0x20, 0x30}, Color{0x40, 0x50, 0x60}}
	d := Frame{Color{3, 2, 1}, Color{6, 5, 4}}
	f.Mix(d, 255/3)
	expected := Frame{Color{0x0C, 0x16, 0x21}, Color{0x2d, 0x37, 0x42}}
	if !f.IsEqual(expected) {
		t.Fatal(f)
	}
}

func TestFrame_Reset(t *testing.T) {
	f := Frame{Color{0x10, 0x20, 0x30}, Color{0x40, 0x50, 0x60}}
	f.Reset(2)
	expected := make(Frame, 2)
	if !f.IsEqual(expected) {
		t.Fatal(f)
	}
}

func TestFrame_FromString(t *testing.T) {
	f := Frame{}
	if err := f.FromString("123456"); err == nil {
		t.Fail()
	}
	if err := f.FromString("L1g3456"); err == nil {
		t.Fail()
	}
	if err := f.FromString("L123456"); err != nil {
		t.Fail()
	}
	if !f.IsEqual(Frame{{0x12, 0x34, 0x56}}) {
		t.Fail()
	}
}

func TestFrame_ToRGB(t *testing.T) {
	f := Frame{{0x12, 0x34, 0x56}}
	b := [3]byte{}
	f.ToRGB(b[:])
	if !bytes.Equal(b[:], []byte{0x12, 0x34, 0x56}) {
		t.Fail()
	}
}

func TestFrame_IsEqual(t *testing.T) {
	f1 := Frame{{0x12, 0x34, 0x56}}
	f2 := Frame{{0x12, 0x34, 0x56}, {}}
	if f1.IsEqual(f2) || f2.IsEqual(f1) {
		t.Fail()
	}
	f2 = Frame{{}}
	if f1.IsEqual(f2) || f2.IsEqual(f1) {
		t.Fail()
	}
}

//

type expectation struct {
	offsetMS uint32
	colors   Frame
}

func testFrames(t *testing.T, p Pattern, expectations []expectation) {
	var pixels Frame
	for frame, e := range expectations {
		pixels.Reset(len(e.colors))
		p.Render(pixels, e.offsetMS)
		if !e.colors.IsEqual(pixels) {
			x := marshalPattern(e.colors)
			t.Fatalf("frame=%d bad expectation:\n%s\n%s", frame, x, marshalPattern(pixels))
		}
	}
}

func testFrame(t *testing.T, p Pattern, e expectation) {
	pixels := make(Frame, len(e.colors))
	p.Render(pixels, e.offsetMS)
	if !e.colors.IsEqual(pixels) {
		t.Fatalf("%s != %s", marshalPattern(e.colors), marshalPattern(pixels))
	}
}

/*
func frameSimilar(lhs, rhs Frame) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, a := range lhs {
		b := rhs[i]
		dR := int(a.R) - int(b.R)
		dG := int(a.G) - int(b.G)
		dB := int(a.B) - int(b.B)
		if dR > 1 || dR < -1 || dG > 1 || dG < -1 || dB > 1 || dB < -1 {
			return false
		}
	}
	return true
}
*/
