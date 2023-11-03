// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// anim1d renders an anim1d animation to the terminal or LEDs strip.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/maruel/anim1d"
	"github.com/maruel/ansi256"
	"periph.io/x/conn/v3/display"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/apa102"
	"periph.io/x/devices/v3/screen1d"
	"periph.io/x/host/v3"
)

func mainImpl() error {
	verbose := flag.Bool("v", false, "verbose mode")
	fake := flag.Bool("terminal", false, "print the animation at the terminal")
	spiID := flag.String("spi", "", "SPI port to use")
	hz := flag.Int("hz", 0, "SPI port speed")
	numPixels := flag.Int("n", apa102.DefaultOpts.NumPixels, "number of pixels on the strip")
	intensity := flag.Int("l", int(apa102.DefaultOpts.Intensity), "light intensity [1-255]")
	temperature := flag.Int("t", int(apa102.DefaultOpts.Temperature), "light temperature in °Kelvin [3500-7500]")
	fps := flag.Int("fps", 30, "frames per second")
	fileName := flag.String("f", "", "file to load the animation from")
	raw := flag.String("r", "", "inline serialized animation")
	flag.Parse()
	if !*verbose {
		log.SetOutput(io.Discard)
	}
	log.SetFlags(log.Lmicroseconds)
	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}
	if *intensity < 1 || *intensity > 255 {
		return errors.New("intensity must be between 1 and 255")
	}
	if *temperature < 0 || *temperature > 65535 {
		return errors.New("temperature must be between 0 and 65535")
	}
	if *numPixels < 1 || *numPixels > 10000 {
		return errors.New("number of pixels must be between 1 and 10000")
	}
	if *fps < 1 || *fps > 200 {
		return errors.New("fps must be between 1 and 200")
	}
	var pat anim1d.SPattern
	if *fileName != "" {
		if *raw != "" {
			return errors.New("can't use both -f and -r")
		}
		c, err := os.ReadFile(*fileName)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(c, &pat); err != nil {
			return fmt.Errorf("bad pattern: %w", err)
		}
	} else if *raw != "" {
		if err := json.Unmarshal([]byte(*raw), &pat); err != nil {
			return fmt.Errorf("bad pattern: %w", err)
		}
	} else {
		return errors.New("use one of -f or -r; try -r '\"#0101ff\"'")
	}

	var display displayWriter
	if *fake {
		// intensity and temperature are ignored.
		display = screen1d.New(&screen1d.Opts{X: *numPixels, Palette: ansi256.Default})
	} else {
		if _, err := host.Init(); err != nil {
			return err
		}
		s, err := spireg.Open(*spiID)
		if err != nil {
			if *spiID == "" {
				return fmt.Errorf("use -terminal if you don't have LEDs; error opening SPI: %v", err)
			}
			return err
		}
		defer s.Close()
		if *hz != 0 {
			if err = s.LimitSpeed(physic.Frequency(*hz) * physic.Hertz); err != nil {
				return err
			}
		}
		opts := apa102.DefaultOpts
		opts.NumPixels = *numPixels
		opts.Intensity = uint8(*intensity)
		opts.Temperature = uint16(*temperature)
		if display, err = apa102.New(s, &opts); err != nil {
			return err
		}
	}
	// TODO(maruel): Handle Ctrl-C to cleanly exit.
	defer display.Halt()
	return runLoop(display, pat.Pattern, *fps)
}

type displayWriter interface {
	display.Drawer
	io.Writer
}

func runLoop(display displayWriter, p anim1d.Pattern, fps int) error {
	// TODO(maruel): Use double-buffering: one goroutine generates the frames,
	// the other transmits the data.
	delta := time.Second / time.Duration(fps)
	numLights := display.Bounds().Dx()
	buf := make([]byte, numLights*3)
	f := make(anim1d.Frame, numLights)
	t := time.NewTicker(delta)
	start := time.Now()
	for {
		// Wraps after 49.71 days.
		p.Render(f, uint32(time.Since(start)/time.Millisecond))
		f.ToRGB(buf)
		if _, err := display.Write(buf); err != nil {
			return err
		}
		<-t.C
	}
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "anim1d: %s.\n", err)
		os.Exit(1)
	}
}
