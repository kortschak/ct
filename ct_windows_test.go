// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !ansi

package ct

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"unsafe"
)

type handle uintptr

type mock struct {
	hist []interface{}
	tab  map[uintptr]*info
}

func (m *mock) Call(_ ...uintptr) (r1, r2 uintptr, lastErr error) {
	h := uintptr(rand.Int63())
	m.tab[h] = &info{}
	m.hist = m.hist[:0]
	return h, 0, nil
}

func (m *mock) Write(p []byte) (int, error) { m.hist = append(m.hist, string(p)); return len(p), nil }

type bufferInfo mock

func (bi *bufferInfo) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	i := (*info)(unsafe.Pointer(a[1]))
	i = bi.tab[a[0]]
	_ = i
	return 1, 0, nil
}

type textAttribute mock

func (ta *textAttribute) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	ta.tab[a[0]].wAttributes = uint16(a[1])
	ta.hist = append(ta.hist, uint16(a[1]))
	return 1, 0, nil
}

var mockConsole = &mock{tab: make(map[uintptr]*info)}

func init() {
	getStdHandle = mockConsole
	getConsoleScreenBufferInfo = (*bufferInfo)(mockConsole)
	setConsoleTextAttribute = (*textAttribute)(mockConsole)
}

type S string

func (s S) String() string {
	return "{" + string(s) + "}"
}

func par(p ...interface{}) []interface{} { return p }

var testModesToWindowsConsole = []struct {
	mode   Mode
	params []interface{}
	want   []interface{}
}{
	{0, par("Hello, 世界"), par("Hello, 世界")},
	{Reset, par("Hello, 世界"), par(uint16(0), "Hello, 世界", uint16(0))},
	{NoResetAfter, par("Hello, 世界"), par("Hello, 世界")},
	{0xf, par("Hello, 世界"), par("Hello, 世界")},
	{0xf << colorWidth, par("Hello, 世界"), par("Hello, 世界")},
	{0xff << (2 * colorWidth), par("Hello, 世界"), par("Hello, 世界")},
	{0xff << (2*colorWidth + xTermColorWidth), par("Hello, 世界"), par("Hello, 世界")},
	{Fg(White) | Bg(Red), par("Hello, 世界"), par(uint16(0x001f), "Hello, 世界", uint16(0))},
	{Fg(White) | Bg(Red), par("simple text"), par(uint16(0x001f), "simple text", uint16(0))},
	{Fg(White) | Bg(Red) | NoResetAfter, par("don't reset"), par(uint16(0x001f), "don't reset")},
	{Fg(White) | Bg(Red), par(Fg(White) | Bg(Red)), par(uint16(0x001f), fmt.Sprint(Fg(White)|Bg(Red)), uint16(0))},
	{XTermFg(21) | XTermBg(196) | Underline | Bold, par("XTerm colors, underline and bold"), par(uint16(0x8010), "XTerm colors, underline and bold", uint16(0))},
	{Fg(White) | Bg(Red) | Bold, par(S("fmt.Stringer")), par(uint16(0x001f), "{fmt.Stringer}", uint16(0))},
	{Fg(White) | Bg(Red) | Bold, par(S("fmt.Stringer"), "simple text", "and", "the answer is:", 4, 2), par(uint16(0x001f), "{fmt.Stringer} simple text and the answer is:42", uint16(0))},
}

func TestRenderToWindowsConsole(t *testing.T) {
	for _, tt := range testModesToWindowsConsole {
		fmt.Fprint(mockConsole, tt.mode.Paint(tt.params...))
		if !reflect.DeepEqual(mockConsole.hist, tt.want) {
			t.Errorf("Render to Windows got: %+v want: %+v", mockConsole.hist, tt.want)
		}
		mockConsole.hist = mockConsole.hist[:0]
	}
}
