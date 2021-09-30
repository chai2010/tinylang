// Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler

import (
	"bytes"
	"fmt"
	gotoken "go/token"

	"github.com/chai2010/tinylang/tiny/ast"
	"github.com/chai2010/tinylang/tiny/token"
	"github.com/chai2010/tinylang/wasm/pkg/wasm/encoding"
	"github.com/chai2010/tinylang/wasm/pkg/wasm/instruction"
	"github.com/chai2010/tinylang/wasm/pkg/wasm/module"
	"github.com/chai2010/tinylang/wasm/pkg/wasm/types"
)

type Compiler struct {
	stages []func() error

	file *ast.File
	fset *gotoken.FileSet

	globalNames []string
	globalNodes map[string]*ast.Ident

	module   *module.Module
	mainBody *module.CodeEntry

	tinyMainIndex  uint32
	tinyReadIndex  uint32
	tinyWriteIndex uint32
}

func NewCompiler() *Compiler {
	p := new(Compiler)
	p.stages = []func() error{
		p.reset,
		p.initModule,
		p.emitGlobals,
		p.emitMainBody,
	}
	return p
}

// Compile 编译 TINY 程序对应的语法树
func (p *Compiler) Compile(f *ast.File) (*module.Module, error) {
	fset := gotoken.NewFileSet()
	if len(f.Data) > 0 {
		fset.AddFile(f.Name, 1, len(f.Data)).SetLinesForContent(f.Data)
	}

	p.file = f
	p.fset = fset

	for _, stage := range p.stages {
		if err := stage(); err != nil {
			return nil, err
		}
	}

	return p.module, nil
}

func (p *Compiler) ModuleString() string {
	if p.module == nil {
		return ""
	}
	var buf bytes.Buffer
	module.Pretty(&buf, p.module, module.PrettyOption{Contents: true})
	return buf.String()
}

func (p *Compiler) reset() error {
	p.globalNames = []string{}
	p.globalNodes = make(map[string]*ast.Ident)
	return nil
}

func (p *Compiler) initModule() error {
	p.module = &module.Module{
		Names: module.NameSection{
			Module: "tinylang",
		},
		Memory: module.MemorySection{
			Memorys: []module.Memory{
				{InitPages: 1, MaxPages: 1},
			},
		},
		Export: module.ExportSection{
			Exports: []module.Export{
				{
					Name: "memory",
					Descriptor: module.ExportDescriptor{
						Type:  module.MemoryExportType,
						Index: 0,
					},
				},
			},
		},
	}

	// init types
	p.tinyReadIndex = 0
	p.tinyWriteIndex = 1
	p.tinyMainIndex = 2

	p.module.Type.Functions = []module.FunctionType{
		// func __tiny_read() int32
		{Results: []types.ValueType{types.I32}},
		// func __tiny_write(x int32)
		{Params: []types.ValueType{types.I32}},
		// func _start
		{},
	}

	// import
	p.module.Import.Imports = []module.Import{
		{
			Module: "env",
			Name:   "__tiny_read",
			Descriptor: module.FunctionImport{
				Func: p.tinyReadIndex,
			},
		},
		{
			Module: "env",
			Name:   "__tiny_write",
			Descriptor: module.FunctionImport{
				Func: p.tinyWriteIndex,
			},
		},
	}

	// _start func
	p.module.Function.TypeIndices = []uint32{
		p.tinyMainIndex,
	}
	p.module.Names.Functions = append(p.module.Names.Functions, module.NameMap{
		Index: p.tinyMainIndex,
		Name:  "_start",
	})

	// _start func body
	{
		var entry = &module.CodeEntry{
			Func: module.Function{
				Locals: []module.LocalDeclaration{},
				Expr: module.Expr{
					Instrs: []instruction.Instruction{
						instruction.I32Const{Value: 42},
						instruction.Call{Index: p.tinyWriteIndex},
						instruction.Return{},
					},
				},
			},
		}

		var buf bytes.Buffer
		if err := encoding.WriteCodeEntry(&buf, entry); err != nil {
			return err
		}

		p.module.Code.Segments = append(p.module.Code.Segments, module.RawCodeSegment{
			Code: buf.Bytes(),
		})
	}

	// export _start
	p.module.Export.Exports = append(p.module.Export.Exports, module.Export{
		Name: "_start",
		Descriptor: module.ExportDescriptor{
			Type:  module.FunctionExportType,
			Index: p.tinyMainIndex,
		},
	})

	return nil
}

func (p *Compiler) emitGlobals() error {
	if len(p.module.Global.Globals) > 0 {
		return nil
	}

	ast.Walk(p.file, func(node ast.Node) bool {
		if n, ok := node.(*ast.Ident); ok {
			if _, ok := p.globalNodes[n.Name]; !ok {
				p.globalNames = append(p.globalNames, n.Name)
				p.globalNodes[n.Name] = n
			}
		}
		return true
	})

	// default global
	defauleGlobal := module.Global{
		Type: types.I32,
		Init: module.Expr{
			Instrs: []instruction.Instruction{
				instruction.I32Const{Value: 0},
			},
		},
	}

	// global[0] is nil
	p.module.Global.Globals = []module.Global{
		defauleGlobal,
	}

	// tiny var index start from 1
	defauleGlobal.Mutable = true
	for range p.globalNames {
		p.module.Global.Globals = append(p.module.Global.Globals, defauleGlobal)
	}

	return nil
}

