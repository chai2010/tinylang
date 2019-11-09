// +build ignore

package main

import (
	"github.com/chai2010/tinylang/pkg/spec/comet"
)

func main() {
	vm := comet.New(nil, nil, 0)
	vm.DebugRun()
}
