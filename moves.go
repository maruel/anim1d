package anim1d

import "github.com/maruel/anim1d/math32"

// MovePerHour is the number of movement per hour.
//
// Can be either positive or negative. Maximum supported value is ±3600000, 1000
// move/sec.
//
// Sample values:
//   - 1: one move per hour
//   - 60: one move per minute
//   - 3600: one move per second
//   - 216000: 60 move per second
type MovePerHour SValue

// Eval is not a Value implementation but it leverages an inner one.
func (m *MovePerHour) Eval(timeMS uint32, l int, cycle int) int {
	s := SValue(*m)
	// Prevent overflows.
	v := math32.MinMax32(s.Eval(timeMS, l), -3600000, 3600000)
	// TODO(maruel): Reduce the amount of int64 code in there yet keeping it from
	// overflowing.
	// offset ranges [0, 3599999]
	offset := timeMS % 3600000
	// (1<<32)/3600000 = 1193 is too low. Temporarily upgrade to int64 to
	// calculate the value.
	low := int64(offset) * int64(v) / 3600000
	hour := timeMS / 3600000
	high := int64(hour) * int64(v)
	if cycle != 0 {
		return int((low + high) % int64(cycle))
	}
	return int(low + high)
}
