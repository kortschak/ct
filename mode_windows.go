// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !ansi

// BUG(kortschak): The behavior of ct on Windows platforms attempts to reasonably
// closely mimic the behavior on ANSI terminals, but because of the disjointed
// approach to console control in Windows, it is currently not possible to decorate
// output to both standard output and standard error; only standard output is
// supported.
package ct

import (
	"fmt"
	"io"
	"sync"
	"syscall"
	"unsafe"
)

const stdoutHandle = -11

func dword(i int) uintptr { return uintptr(uint32(i)) }

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	getStdHandle = kernel32.NewProc("GetStdHandle")

	getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	setConsoleTextAttribute    = kernel32.NewProc("SetConsoleTextAttribute")
)

type (
	coord     [2]int16
	shortRect [4]int16

	// CONSOLE_SCREEN_BUFFER_INFO
	info struct {
		_           coord
		_           coord
		wAttributes uint16
		_           shortRect
		_           coord
	}
)

var (
	consoleAttrs = make(map[console]uint16)
	attrLock     sync.Mutex
)

type console uintptr

func (c console) getConsoleScreenBufferInfo() *info {
	i := &info{}
	r, _, _ := getConsoleScreenBufferInfo.Call(uintptr(c), uintptr(unsafe.Pointer(i)))
	if r == 0 {
		return nil
	}
	return i
}

func (c console) setConsoleTextAttribute(attr uint16) {
	setConsoleTextAttribute.Call(uintptr(c), uintptr(attr))
}

type state struct {
	fmt.State
	console
	attr uint16
}

func hook(fs fmt.State) fmt.State {
	c, _, err := getStdHandle.Call(dword(stdoutHandle))
	if err != nil {
		panic(err)
	}

	s := state{State: fs, console: console(c)}

	i := s.getConsoleScreenBufferInfo()
	if i == nil {
		return fs
	}
	s.attr = i.wAttributes

	attrLock.Lock()
	if _, ok := consoleAttrs[s.console]; !ok {
		consoleAttrs[s.console] = i.wAttributes
	}
	attrLock.Unlock()

	return s
}

const (
	reverse    = 0x4000
	underscore = 0x8000
)

func (m Mode) set(w io.Writer) {
	if m&(colorSet|(colorSet<<colorWidth)|Reset|Bold|Negative|Underline) == 0 {
		return
	}
	s, ok := w.(state)
	if !ok {
		return
	}

	if m&Reset != 0 {
		attrLock.Lock()
		s.setConsoleTextAttribute(consoleAttrs[s.console])
		if m&activeBits == Reset {
			delete(consoleAttrs, s.console)
		}
		attrLock.Unlock()
	}

	if m&colorSet != 0 {
		s.attr &^= colorMask
		s.attr |= uint16(m & colorMask)
	}
	if m&(colorSet<<colorWidth) != 0 {
		s.attr &^= colorMask << (colorWidth - 1)
		s.attr |= uint16(m&(colorMask<<colorWidth)) >> 1
	}

	if m&Bold != 0 {
		s.attr |= 1 << (colorWidth - 1)
	}
	if m&Negative != 0 {
		s.attr |= reverse
	}
	if m&Underline != 0 {
		s.attr |= underscore
	}

	s.setConsoleTextAttribute(s.attr)
}

func (m Mode) reset(w io.Writer) {
	if m&(colorSet|(colorSet<<colorWidth)|Bold|Negative|Underline) == 0 {
		return
	}
	if s, ok := w.(state); ok {
		attrLock.Lock()
		if attr, ok := consoleAttrs[s.console]; ok {
			s.setConsoleTextAttribute(attr)
			delete(consoleAttrs, s.console)
			attrLock.Unlock()
		}
	}
}
