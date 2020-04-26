// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"

	"github.com/chai2010/tinylang/comet"
)

// 字节码
type Program struct {
	Globals      []comet.DataStore    // 全局变量, 在开头, 从 0 开始
	Instructions []*comet.Instruction // 指令列表
}

// NewProgram 返回 程序对象(用于产生字节码)
func NewProgram() *Program {
	return &Program{}
}

// 获取启动地址
func (p *Program) GetStartPC() uint16 {
	var pc uint16
	for _, x := range p.Globals {
		pc += x.Size
	}
	return pc
}

// 定义一个变量, 不能重复定义
func (p *Program) DefineName(name string, size uint16, comment string, pos token.Pos) (uint16, error) {
	for _, x := range p.Globals {
		if x.Name == name {
			return 0, fmt.Errorf("%s exists", name)
		}
	}
	off := p.GetStartPC()
	p.Globals = append(p.Globals, comet.DataStore{
		Name:    name,
		Size:    size,
		Comment: comment,
		Pos:     pos,
	})
	return off, nil
}

// 查询标号的地址
func (p *Program) LookupName(name string) (addr uint16, ok bool) {
	for i, x := range p.Globals {
		if x.Name == name {
			return addr, true
		}
		addr += p.Globals[i].Size
	}
	for _, ins := range p.Instructions {
		if ins.Label == name {
			return addr, true
		}
		addr += ins.Op.Size()
	}
	return 0, false
}

// 添加指令
func (p *Program) AppendInstruction(ins *comet.Instruction) {
	p.Instructions = append(p.Instructions, ins)
}

// 填充标号地址
func (p *Program) FixAdrLabels() error {
	for i, ins := range p.Instructions {
		if ins.ADRLabel != "" {
			adr, ok := p.LookupName(ins.ADRLabel)
			if !ok {
				return fmt.Errorf("label '%s' not found", ins.ADRLabel)
			}
			p.Instructions[i].ADR = adr
		}
	}
	return nil
}

// 全部数据和指令的文本格式, 采用 CASL 汇编语言表示
func (p *Program) String() string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, ";      CASL Assmebly Language")
	fmt.Fprintln(&buf)

	s := ("       START CASL00" + strings.Repeat(" ", 24+7))[:24+7]
	fmt.Fprintf(&buf, "%s; start pc=0x%04x\n", s, p.GetStartPC())
	fmt.Fprintln(&buf)

	if len(p.Globals) > 0 {
		for _, x := range p.Globals {
			fmt.Fprintln(&buf, x.String())
		}
		fmt.Fprintln(&buf)
	}

	// CASL00 DS 0
	fmt.Fprintln(&buf, &comet.DataStore{Name: "CASL00"})
	fmt.Fprintln(&buf)

	for _, x := range p.Instructions {
		fmt.Fprintln(&buf, x.String())
	}
	fmt.Fprintln(&buf)

	fmt.Fprintln(&buf, "       END")
	return buf.String()
}

// 字节码
func (p *Program) Bytes() []byte {
	var pc = p.GetStartPC()

	// 开头4个字节预留给 comet.BytecodeHeader
	var data = make([]byte, 4)

	// 生成 DS 部分字节码
	data = append(data, make([]byte, pc*2)...)

	// 修复地址
	if err := p.FixAdrLabels(); err != nil {
		panic(err)
	}

	// 生成字节码
	for _, x := range p.Instructions {
		data = append(data, x.Bytes()...)
	}

	// 重写头部 comet.BytecodeHeader
	data[0] = byte(pc)
	data[1] = byte(pc >> 8)
	data[2] = byte((len(data) - 4) / 2)
	data[3] = byte(((len(data) - 4) / 2) >> 8)

	// OK
	return data
}
