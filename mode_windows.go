// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !ansi

package ct

import (
	"io"
)

func (m Mode) set(w io.Writer) {}

func (m Mode) reset(w io.Writer) {}
