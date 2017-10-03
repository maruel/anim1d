// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

// Value defines a value that may be constant or that may evolve over time.
type Value interface {
	Eval(timeMS uint32, l int) int32
}
