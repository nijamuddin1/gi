// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windriver

import (
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/driver/internal/errapp"
	"github.com/goki/gi/oswin/driver/internal/win32"
)

// Main is called by the program's main function to run the graphical
// application.
//
// It calls f on the App, possibly in a separate goroutine, as some OS-
// specific libraries require being on 'the main thread'. It returns when f
// returns.
func Main(f func(oswin.App)) {
	oswin.TheApp = theApp
	theApp.initScreens()
	if err := win32.Main(func() { f(theApp) }); err != nil {
		f(errapp.Stub(err))
	}
}
