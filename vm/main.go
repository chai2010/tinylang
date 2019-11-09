// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//#include "./comet.h"
import "C"

import (
	"os"
)

func main() {
	argc := C.int(len(os.Args))
	argv := make([]*C.char, len(os.Args))

	for i, s := range os.Args {
		argv[i] = C.CString(s)
	}

	C.cometMain(argc, &argv[0])
}
