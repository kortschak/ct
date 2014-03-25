// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gopher.png Copyright ©2009 The Go Authors. All rights reserved.
// Used under the Go LICENSE available at http://golang.org/LICENSE
// Gopher artwork originally by Renée French. Used under the Creative
// Commons Attributions 3.0 license.

// +build ignore

package main

import (
	"fmt"
	"github.com/kortschak/ct"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

var (
	valueRange = [...]byte{0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
	table      = [16]color.RGBA{
		0x0: {R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		0x1: {R: 0x80, G: 0x00, B: 0x00, A: 0xff},
		0x2: {R: 0x00, G: 0x80, B: 0x00, A: 0xff},
		0x3: {R: 0x80, G: 0x80, B: 0x00, A: 0xff},
		0x4: {R: 0x00, G: 0x00, B: 0x80, A: 0xff},
		0x5: {R: 0x80, G: 0x00, B: 0x80, A: 0xff},
		0x6: {R: 0x00, G: 0x80, B: 0x80, A: 0xff},
		0x7: {R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff},
		0x8: {R: 0x80, G: 0x80, B: 0x80, A: 0xff},
		0x9: {R: 0xff, G: 0x00, B: 0x00, A: 0xff},
		0xa: {R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		0xb: {R: 0xff, G: 0xff, B: 0x00, A: 0xff},
		0xc: {R: 0x00, G: 0x00, B: 0xff, A: 0xff},
		0xd: {R: 0xff, G: 0x00, B: 0xff, A: 0xff},
		0xe: {R: 0x00, G: 0xff, B: 0xff, A: 0xff},
		0xf: {R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}
	xterm = make(color.Palette, 255)
)

type xTermColor byte

func (c xTermColor) RGBA() (r, g, b, a uint32) {
	if c < 16 {
		return table[c].RGBA()
	}
	if c < 232 {
		c -= 16
		return color.RGBA{
			R: valueRange[(c/36)%6],
			G: valueRange[(c/6)%6],
			B: valueRange[c%6],
			A: 0xff,
		}.RGBA()
	}
	cb := byte(8 + (c-232)*10)
	return color.RGBA{cb, cb, cb, 0xff}.RGBA()
}

func init() {
	for i := range xterm {
		xterm[i] = xTermColor(i)
	}
}

func main() {
	var fn = "gopher.png"
	if len(os.Args) > 1 {
		fn = os.Args[1]
	}
	f, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Couldn't find image file %q: %v\n", fn, err)
		os.Exit(1)
	}
	m, s, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Couldn't decode image file for %q: %v\n", s, err)
		os.Exit(1)
	}
	b := m.Bounds()
	for y := 0; y < b.Dy(); y++ {
		fmt.Print(" ")
		for x := 0; x < b.Dx(); x++ {
			c := m.At(x, y)
			if _, _, _, a := c.RGBA(); a < 0x80 {
				fmt.Print(" ")
				continue
			}
			fmt.Print(ct.XTermBg(byte(xterm.Index(c))).Paint(" "))
		}
		fmt.Println()
	}
}
