// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows ansi

package ct

import (
	"fmt"
	"testing"
)

type S string

func (s S) String() string {
	return "{" + string(s) + "}"
}

func par(p ...interface{}) []interface{} { return p }

var testModesToANSI = []struct {
	mode   Mode
	params []interface{}
	want   string
}{
	{0, par("Hello, 世界"), "Hello, 世界"},
	{Reset, par("Hello, 世界"), "\x1b[0mHello, 世界"},
	{NoResetAfter, par("Hello, 世界"), "Hello, 世界"},
	{0xf, par("Hello, 世界"), "Hello, 世界"},
	{0xf << colorWidth, par("Hello, 世界"), "Hello, 世界"},
	{0xff << (2 * colorWidth), par("Hello, 世界"), "Hello, 世界"},
	{0xff << (2*colorWidth + xTermColorWidth), par("Hello, 世界"), "Hello, 世界"},
	{Fg(White) | Bg(Red), par("Hello, 世界"), "\x1b[37;41;1mHello, 世界\x1b[0m"},
	{Fg(White) | Bg(Red), par("simple text"), "\x1b[37;41;1msimple text\x1b[0m"},
	{Fg(White) | Bg(Red) | NoResetAfter, par("don't reset"), "\x1b[37;41;1mdon't reset"},
	{Fg(White) | Bg(Red), par(Fg(White) | Bg(Red)), csi + "37;41;1m" + fmt.Sprint(Fg(White)|Bg(Red)) + resetSgr},
	{XTermFg(21) | XTermBg(196) | Underline | Bold, par("XTerm colors, underline and bold"), "\x1b[38;5;21;48;5;196;1;4mXTerm colors, underline and bold\x1b[0m"},
	{Fg(White) | Bg(Red) | Bold, par(S("fmt.Stringer")), "\x1b[37;41;1m{fmt.Stringer}\x1b[0m"},
	{Fg(White) | Bg(Red) | Bold, par(S("fmt.Stringer"), "simple text", "and", "the answer is:", 4, 2), "\x1b[37;41;1m{fmt.Stringer} simple text and the answer is:42\x1b[0m"},
}

func TestRenderToANSI(t *testing.T) {
	for _, tt := range testModesToANSI {
		got := fmt.Sprint(tt.mode.Paint(tt.params...))
		if got != tt.want {
			t.Errorf("Render to ANSI got: %q want: %q", got, tt.want)
		}
	}
}
