// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// mixers is all the patterns that are constructions of other patterns.

package patterns

import (
	"github.com/maruel/anim1d"
	"github.com/maruel/anim1d/math32"
)

// Gradient does a gradient between 2 patterns.
//
// A good example is using two colors but it can also be animations.
//
// TODO(maruel): Support N colors at M positions.
type Gradient struct {
	Left  anim1d.SPattern
	Right anim1d.SPattern
	Curve anim1d.Curve
	buf   anim1d.Frame
}

// Render implements Pattern.
func (g *Gradient) Render(pixels anim1d.Frame, timeMS uint32) {
	l := len(pixels)
	if l == 0 {
		return
	}
	g.buf.Reset(l)
	g.Left.Render(pixels, timeMS)
	g.Right.Render(g.buf, timeMS)
	if l == 1 {
		pixels.Mix(g.buf, g.Curve.Scale8(65535>>1))
	} else {
		max := l - 1
		for i := range pixels {
			intensity := uint16(i * 65535 / max)
			pixels[i].Mix(g.buf[i], g.Curve.Scale8(intensity))
		}
	}
}

// Split splits the strip in two.
//
// Unlike gradient, this create 2 logical independent subsets.
type Split struct {
	Left   anim1d.SPattern
	Right  anim1d.SPattern
	Offset anim1d.SValue // Point to split between both sides.
}

// Render implements Pattern.
func (s *Split) Render(pixels anim1d.Frame, timeMS uint32) {
	offset := math32.MinMax(int(s.Offset.Eval(timeMS, len(pixels))), 0, len(pixels))
	if s.Left.Pattern != nil && offset != 0 {
		s.Left.Render(pixels[:offset], timeMS)
	}
	if s.Right.Pattern != nil && offset != len(pixels) {
		s.Right.Render(pixels[offset:], timeMS)
	}
}

// Transition changes from Before to After over time. It doesn't repeat.
//
// In gets timeMS that is subtracted by OffsetMS.
type Transition struct {
	Before       anim1d.SPattern // Old pattern that is disappearing
	After        anim1d.SPattern // New pattern to show
	OffsetMS     uint32          // Offset at which the transiton from Before->In starts
	TransitionMS uint32          // Duration of the transition while both are rendered
	Curve        anim1d.Curve    // Type of transition, defaults to EaseOut if not set
	buf          anim1d.Frame
}

// Render implements Pattern.
func (t *Transition) Render(pixels anim1d.Frame, timeMS uint32) {
	if timeMS <= t.OffsetMS {
		// Before transition.
		t.Before.Render(pixels, timeMS)
		return
	}
	t.After.Render(pixels, timeMS-t.OffsetMS)
	if timeMS >= t.OffsetMS+t.TransitionMS {
		// After transition.
		t.buf = nil
		return
	}
	t.buf.Reset(len(pixels))

	// TODO(maruel): Add lateral animation and others.
	t.Before.Render(t.buf, timeMS)
	intensity := uint16((timeMS - t.OffsetMS) * 65535 / (t.TransitionMS))
	pixels.Mix(t.buf, 255.-t.Curve.Scale8(intensity))
}

// Loop rotates between all the animations.
//
// Display starts with one ShowMS for Patterns[0], then starts looping.
// timeMS is not modified so it's like as all animations continued animating
// behind.
// TODO(maruel): Add lateral transition and others.
type Loop struct {
	Patterns     []anim1d.SPattern
	ShowMS       uint32       // Duration for each pattern to be shown as pure
	TransitionMS uint32       // Duration of the transition between two patterns, can be 0
	Curve        anim1d.Curve // Type of transition, defaults to EaseOut if not set
	buf          anim1d.Frame
}

// Render implements Pattern.
func (l *Loop) Render(pixels anim1d.Frame, timeMS uint32) {
	lp := uint32(len(l.Patterns))
	if lp == 0 {
		return
	}
	cycleDuration := l.ShowMS + l.TransitionMS
	if cycleDuration == 0 {
		// Misconfigured. Lock to the first pattern.
		l.Patterns[0].Render(pixels, timeMS)
		return
	}

	base := timeMS / cycleDuration
	index := base % lp
	a := l.Patterns[index]
	a.Render(pixels, timeMS)
	offset := timeMS - (base * cycleDuration)
	if offset <= l.ShowMS {
		return
	}

	// Transition.
	l.buf.Reset(len(pixels))
	b := l.Patterns[(index+1)%lp]
	b.Render(l.buf, timeMS)
	offset -= l.ShowMS
	intensity := uint16((l.TransitionMS - offset) * 65535 / l.TransitionMS)
	pixels.Mix(l.buf, l.Curve.Scale8(65535-intensity))
}

