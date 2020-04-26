// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TINY 语言编译器, 生成 COMET 字节码.
//
// 每次运算的结果和返回值均保存到 GR0 寄存器.
//
// 字节码布局:
//  +---------------+ 0
//  | JMP start     | 跳转到开始地址, 该指令2个字节
//  +---------------+ 2
//  | __AC__        | 临时变量, 内部实用
//  | global_var_0  | 全局变量, 地址从 2 开始计算
//  | global_var_1  |
//  | ...           |
//  | global_var_n  |
//  +---------------+ 2*X
//  | ...           |
//  | instruction_1 | 指令序列, 地址需要对齐到2的倍数
//  | instruction_2 |
//  | ...           |
//  | instruction_n |
//  +---------------+
//  | ...           | 堆空间
//  | HEAP          |
//  | ...           |
//  +---------------+
//  | ...           |
//  | ...           |
//  | STACK         | 栈空间
//  +---------------+ 0xFC00
//  | ...           | 外设地址
//  +---------------+ 0xFFFF
//
package compiler

import (
	"fmt"
	gotoken "go/token"

	"github.com/chai2010/tinylang/comet"
	"github.com/chai2010/tinylang/tiny/ast"
	"github.com/chai2010/tinylang/tiny/token"
)

var _ = fmt.Sprint

// 内部临时变量
const __AC__ = "__AC__"

// TINY 编译器
type Compiler struct {
	fset   *gotoken.FileSet
	prog   *Program
	nextId int
}

// NewCompiler 构造编译器对象
func NewCompiler() *Compiler {
	p := &Compiler{
		prog: NewProgram(),
	}
	return p.reset()
}

// Compile 编译 TINY 程序对应的语法树
func (p *Compiler) Compile(f *ast.File) ([]byte, error) {
	fset := gotoken.NewFileSet()
	if len(f.Data) > 0 {
		fset.AddFile(f.Name, 1, len(f.Data)).SetLinesForContent(f.Data)
	}

	p.fset = fset
	if err := p.reset().compile(f); err != nil {
		return nil, err
	}

	p.emit(&comet.Instruction{Op: comet.HALT})
	return p.prog.Bytes(), nil
}

func (p *Compiler) CASLString() string {
	return p.prog.String()
}

// 清空内部状态
func (p *Compiler) reset() *Compiler {
	p.nextId = 0
	p.prog = NewProgram()
	p.prog.DefineName(__AC__, 1, "; temp", gotoken.NoPos)
	return p
}

