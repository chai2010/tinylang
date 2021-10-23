// Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	gotoken "go/token"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/chai2010/tinylang/tiny/ast"
	"github.com/chai2010/tinylang/tiny/token"
)

var (
	i32  = types.I32
	void = types.Void
)

type Options struct {
	ModulePath  string
	ModuleBytes []byte
}

// TINY 编译器(到 LLIR)
type Compiler struct {
	opt    *Options
	stages []func() error

	file *ast.File
	fset *gotoken.FileSet

	module *ir.Module

	tinyMain  *ir.Func
	tinyRead  *ir.Func
	tinyWrite *ir.Func

	globalNames []string
	globalNodes map[string]*ast.Ident

	curBlock *ir.Block
	nextId   int
}

// NewCompiler 构造编译器对象
func NewCompiler(opt *Options) *Compiler {
	p := &Compiler{opt: opt}
	p.stages = []func() error{
		p.reset,
		p.initModule,
		p.emitGlobals,
		p.emitMainBody,
	}
	return p
}

// Compile 编译 TINY 程序对应的语法树
func (p *Compiler) Compile(f *ast.File) (*ir.Module, error) {
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

// 输出 LLIR 模块文本格式
func (p *Compiler) LLString() string {
	return p.module.String()
}

func (p *Compiler) reset() error {
	p.module = nil

	p.tinyMain = nil
	p.tinyRead = nil
	p.tinyWrite = nil

	p.globalNames = []string{}
	p.globalNodes = make(map[string]*ast.Ident)

	p.nextId = 1
	return nil
}

func (p *Compiler) initModule() error {
	p.module = ir.NewModule()

	if p.opt != nil && p.opt.ModulePath != "" {
		m, err := asm.ParseBytes(p.opt.ModulePath, p.opt.ModuleBytes)
		if err != nil {
			return err
		}
		p.module = m
	}

	// i32 main()
	p.tinyMain = p.defFunc("main", i32)

	// i32 __tiny_read()
	p.tinyRead = p.defFunc("__tiny_read", i32)

	// void __tiny_write(i32 x)
	p.tinyWrite = p.defFunc("__tiny_write", void, ir.NewParam("x", i32))

	return nil
}

func (p *Compiler) emitGlobals() error {
	ast.Walk(p.file, func(node ast.Node) bool {
		if n, ok := node.(*ast.Ident); ok {
			if _, ok := p.globalNodes[n.Name]; !ok {
				p.globalNames = append(p.globalNames, n.Name)
				p.globalNodes[n.Name] = n
			}
		}
		return true
	})

	for _, name := range p.globalNames {
		p.defGlobal(name, constant.NewInt(types.I32, 0))
	}

	return nil
}

func (p *Compiler) emitMainBody() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	p.curBlock = p.tinyMain.NewBlock(p.genLocalName("entry"))
	p.compileNone(p.tinyMain, p.file)

	p.curBlock.NewRet(constant.NewInt(i32, 0))
	return nil
}

// 递归编译节点
func (p *Compiler) compileNone(fn *ir.Func, node ast.Node) value.Value {
	if node == nil {
		return nil
	}

	switch node := node.(type) {
	case *ast.File:
		for _, n := range node.List {
			p.compileNone(fn, n)
		}
		return nil
	case *ast.BlockStmt:
		for _, n := range node.List {
			p.compileNone(fn, n)
		}
		return nil
	case *ast.IfStmt:
		ifTrue := fn.NewBlock(p.genLocalName("if.true"))
		ifEnd := fn.NewBlock(p.genLocalName("if.end"))
		ifFalse := ifEnd

		if node.Else != nil {
			ifFalse = fn.NewBlock(p.genLocalName("if.false"))
		}

		cond := p.compileNone(fn, node.Cond)
		p.curBlock.Term = ir.NewCondBr(cond, ifTrue, ifFalse)
		p.curBlock = ifTrue

		for _, n := range node.Body.List {
			p.compileNone(fn, n)
		}

		if p.curBlock != nil {
			p.curBlock.Term = ir.NewBr(ifEnd)
		}
		if node.Else != nil {
			p.curBlock = ifFalse
			p.compileNone(fn, node.Else)
			if p.curBlock != nil {
				p.curBlock.Term = ir.NewBr(ifEnd)
			}
		}

		p.curBlock = ifEnd
		return nil

	case *ast.RepeatStmt:
		repeatBody := fn.NewBlock(p.genLocalName("repeat.body"))
		repeatEnd := fn.NewBlock(p.genLocalName("repeat.end"))

		p.curBlock.Term = ir.NewBr(repeatBody)
		p.curBlock = repeatBody

		for _, n := range node.Body.List {
			p.compileNone(fn, n)
		}

		cond := p.compileNone(fn, node.Until)
		p.curBlock.Term = ir.NewCondBr(cond, repeatEnd, repeatBody)

		p.curBlock = repeatEnd
		return nil
	case *ast.AssignStmt:
		v := p.compileNone(fn, node.Value)
		p.curBlock.NewStore(v, p.getGlobal(node.Target.Name))
		return nil
	case *ast.ReadStmt:
		v := p.curBlock.NewCall(p.tinyRead)
		p.curBlock.NewStore(v, p.getGlobal(node.Target.Name))
		return nil
	case *ast.WriteStmt:
		x := p.compileNone(fn, node.Value)
		p.curBlock.NewCall(p.tinyWrite, x)
		return nil
	case *ast.Ident:
		return p.curBlock.NewLoad(i32, p.getGlobal(node.Name))
	case *ast.Number:
		return constant.NewInt(i32, int64(node.Value))
	case *ast.ParenExpr:
		return p.compileNone(fn, node.X)
	case *ast.BinaryExpr:
		left := p.compileNone(fn, node.X)
		right := p.compileNone(fn, node.Y)
		switch node.Op {
		case token.LT:
			cond := p.curBlock.NewICmp(enum.IPredSLE, left, right)
			cond = p.curBlock.NewICmp(enum.IPredNE, cond, constant.NewInt(i32, 0))
			return cond
		case token.EQ:
			cond := p.curBlock.NewICmp(enum.IPredEQ, left, right)
			cond = p.curBlock.NewICmp(enum.IPredNE, cond, constant.NewInt(i32, 0))
			return cond
		case token.PLUS:
			return p.curBlock.NewAdd(left, right)
		case token.MINUS:
			return p.curBlock.NewSub(left, right)
		case token.TIMES:
			return p.curBlock.NewMul(left, right)
		case token.OVER:
			return p.curBlock.NewSDiv(left, right)
		default:
			panic(fmt.Sprintf("unreachable: op = %v", node.Op))
		}
	default:
		panic("unreachable")
	}
}

// 生成局部变量名
func (p *Compiler) genLocalName(prefix string) string {
	if prefix == "" {
		prefix = "local"
	}
	id := p.nextId
	p.nextId++
	return fmt.Sprintf("%s.%04d", prefix, id)
}
