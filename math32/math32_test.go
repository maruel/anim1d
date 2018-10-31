// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package math32

import (
	"testing"
)

func TestMinMax(t *testing.T) {
	if MinMax(2, 0, 3) != 2 {
		t.Fail()
	}
	if MinMax(-2, 0, 3) != 0 {
		t.Fail()
	}
	if MinMax(4, 0, 3) != 3 {
		t.Fail()
	}
}

func TestMinMax32(t *testing.T) {
	if MinMax32(2, 0, 3) != 2 {
		t.Fail()
	}
	if MinMax32(-2, 0, 3) != 0 {
		t.Fail()
	}
	if MinMax32(4, 0, 3) != 3 {
		t.Fail()
	}
}
