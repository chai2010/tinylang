// Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (p *Compiler) getGlobal(name string) *ir.Global {
	for _, g := range p.module.Globals {
		if g.Name() == name {
			return g
		}
	}
	return nil
}

func (p *Compiler) defGlobal(name string, init constant.Constant) *ir.Global {
	if g := p.getGlobal(name); g != nil {
		return g
	}
	g := p.module.NewGlobalDef(name, init)
	return g
}

func (p *Compiler) getFunc(name string) *ir.Func {
	for _, fn := range p.module.Funcs {
		if fn.Name() == name {
			return fn
		}
	}
	return nil
}

func (p *Compiler) defFunc(name string, retType types.Type, params ...*ir.Param) *ir.Func {
	if fn := p.getFunc(name); fn != nil {
		return fn
	}
	fn := p.module.NewFunc(name, retType, params...)
	return fn
}
