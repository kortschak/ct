// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows ansi

package ct

import (
	"fmt"
	"io"
)

var modeTable = [...]int{
	reset:      0,
	bold:       1,
	faint:      2,
	italic:     3,
	underline:  4,
	blinkSlow:  5,
	blinkRapid: 6,
	negative:   7,
	conceal:    8,
	crossedOut: 9,
}

const (
	csi      = "\x1b["
	resetSgr = csi + "0m"
)

var semicolon = []byte{';'}

func (m Mode) set(w io.Writer) {
	w.Write([]byte(csi))
	if m&(colorSet|colorMask)>>(colorWidth-2) == 0x3 {
		m |= Bold
	}
	var printed bool
	for _, d := range []Mode{30, 40} {
		if m&colorSet != 0 {
			if printed {
				w.Write(semicolon)
			}
			fmt.Fprintf(w, "%d", m&(colorMask>>1)+d)
			printed = true
		}
		m >>= colorWidth
	}
	for _, d := range []int{38, 48} {
		if m&xTermColorSet != 0 {
			if printed {
				w.Write(semicolon)
			}
			fmt.Fprintf(w, "%d;5;%d", d, m&xTermColorMask)
			printed = true
		}
		m >>= xTermColorWidth
	}
	for _, v := range modeTable {
		if m&1 != 0 {
			if printed {
				w.Write(semicolon)
			}
			fmt.Fprintf(w, "%d", v)
			printed = true
		}
		m >>= 1
	}
	w.Write([]byte{'m'})
}

func (m Mode) reset(w io.Writer) { w.Write([]byte(resetSgr)) }
