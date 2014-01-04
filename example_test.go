// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux unix posix ansi

package ct_test

import (
	"fmt"
	"github.com/kortschak/ct"
)

var warn = (ct.Fg(ct.White) | ct.Bg(ct.Red)).Render

func ExampleMode_Render() {
	fmt.Println(warn("WARNING:"), "Danger Will Robinson!")
}
