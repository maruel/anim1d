// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestNilObject(t *testing.T) {
	c := Frame{}
	d := Frame{{}}
	for _, r := range patternsStr {
		p := reflect.New(r).Interface().(Pattern)
		p.Render(nil, 0)
		p.Render(c, 0)
		p.Render(d, 0)
	}
}

func TestSPattern_Render(t *testing.T) {
	p := SPattern{}
	p.Render(nil, 0)
}

func TestSPattern_Marshal(t *testing.T) {
	p := SPattern{}
	if b, err := p.MarshalJSON(); err != nil || string(b) != "{}" {
		t.Fatal(b, err)
	}
	p.Pattern = &invalid{}
	if _, err := p.MarshalJSON(); err == nil {
		t.Fatal("expected failure")
	}
}

func TestSPattern_Umarshal(t *testing.T) {
	p := SPattern{}
	if err := p.UnmarshalJSON(nil); err == nil || p.Pattern != nil {
		t.Fatal("expected failure")
	}
	if err := p.UnmarshalJSON([]byte("\"#\"")); err == nil || p.Pattern != nil {
		t.Fatal(p, err)
	}
	if err := p.UnmarshalJSON([]byte("\"#102030\"")); err != nil || p.Pattern == nil {
		t.Fatal(p, err)
	}
}

func TestSPattern_Umarshal_Dict(t *testing.T) {
	defer func() {
		patternsDictLookup = map[string]reflect.Type{}
	}()

	if err := RegisterPattern(&pat{Dict}); err != nil {
		t.Fatal(err)
	}
	/*
		if err := p.UnmarshalJSON([]byte(`{"_type":`)); err != nil || p.Pattern == nil {
			t.Fatal(p, err)
		}
	*/
}

func TestJSONPatterns(t *testing.T) {
	for _, r := range patternsStr {
		p := reflect.New(r).Interface().(Pattern)
		p2 := &SPattern{p}
		b, err := json.Marshal(p2)
		if err != nil {
			t.Fatal(err)
		}
		if c := b[0]; c != uint8('"') {
			t.Fatalf("Expected '\"', got %q", c)
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
}

func TestRegisterPattern(t *testing.T) {
	org := patternsStr
	patternsStr = nil
	defer func() {
		patternsStr = org
		patternsDictLookup = map[string]reflect.Type{}
	}()
	if RegisterPattern(&pat{Int}) == nil {
		t.Fatal("Int is not supported")
	}
	if err := RegisterPattern(&pat{Str}); err != nil {
		t.Fatal(err)
	}
	if err := RegisterPattern(&pat{Dict}); err != nil {
		t.Fatal(err)
	}
	/*
		if RegisterPattern(&invalid{}) == nil {
			t.Fatal("doesn't implement JsonEncoding")
		}
	*/
}

func TestColor_JSON(t *testing.T) {
	c := Color{}
	if c.UnmarshalJSON([]byte(nil)) == nil {
		t.Fatal("expected error")
	}
	if err := c.UnmarshalJSON([]byte("\"#102030\"")); err != nil {
		t.Fatal(err)
	}
}

func TestFrame_JSON(t *testing.T) {
	f := Frame{}
	if f.UnmarshalJSON([]byte(nil)) == nil {
		t.Fatal("expected error")
	}
}

//

type pat struct {
	t Type
}

func (p *pat) Type() Type                    { return p.t }
func (p *pat) Render(f Frame, timeMS uint32) {}
func (p *pat) MarshalJSON() ([]byte, error)  { return nil, errors.New("fake") }
func (p *pat) UnmarshalJSON(b []byte) error  { return errors.New("fake") }

type invalid struct{}

func (i *invalid) Render(f Frame, timeMS uint32) {}

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

// marshalPattern is a shorthand to JSON encode a pattern.
func marshalPattern(p Pattern) []byte {
	b, err := json.Marshal(&SPattern{p})
	if err != nil {
		panic(err)
	}
	return b
}
