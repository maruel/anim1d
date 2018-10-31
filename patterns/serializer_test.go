// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package patterns

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/maruel/anim1d"
)

func TestJSONPatternsSpotCheck(t *testing.T) {
	serializePattern(t, &Rainbow{}, `"Rainbow"`)
	serializePattern(t, &PingPong{}, `{"Child":{},"MovePerHour":0,"_type":"PingPong"}`)
	serializePattern(t, &Chronometer{}, `{"Child":{},"_type":"Chronometer"}`)

	// Create one more complex. Assert that int64 is not mangled.
	p := &Transition{
		Before: anim1d.SPattern{
			&Transition{
				After:        anim1d.SPattern{&anim1d.Color{255, 255, 255}},
				OffsetMS:     600000,
				TransitionMS: 600000,
				Curve:        anim1d.Direct,
			},
		},
		After:        anim1d.SPattern{&anim1d.Color{}},
		OffsetMS:     30 * 60000,
		TransitionMS: 600000,
		Curve:        anim1d.Direct,
	}
	expected := `{"After":"#000000","Before":{"After":"#ffffff","Before":{},"Curve":"direct","OffsetMS":600000,"TransitionMS":600000,"_type":"Transition"},"Curve":"direct","OffsetMS":1800000,"TransitionMS":600000,"_type":"Transition"}`
	serializePattern(t, p, expected)
}

//

func serializePattern(t *testing.T, p anim1d.Pattern, expected string) {
	p2 := anim1d.SPattern{p}
	b, err := json.Marshal(&p2)
	if err != nil {
		t.Fatalf("%s: %v", reflect.TypeOf(p), err)
	}
	if s := string(b); s != expected {
		t.Fatalf("%s != %s", s, expected)
	}
	p2.Pattern = nil
	if err = json.Unmarshal(b, &p2); err != nil {
		t.Fatal(err)
	}
}
