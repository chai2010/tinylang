// Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/wasmerio/wasmer-go/wasmer"
)

func Run(wasmBytes []byte) error {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		return err
	}

	tinyRead := wasmer.NewFunction(store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(), wasmer.NewValueTypes(wasmer.I32)),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			var x int32
			fmt.Print("READ: ")
			fmt.Scanf("%d", &x)
			return []wasmer.Value{wasmer.NewI32(x)}, nil
		},
	)
	tinyWrite := wasmer.NewFunction(store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			if len(args) > 0 {
				fmt.Println(args[0].I32())
			}
			return nil, nil
		},
	)

	importObject := wasmer.NewImportObject()
	importObject.Register("env", map[string]wasmer.IntoExtern{
		"__tiny_read":  tinyRead,
		"__tiny_write": tinyWrite,
	})

	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return err
	}

	_start_fn, err := instance.Exports.GetFunction("_start")
	if err != nil {
		return err
	}

	if _, err = _start_fn(); err != nil {
		return err
	}

	return nil
}
