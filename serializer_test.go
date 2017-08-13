// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestNilObject(t *testing.T) {
	c := Frame{}
	d := Frame{{}}
	for _, p := range knownPatterns {
		p.Render(nil, 0)
		p.Render(c, 0)
		p.Render(d, 0)
	}
}

func TestJSONPatterns(t *testing.T) {
	for _, p := range knownPatterns {
		p2 := &SPattern{p}
		b, err := json.Marshal(p2)
		if err != nil {
			t.Fatal(err)
		}
		if isColorOrFrameOrRainbow(p) {
			if c := b[0]; c != uint8('"') {
				t.Fatalf("Expected '\"', got %q", c)
			}
		} else {
			if c := b[0]; c != uint8('{') {
				t.Fatalf("Expected '{', got %q", c)
			}
		}
		// Must not crash on nil members and empty frame.
		p2.Render(Frame{}, 0)
		p2.Pattern = nil
		if err := json.Unmarshal(b, p2); err != nil {
			t.Fatalf("%s, %vv", b, err)
		}
	}
}

func TestJSONPatternsSpotCheck(t *testing.T) {
	// Increase coverage of edge cases.
	serializePattern(t, &Color{1, 2, 3}, `"#010203"`)
	serializePattern(t, &Frame{}, `"L"`)
	serializePattern(t, &Frame{{1, 2, 3}, {4, 5, 6}}, `"L010203040506"`)
	serializePattern(t, &Rainbow{}, `"Rainbow"`)
	serializePattern(t, &PingPong{}, `{"Child":{},"MovePerHour":0,"_type":"PingPong"}`)
	serializePattern(t, &Chronometer{}, `{"Child":{},"_type":"Chronometer"}`)

	// Create one more complex. Assert that int64 is not mangled.
	p := &Transition{
		Before: SPattern{
			&Transition{
				After:        SPattern{&Color{255, 255, 255}},
				OffsetMS:     600000,
				TransitionMS: 600000,
				Curve:        Direct,
			},
		},
		After:        SPattern{&Color{}},
		OffsetMS:     30 * 60000,
		TransitionMS: 600000,
		Curve:        Direct,
	}
	expected := `{"After":"#000000","Before":{"After":"#ffffff","Before":{},"Curve":"direct","OffsetMS":600000,"TransitionMS":600000,"_type":"Transition"},"Curve":"direct","OffsetMS":1800000,"TransitionMS":600000,"_type":"Transition"}`
	serializePattern(t, p, expected)
}

func TestJSONValues(t *testing.T) {
	for _, v := range knownValues {
		name := reflect.TypeOf(v).Elem().Name()
		v := v
		t.Run(name, func(t *testing.T) {
			v2 := &SValue{v}
			b, err := json.Marshal(v2)
			if err != nil {
				t.Fatal(err)
			}
			if isConst(v) {
				if _, err = strconv.ParseInt(string(b), 10, 32); err != nil {
					t.Fatalf("%v", err)
				}
			} else if isPercent(v) {
				// Skip the %.
				f := 0.
				if _, err = fmt.Sscanf(string(b), "\"%g%%\"", &f); err != nil {
					t.Fatalf("%v", err)
				}
			} else if isOpAdd(v) {
				// Skip the +.
				i := 0
				if _, err = fmt.Sscanf(string(b), "\"+%d\"", &i); err != nil {
					t.Fatalf("%v", err)
				}
			} else if isOpMod(v) {
				// Skip the %.
				i := 0
				if _, err = fmt.Sscanf(string(b), "\"%%%d\"", &i); err != nil {
					t.Fatalf("%v", err)
				}
			} else if isRand(v) && string(b) == "\""+randKey+"\"" {
				// Ok.
			} else {
				if c := b[0]; c != uint8('{') {
					t.Errorf("want '{', got %q", c)
				}
			}
			/*
				if _, err := parseValueString(s); err != nil {
					t.Fatalf("%d: %v", i, err)
				}
			*/
		})
	}
}

func TestStringValues(t *testing.T) {
	//c1 := Const(1)
	//c2 := Const(2)
	//c3 := Const(3)
	c10 := Const(10)
	p10 := Percent(6554)
	pm10 := Percent(-6554)
	pm99 := Percent(-6553)
	data := []struct {
		v        Value
		expected string
	}{
		{&TimeMS{}, "t"},
		{&Length{}, "l"},
		{&c10, "10"},
		{&p10, "10%"},
		{&pm10, "-10%"},
		{&pm99, "-9.999%"},
		/*
			{&OpAdd{SValue{&c1}, SValue{&c2}}, "1+2"},
			{&OpGroup{SValue{&c1}}, ""},
			{&OpMod{SValue{&c3}, SValue{&c2}}, ""},
			{&OpMul{SValue{&c2}, SValue{&c3}}, ""},
			{&OpStep{SValue{&c2}}, ""},
			{&OpSub{SValue{&c2}, SValue{&c1}}, ""},
		*/
		{&Rand{}, "Rand"},
	}
	for i, line := range data {
		s := line.v.String()
		if line.expected != s {
			t.Fatalf("%d: %s != %s", i, line.expected, s)
		}
		if _, err := parseValueString(s); err != nil {
			t.Fatalf("%d: %v", i, err)
		}
	}
}

//

func serializePattern(t *testing.T, p Pattern, expected string) {
	p2 := &SPattern{p}
	b, err := json.Marshal(p2)
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); s != expected {
		t.Fatalf("%s != %s", s, expected)
	}
	p2.Pattern = nil
	if err = json.Unmarshal(b, p2); err != nil {
		t.Fatal(err)
	}
}

func isColorOrFrameOrRainbow(p Pattern) bool {
	if _, ok := p.(*Color); ok {
		return ok
	}
	if _, ok := p.(*Frame); ok {
		return ok
	}
	_, ok := p.(*Rainbow)
	return ok
}

// serializeValue tests tests both JSON marshalling and unmarshalling.
func serializeValue(t *testing.T, v Value, expected string) {
	v2 := &SValue{v}
	b, err := json.Marshal(v2)
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); s != expected {
		t.Fatalf("%s != %s", s, expected)
	}
	v2.Value = nil
	if err = json.Unmarshal(b, v2); err != nil {
		t.Fatal(err)
	}
}

func isTimeMS(v Value) bool {
	_, ok := v.(*TimeMS)
	return ok
}

func isLength(v Value) bool {
	_, ok := v.(*Length)
	return ok
}

func isConst(v Value) bool {
	_, ok := v.(*Const)
	return ok
}

func isPercent(v Value) bool {
	_, ok := v.(*Percent)
	return ok
}

func isOpAdd(v Value) bool {
	_, ok := v.(*OpAdd)
	return ok
}

func isOpMod(v Value) bool {
	_, ok := v.(*OpMod)
	return ok
}

func isRand(v Value) bool {
	_, ok := v.(*Rand)
	return ok
}

// marshalPattern is a shorthand to JSON encode a pattern.
func marshalPattern(p Pattern) []byte {
	b, err := json.Marshal(&SPattern{p})
	if err != nil {
		panic(err)
	}
	return b
}
