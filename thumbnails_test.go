// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"testing"
)

func TestThumbnailsCache(t *testing.T) {
	tb := ThumbnailsCache{
		NumberLEDs:       10,
		ThumbnailHz:      10,
		ThumbnailSeconds: 1,
	}
	c := Color{0x10, 0x20, 0x03}
	b, err := c.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	g1, err := tb.GIF(b)
	if err != nil {
		t.Fatal(err)
	}
	g2, err := tb.GIF(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(g1, g2) {
		t.Fatal("different encoded bytes")
	}

	if _, err = tb.GIF([]byte("foo")); err == nil {
		t.Fatal("expected encoding error")
	}
}
