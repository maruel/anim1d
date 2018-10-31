// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

func TestJSONValues(t *testing.T) {
	for _, v := range knownValues {
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
				t.Fatalf("Expected '{', got %q", c)
			}
		}
		v2.Value = nil
		if err = json.Unmarshal(b, v2); err != nil {
			t.Fatalf("%q: %v", b, err)
		}
	}
}

func TestJSONValuesSpotCheck(t *testing.T) {
	// Increase coverage of edge cases.
	c := Const(10)
	serializeValue(t, &c, `10`)
	p := Percent(65536)
	serializeValue(t, &p, `"100%"`)
	p = 6554
	serializeValue(t, &p, `"10%"`)
	p = 6553
	serializeValue(t, &p, `"9.999%"`)
	p = -6554
	serializeValue(t, &p, `"-10%"`)
	serializeValue(t, &Rand{}, `"rand"`)
	serializeValue(t, &Rand{TickMS: 43}, `{"TickMS":43,"_type":"Rand"}`)
}

//

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