func (p *Compiler) emitMainBody() error {
	p.mainBody = &module.CodeEntry{
		Func: module.Function{
			Locals: []module.LocalDeclaration{},
			Expr: module.Expr{
				Instrs: []instruction.Instruction{
					// instruction.I32Const{Value: 42 + 1},
					// instruction.Call{Index: p.tinyWriteIndex},
				},
			},
		},
	}

	if err := p.compileNone(&p.mainBody.Func.Expr, p.file); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := encoding.WriteCodeEntry(&buf, p.mainBody); err != nil {
		return err
	}

	// replace main body
	p.module.Code.Segments = []module.RawCodeSegment{
		{Code: buf.Bytes()},
	}

	return nil
}

func (p *Compiler) compileNone(ctx *module.Expr, node ast.Node) error {
	if node == nil {
		return nil
	}
	switch node := node.(type) {
	case *ast.File:
		for _, n := range node.List {
			if err := p.compileNone(ctx, n); err != nil {
				return err
			}
		}
	case *ast.BlockStmt:
		for _, n := range node.List {
			if err := p.compileNone(ctx, n); err != nil {
				return err
			}
		}
	case *ast.IfStmt:
		if err := p.compileNone(ctx, node.Cond); err != nil {
			return err
		}

		var ctxIf module.Expr
		for _, n := range node.Body.List {
			if err := p.compileNone(&ctxIf, n); err != nil {
				return err
			}
		}
		if node.Else != nil {
			p.emit(&ctxIf, instruction.Else{})
			if err := p.compileNone(&ctxIf, node.Else); err != nil {
				return err
			}
		}
		p.emit(ctx, instruction.If{Instrs: ctxIf.Instrs})

	case *ast.RepeatStmt:
		var ctxLoop module.Expr
		for _, n := range node.Body.List {
			if err := p.compileNone(&ctxLoop, n); err != nil {
				return err
			}
		}
		if err := p.compileNone(&ctxLoop, node.Until); err != nil {
			return err
		}

		// if !cond { BrIf: continue }
		p.emit(&ctxLoop, instruction.I32Eqz{})
		p.emit(&ctxLoop, instruction.BrIf{})

		p.emit(ctx, instruction.Loop{Instrs: ctxLoop.Instrs})

	case *ast.AssignStmt:
		if err := p.compileNone(ctx, node.Value); err != nil {
			return err
		}
		p.emit(ctx, instruction.SetGlobal{Index: p.globalIndexByName(node.Target.Name)})
	case *ast.ReadStmt:
		p.emit(ctx, instruction.Call{Index: p.tinyReadIndex})
		p.emit(ctx, instruction.SetGlobal{Index: p.globalIndexByName(node.Target.Name)})
	case *ast.WriteStmt:
		if err := p.compileNone(ctx, node.Value); err != nil {
			return err
		}
		p.emit(ctx, instruction.Call{Index: p.tinyWriteIndex})
	case *ast.Ident:
		p.emit(ctx, instruction.GetGlobal{Index: p.globalIndexByName(node.Name)})
	case *ast.Number:
		p.emit(ctx, instruction.I32Const{Value: int32(node.Value)})
	case *ast.ParenExpr:
		if err := p.compileNone(ctx, node.X); err != nil {
			return err
		}
	case *ast.BinaryExpr:
		if err := p.compileNone(ctx, node.X); err != nil {
			return err
		}
		if err := p.compileNone(ctx, node.Y); err != nil {
			return err
		}
		switch node.Op {
		case token.LT:
			p.emit(ctx, instruction.I32LtS{})
		case token.EQ:
			p.emit(ctx, instruction.I32Eq{})
		case token.PLUS:
			p.emit(ctx, instruction.I32Add{})
		case token.MINUS:
			p.emit(ctx, instruction.I32Sub{})
		case token.TIMES:
			p.emit(ctx, instruction.I32Mul{})
		case token.OVER:
			p.emit(ctx, instruction.I32DivS{})
		default:
			panic(fmt.Sprintf("unreachable: op = %v", node.Op))
		}
	default:
		panic(fmt.Sprintf("unreachable: type = %T", node))
	}

	return nil
}

func (p *Compiler) globalIndexByName(name string) uint32 {
	for i, s := range p.globalNames {
		if s == name {
			return uint32(i) + 1
		}
	}
	return 0
}

func (p *Compiler) emit(ctx *module.Expr, x ...instruction.Instruction) {
	ctx.Instrs = append(ctx.Instrs, x...)
}
