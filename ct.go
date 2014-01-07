// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ct provides ANSI terminal text coloring and decoration support.
package ct

import (
	"bytes"
	"fmt"
	"reflect"
)

// Color is a basic color list selector. Not all colors are available on all platforms.
type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	LightGray
	Gray
	BoldRed
	BoldGreen
	BoldYellow
	BoldBlue
	BoldMagenta
	BoldCyan
	White
)

// Fg returns a foreground color Mode based on the provided Color.
func Fg(c Color) Mode { return colorSet | Mode(c)&colorMask }

// Bg returns a background color Mode based on the provided Color.
func Bg(c Color) Mode { return (colorSet | Mode(c)&colorMask) << colorWidth }

// XTermFg returns an XTerm foreground color Mode based on the provided XTerm color.
func XTermFg(c byte) Mode { return (xTermColorSet | Mode(c)) << (2 * colorWidth) }

// XTermBg returns an XTerm background color Mode based on the provided XTerm color.
func XTermBg(c byte) Mode { return (xTermColorSet | Mode(c)) << (2*colorWidth + xTermColorWidth) }

const (
	colorWidth = 5
	colorSet   = 1 << (colorWidth - 1)
	colorMask  = colorSet - 1

	xTermColorWidth = 9
	xTermColorSet   = 1 << (xTermColorWidth - 1)
	xTermColorMask  = xTermColorSet - 1

	colorBits = 2*colorWidth + 2*xTermColorWidth // Fg + Bg + XtermFg + XtermBg

	activeBits = ^Mode(0)&^(1<<colorBits-1)&^NoResetAfter |
		colorSet | colorSet<<colorWidth | xTermColorSet<<(2*colorWidth) | xTermColorSet<<(2*colorWidth+xTermColorWidth)
)

// Mode specifies terminal rendering modes via its Render method. Modes should be bitwise or'd together.
// Not all modes are available on all platforms.
//
// A Mode is essentially a bit field of colors and decoration flags. The layout of a Mode is as follows:
//
//  bits 0-3    - color of foreground
//  bit 4       - set foreground color
//  bits 5-7    - color of background
//  bit 8       - set background color
//  bits 9-16   - XTerm color of foreground
//  bit 17      - set XTerm foreground color
//  bits 18-25  - XTerm color of background
//  bit 26      - set XTerm background color
//  bits 27-    - text decoration flags
//  highest bit - do not reset terminal after rendering
//
// XTerm colors take priority over normal terminal colors if supported by the platform.
type Mode uint64

const (
	Reset Mode = 1 << (iota + colorBits)
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	Negative
	Conceal
	CrossedOut
	NoResetAfter Mode = 1 << 63
)

// This set of constants must be kept in sync with the Mode constants above.
const (
	reset = iota
	bold
	faint
	italic
	underline
	blinkSlow
	blinkRapid
	negative
	conceal
	crossedOut
)

// Paint returns a fmt.Formatter that will apply the mode to the
// printed values of the parameters.
func (m Mode) Paint(v ...interface{}) fmt.Formatter {
	return text{Mode: m, v: v}
}

type text struct {
	Mode
	v []interface{}
}

func doesString(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	case fmt.Stringer:
		return true
	case fmt.GoStringer:
		return true
	case fmt.Formatter:
		return true
	default:
		return reflect.ValueOf(v).Kind() == reflect.String
	}
}

// Format allows text to satisfy the fmt.Formatter interface. The format
// behaviour is the same as for fmt.Print.
func (t text) Format(fs fmt.State, c rune) {
	if t.Mode&activeBits != 0 {
		t.Mode.set(fs)
	}

	w, wOk := fs.Width()
	p, pOk := fs.Precision()
	var (
		b          bytes.Buffer
		prevString bool
	)
	b.WriteByte('%')
	for _, f := range "+-# 0" {
		if fs.Flag(int(f)) {
			b.WriteRune(f)
		}
	}
	if wOk {
		fmt.Fprint(&b, w)
	}
	if pOk {
		b.WriteByte('.')
		fmt.Fprint(&b, p)
	}
	b.WriteRune(c)
	format := b.String()

	for _, v := range t.v {
		isString := v != nil && doesString(v)
		if isString && prevString {
			fs.Write([]byte{' '})
		}
		prevString = isString
		fmt.Fprintf(fs, format, v)
	}

	if t.Mode&activeBits != 0 && t.Mode&activeBits != Reset && t.Mode&NoResetAfter == 0 {
		t.reset(fs)
	}
}