// Rotate rotates a pattern that can also cycle either way.
//
// Use negative to go left. Can be used for 'candy bar'.
//
// Similar to PingPong{} except that it doesn't bounce.
//
// Use 5x oversampling with Scale{} to create smoother animation.
type Rotate struct {
	Child       anim1d.SPattern
	MovePerHour anim1d.MovePerHour // Expressed in number of light jumps per hour.
	buf         anim1d.Frame
}

// Render implements Pattern.
func (r *Rotate) Render(pixels anim1d.Frame, timeMS uint32) {
	l := len(pixels)
	r.buf.Reset(l)
	r.Child.Render(r.buf, timeMS)
	offset := r.MovePerHour.Eval(timeMS, len(pixels), l)
	if offset < 0 {
		// Reverse direction.
		offset = l + offset
	}
	copy(pixels[offset:], r.buf)
	copy(pixels[:offset], r.buf[l-offset:])
}

// Chronometer moves 3 lights to the right, each indicating second, minute, and
// hour passed since the start.
//
// Child has 4 pixels used in this order: [default, second, minute, hour].
type Chronometer struct {
	Child anim1d.SPattern
	buf   anim1d.Frame
}

// Render implements Pattern.
func (r *Chronometer) Render(pixels anim1d.Frame, timeMS uint32) {
	l := uint32(len(pixels))
	if l == 0 {
		return
	}
	r.buf.Reset(4)
	r.Child.Render(r.buf, timeMS)

	seconds := timeMS / 1000
	mins := seconds / 60
	hours := mins / 60

	secPos := (l*(seconds%60) + 30) / 60
	minPos := (l*(mins%60) + 30) / 60
	hourPos := hours % l

	for i := range pixels {
		switch uint32(i) {
		case secPos:
			pixels[i] = r.buf[1]
		case minPos:
			pixels[i] = r.buf[2]
		case hourPos:
			pixels[i] = r.buf[3]
		default:
			pixels[i] = r.buf[0]
		}
	}
}

// PingPong shows a 'ball' with a trail that bounces from one side to
// the other.
//
// Can be used for a ball, a water wave or K2000 (Knight Rider) style light.
// The trail can be a anim1d.Frame or a dynamic pattern.
//
// To get smoothed movement, use Scale{} with a 5x factor or so.
// TODO(maruel): That's a bit inefficient, enable Interpolation here.
type PingPong struct {
	Child       anim1d.SPattern    // [0] is the front pixel so the pixels are effectively drawn in reverse order
	MovePerHour anim1d.MovePerHour // Expressed in number of light jumps per hour
	buf         anim1d.Frame
}

// Render implements Pattern.
func (p *PingPong) Render(pixels anim1d.Frame, timeMS uint32) {
	if len(pixels) == 0 {
		return
	}
	p.buf.Reset(len(pixels)*2 - 1)
	p.Child.Render(p.buf, timeMS)
	// The last point of each extremity is only lit on one tick but every other
	// points are lit twice during a full cycle. This means the full cycle is
	// 2*(len(pixels)-1). For a 3 pixels line, the cycle is: x00, 0x0, 00x, 0x0.
	//
	// For Child being anim1d.Frame "01234567":
	//   move == 0  -> "01234567"
	//   move == 2  -> "21056789"
	//   move == 5  -> "543210ab"
	//   move == 7  -> "76543210"
	//   move == 9  -> "98765012"
	//   move == 11 -> "ba901234"
	//   move == 13 -> "d0123456"
	//   move 14 -> move 0; "2*(8-1)"
	cycle := 2 * (len(pixels) - 1)
	// TODO(maruel): Smoothing with Curve, defaults to Step.
	pos := p.MovePerHour.Eval(timeMS, len(pixels), cycle)

	// Once it works the following code looks trivial but everytime it takes me
	// an absurd amount of time to rewrite it.
	if pos >= len(pixels)-1 {
		// Head runs left.
		// pos2 is the position from the right.
		pos2 := pos + 1 - len(pixels)
		// limit is the offset at which order change.
		limit := len(pixels) - pos2 - 1
		for i := range pixels {
			if i < limit {
				// Going right.
				pixels[i] = p.buf[len(pixels)-i+pos2-1]
			} else {
				// Going left.
				pixels[i] = p.buf[i-limit]
			}
		}
	} else {
		// Head runs right.
		for i := range pixels {
			if i <= pos {
				// Going right.
				pixels[i] = p.buf[pos-i]
			} else {
				// Going left.
				pixels[i] = p.buf[pos+i]
			}
		}
	}
}