// 递归编译节点
func (p *Compiler) compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.File:
		for _, s := range node.List {
			if err := p.compile(s); err != nil {
				return err
			}
		}

	case *ast.BlockStmt:
		for _, s := range node.List {
			if err := p.compile(s); err != nil {
				return err
			}
		}

	case *ast.IfStmt:
		bin, ok := node.Cond.(*ast.BinaryExpr)
		if !ok {
			return fmt.Errorf("%v: if.cond must be compare expr ", p.fset.Position(node.Pos()))
		}
		if bin.Op != token.LT && bin.Op != token.EQ {
			return fmt.Errorf("%v: if.cond must be compare expr ", p.fset.Position(node.Pos()))
		}

		// __AC__ 是临时变量的地址
		__AC__, ok := p.prog.LookupName(__AC__)
		if !ok {
			panic("__AC__ not found")
		}

		// 右边的表达式结果存在 GR0
		if err := p.compile(bin.Y); err != nil {
			return err
		}

		// PUSH GRO
		p.emit(&comet.Instruction{Op: comet.ST, GR: 0, ADR: __AC__, ADRLabel: "__AC__"})
		p.emit(&comet.Instruction{Op: comet.PUSH, ADR: __AC__, ADRLabel: "__AC__"})

		// 左边的表达式结果存在 GR0
		if err := p.compile(bin.X); err != nil {
			return err
		}

		// POP AC
		p.emit(&comet.Instruction{Op: comet.POP, GR: 1})
		p.emit(&comet.Instruction{Op: comet.ST, GR: 1, ADR: __AC__, ADRLabel: "__AC__"})

		var (
			label_if_else = p.genLabel()
			label_if_end  = p.genLabel()
		)

		// 比较
		p.emit(&comet.Instruction{
			Op: comet.CPA, GR: 0, ADR: __AC__, ADRLabel: "__AC__",
			Comment: fmt.Sprintf("; %v", p.fset.Position(node.Pos())),
		})
		switch bin.Op {
		case token.LT:
			if node.Else != nil {
				p.emit(&comet.Instruction{
					Op: comet.JPZ, GR: 0, ADR: 0, ADRLabel: label_if_else,
					Comment: fmt.Sprintf("; %v goto if_else", p.fset.Position(node.Pos())),
				})
			} else {
				p.emit(&comet.Instruction{
					Op: comet.JPZ, GR: 0, ADR: 0, ADRLabel: label_if_end,
					Comment: fmt.Sprintf("; %v goto if_end", p.fset.Position(node.Pos())),
				})
			}
		case token.EQ:
			if node.Else != nil {
				p.emit(&comet.Instruction{
					Op: comet.JNZ, GR: 0, ADR: 0, ADRLabel: label_if_else,
					Comment: fmt.Sprintf("; %v goto if_else", p.fset.Position(node.Pos())),
				})
			} else {
				p.emit(&comet.Instruction{
					Op: comet.JNZ, GR: 0, ADR: 0, ADRLabel: label_if_end,
					Comment: fmt.Sprintf("; %v goto if_else", p.fset.Position(node.Pos())),
				})
			}
		}

		if err := p.compile(node.Body); err != nil {
			return err
		}

		if node.Else != nil {
			p.emit(&comet.Instruction{
				Op: comet.JMP, GR: 0, ADR: 0, ADRLabel: label_if_end,
				Comment: fmt.Sprintf("; %v goto if_else", p.fset.Position(node.Pos())),
			})
			p.emit(&comet.Instruction{
				Label:   label_if_else,
				Comment: fmt.Sprintf("; %v if_else", p.fset.Position(node.Pos())),
			})
			if err := p.compile(node.Else); err != nil {
				return err
			}
		}
		p.emit(&comet.Instruction{
			Label:   label_if_end,
			Comment: fmt.Sprintf("; %v if_end", p.fset.Position(node.Pos())),
		})

	case *ast.RepeatStmt:
		bin, ok := node.Until.(*ast.BinaryExpr)
		if !ok {
			return fmt.Errorf("%v: repeat.until must be compare expr ", p.fset.Position(node.Pos()))
		}
		if bin.Op != token.LT && bin.Op != token.EQ {
			return fmt.Errorf("%v: repeat.until must be compare expr ", p.fset.Position(node.Pos()))
		}

		label_repeat_begin := p.genLabel()
		p.emit(&comet.Instruction{
			Label:   label_repeat_begin,
			Comment: fmt.Sprintf("; %v repeat_begin", p.fset.Position(node.Pos())),
		})

		if err := p.compile(node.Body); err != nil {
			return err
		}

		// __AC__ 是临时变量的地址
		__AC__, ok := p.prog.LookupName(__AC__)
		if !ok {
			panic("__AC__ not found")
		}

		// 右边的表达式结果存在 GR0
		if err := p.compile(bin.Y); err != nil {
			return err
		}

		// PUSH GRO
		p.emit(&comet.Instruction{Op: comet.ST, GR: 0, ADR: __AC__, ADRLabel: "__AC__"})
		p.emit(&comet.Instruction{Op: comet.PUSH, ADR: __AC__, ADRLabel: "__AC__"})

		// 左边的表达式结果存在 GR0
		if err := p.compile(bin.X); err != nil {
			return err
		}

		// POP AC
		p.emit(&comet.Instruction{Op: comet.POP, GR: 1})
		p.emit(&comet.Instruction{Op: comet.ST, GR: 1, ADR: __AC__, ADRLabel: "__AC__"})

		// 比较
		p.emit(&comet.Instruction{
			Op: comet.CPA, GR: 0, ADR: __AC__, ADRLabel: "__AC__",
			Comment: fmt.Sprintf("; %v", p.fset.Position(node.Pos())),
		})
		switch bin.Op {
		case token.LT:
			p.emit(&comet.Instruction{
				Op: comet.JPZ, GR: 0, ADR: 0, ADRLabel: label_repeat_begin,
				Comment: fmt.Sprintf("; %v goto repeat_begin", p.fset.Position(node.Pos())),
			})
		case token.EQ:
			p.emit(&comet.Instruction{
				Op: comet.JNZ, GR: 0, ADR: 0, ADRLabel: label_repeat_begin,
				Comment: fmt.Sprintf("; %v goto if_else", p.fset.Position(node.Pos())),
			})
		}

	case *ast.AssignStmt:
		// 定义变量
		adr, ok := p.prog.LookupName(node.Target.Name)
		if !ok {
			comment := fmt.Sprintf("; %s", p.fset.Position(node.Target.Pos()))
			adr, _ = p.prog.DefineName(node.Target.Name, 1, comment, gotoken.NoPos)
		}

		// 结果保存到 GR0
		if err := p.compile(node.Value); err != nil {
			return err
		}

		// 结果保存到变量
		p.emit(&comet.Instruction{
			Op:       comet.ST,
			GR:       0,
			ADR:      adr,
			ADRLabel: node.Target.Name,
			Comment:  fmt.Sprintf("; %v", p.fset.Position(node.Pos())),
		})

	case *ast.ReadStmt:
		// 获取变量的地址
		adr, ok := p.prog.LookupName(node.Target.Name)
		if !ok {
			comment := fmt.Sprintf("; %s", p.fset.Position(node.Target.Pos()))
			adr, _ = p.prog.DefineName(node.Target.Name, 1, comment, gotoken.NoPos)
		}

		// 系统调用 SYSCALL_READ, 读取到 GR0
		p.emit(&comet.Instruction{
			Op:      comet.SYSCALL,
			ADR:     uint16(comet.SYSCALL_READ),
			Comment: fmt.Sprintf("; %v", p.fset.Position(node.Pos())),
		})

		// 保存到变量
		p.emit(&comet.Instruction{
			Op:       comet.ST,
			GR:       0,
			ADR:      adr,
			ADRLabel: node.Target.Name,
			Comment:  fmt.Sprintf("; %v", p.fset.Position(node.Pos())),
		})

	case *ast.WriteStmt:
		// 结果保存到 GR0
		if err := p.compile(node.Value); err != nil {
			return err
		}

		// 系统调用 SYSCALL_WRITE, 输出 GR0
		p.emit(&comet.Instruction{
			Op:      comet.SYSCALL,
			ADR:     uint16(comet.SYSCALL_WRITE),
			Comment: fmt.Sprintf("; %v", p.fset.Position(node.Pos())),
		})

	case *ast.Ident:
		// 获取变量的地址
		adr, ok := p.prog.LookupName(node.Name)
		if !ok {
			return fmt.Errorf("%v not found", node)
		}

		// LD GR0, ADR
		p.emit(&comet.Instruction{
			Op:       comet.LD,
			GR:       0,
			ADR:      adr,
			ADRLabel: node.Name,
			Comment:  fmt.Sprintf("; %v ident=%s", p.fset.Position(node.Pos()), node.Name),
		})

	case *ast.Number:
		p.emit(&comet.Instruction{
			Op:      comet.LEA,
			GR:      0,
			ADR:     uint16(node.Value),
			Comment: fmt.Sprintf("; %v num=%d", p.fset.Position(node.Pos()), node.Value),
		})

	case *ast.ParenExpr:
		if err := p.compile(node.X); err != nil {
			return err
		}

	case *ast.BinaryExpr:
		switch node.Op {
		case token.LT, token.EQ:
			return fmt.Errorf("%v: compare must be if.cond or repeat.until", p.fset.Position(node.Pos()))

		case token.PLUS, token.MINUS, token.TIMES, token.OVER:
			// __AC__ 是临时变量的地址
			__AC__, ok := p.prog.LookupName(__AC__)
			if !ok {
				panic("__AC__ not found")
			}

			// 右边的表达式结果存在 GR0
			if err := p.compile(node.Y); err != nil {
				return err
			}

			// PUSH GRO
			p.emit(&comet.Instruction{Op: comet.ST, GR: 0, ADR: __AC__, ADRLabel: "__AC__"})
			p.emit(&comet.Instruction{Op: comet.PUSH, ADR: __AC__, ADRLabel: "__AC__"})

			// 左边的表达式结果存在 GR0
			if err := p.compile(node.X); err != nil {
				return err
			}

			// POP AC
			p.emit(&comet.Instruction{Op: comet.POP, GR: 1})
			p.emit(&comet.Instruction{Op: comet.ST, GR: 1, ADR: __AC__, ADRLabel: "__AC__"})

			// 进行运算
			switch node.Op {
			case token.PLUS:
				p.emit(&comet.Instruction{
					Op: comet.ADD, GR: 0, ADR: __AC__, ADRLabel: "__AC__",
					Comment: fmt.Sprintf("; %v binop:'+'", p.fset.Position(node.Pos())),
				})
			case token.MINUS:
				p.emit(&comet.Instruction{
					Op: comet.SUB, GR: 0, ADR: __AC__, ADRLabel: "__AC__",
					Comment: fmt.Sprintf("; %v binop:'-'", p.fset.Position(node.Pos())),
				})
			case token.TIMES:
				p.emit(&comet.Instruction{
					Op: comet.MUL, GR: 0, ADR: __AC__, ADRLabel: "__AC__",
					Comment: fmt.Sprintf("; %v binop:'*'", p.fset.Position(node.Pos())),
				})
			case token.OVER:
				p.emit(&comet.Instruction{
					Op: comet.DIV, GR: 0, ADR: __AC__, ADRLabel: "__AC__",
					Comment: fmt.Sprintf("; %v binop:'/'", p.fset.Position(node.Pos())),
				})
			}
		}
	}
	return nil
}

func (p *Compiler) emit(ins *comet.Instruction) {
	p.prog.AppendInstruction(ins)
}

// 生成新的CASL标号
func (p *Compiler) genLabel() string {
	id := p.nextId
	p.nextId++
	return fmt.Sprintf("_L%04d", id)
}
