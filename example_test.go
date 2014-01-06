// Copyright ©2014 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ct_test

import (
	"fmt"
	"github.com/kortschak/ct"
)

var (
	info = (ct.Fg(ct.Black) | ct.XTermFg(16) | ct.Bold).Paint
	warn = (ct.Fg(ct.White) | ct.Bg(ct.Red)).Paint
)

func ExampleMode_Paint() {
	fmt.Println(warn("WARNING:"), "Danger, Will Robinson! Danger! ")
	fmt.Println(info("INFO:"), "Doctor Smith, please. You're making the Robot very unhappy!")
}