// Crop skips the beginning and the end of the source.
type Crop struct {
	Child  anim1d.SPattern
	Before anim1d.SValue // Starting pixels to skip
	After  anim1d.SValue // Ending pixels to skip
	buf    anim1d.Frame
}

// Render implements Pattern.
func (c *Crop) Render(pixels anim1d.Frame, timeMS uint32) {
	b := int(math32.MinMax32(c.Before.Eval(timeMS, len(pixels)), 0, 1000))
	a := int(math32.MinMax32(c.After.Eval(timeMS, len(pixels)), 0, 1000))
	// This is slightly wasteful as pixels are drawn just to be ditched.
	c.buf.Reset(len(pixels) + b + a)
	c.Child.Render(c.buf, timeMS)
	copy(pixels, c.buf[b:])
}

// Subset skips the beginning and the end of the destination.
type Subset struct {
	Child  anim1d.SPattern
	Offset anim1d.SValue // Starting pixels to skip
	Length anim1d.SValue // Length of the pixels to carry over
}

// Render implements Pattern.
func (s *Subset) Render(pixels anim1d.Frame, timeMS uint32) {
	if s.Child.Pattern == nil {
		return
	}
	o := math32.MinMax(int(s.Offset.Eval(timeMS, len(pixels))), 0, len(pixels)-1)
	l := math32.MinMax(int(s.Length.Eval(timeMS, len(pixels))), 0, len(pixels)-1-o)
	s.Child.Render(pixels[o:o+l], timeMS)
}

// Dim is a filter that dim the intensity of a buffer.
type Dim struct {
	Child     anim1d.SPattern //
	Intensity anim1d.SValue   // 0 is transparent, 255 is fully opaque with original colors.
}

// Render implements Pattern.
func (d *Dim) Render(pixels anim1d.Frame, timeMS uint32) {
	d.Child.Render(pixels, timeMS)
	i := math32.MinMax32(d.Intensity.Eval(timeMS, len(pixels)), 0, 255)
	pixels.Dim(uint8(i))
}

// Add is a generic mixer that merges the output from multiple patterns with
// saturation.
type Add struct {
	Patterns []anim1d.SPattern // It should be a list of Dim{} with their corresponding weight.
	buf      anim1d.Frame      //
}

// Render implements Pattern.
func (a *Add) Render(pixels anim1d.Frame, timeMS uint32) {
	a.buf.Reset(len(pixels))
	// Draw and merge each pattern.
	a.buf.reset(len(pixels))
	for i := range pixels {
		pixels[i] = Color{}
	}
	for i := range a.Patterns {
		a.Patterns[i].Render(a.buf, timeMS)
		pixels.Add(a.buf)
	}
}

// Scale adapts a larger or smaller patterns to the Strip size
//
// This is useful to create smoother horizontal movement animation or to scale
// up/down images.
type Scale struct {
	Child anim1d.SPattern
	// Defaults to Linear
	Interpolation Interpolation
	// A buffer of this len(buffer)*RatioMilli/1000 will be provided to Child and
	// will be scaled; 500 means smaller, 2000 is larger.
	//
	// Can be set to 0 when Child is a anim1d.Frame. In this case it is stretched
	// to the strip size.
	RatioMilli anim1d.SValue
	buf        anim1d.Frame
}

// Render implements Pattern.
func (s *Scale) Render(pixels anim1d.Frame, timeMS uint32) {
	if f, ok := s.Child.Pattern.(anim1d.Frame); ok {
		if s.RatioMilli.Eval(timeMS, len(pixels)) == 0 {
			s.Interpolation.Scale(f, pixels)
			return
		}
	}
	v := math32.MinMax32(s.RatioMilli.Eval(timeMS, len(pixels)), 1, 1000000)
	s.buf.Reset((int(v)*len(pixels) + 500) / 1000)
	s.Child.Render(s.buf, timeMS)
	s.Interpolation.Scale(s.buf, pixels)
}
