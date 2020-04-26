// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"
)

// 完整指令
type Instruction struct {
	Label    string    // 标号, 独占一个指令, 对应 `Label DS 0`
	Op       OpType    // 指令码
	GR       byte      // 通用寄存器
	XR       byte      // 寻址寄存器
	ADR      uint16    // 跳转的目标地址
	ADRLabel string    // 跳转的目标标号
	Comment  string    // 注释
	Pos      token.Pos // 位置
}

// 解码指令
func (p *CPU) ParseInstruction(pc uint16) (ins *Instruction, ok bool) {
	ins = &Instruction{
		Op:  OpType(p.Mem[pc] / 0x100),
		GR:  byte(p.Mem[pc] % 0x100 / 0x10),
		XR:  byte(p.Mem[pc] % 0x10),
		ADR: uint16(p.Mem[pc+1]),
	}

	if !ins.Valid() {
		return nil, false
	}

	// OK
	return ins, true
}

// 有效的指令
func (p *Instruction) Valid() bool {
	if !p.Op.Valid() {
		return false
	}
	if p.GR > 4 || p.XR > 4 {
		return false
	}
	return true
}

// 字节码
func (p *Instruction) Bytes() []byte {
	if p.Label != "" {
		return []byte{} // Label DS 0
	}
	// 采用小端序
	switch p.Op.Size() {
	case 1:
		return []byte{
			byte(p.GR<<4) + byte(p.XR),
			byte(p.Op),
		}
	case 2:
		return []byte{
			byte(p.GR<<4) + byte(p.XR),
			byte(p.Op),
			byte(p.ADR),
			byte(p.ADR >> 8),
		}
	}
	return nil
}

// 字节码大小
func (p *Instruction) Size() uint16 {
	if p.Label != "" {
		return 0
	}
	return p.Op.Size()
}

// 格式化指令
func (p *Instruction) String() string {
	var buf bytes.Buffer

	// 无效指令
	if !p.Valid() {
		if p.Label == "" {
			fmt.Fprint(&buf, "invalid")
		} else {
			fmt.Fprint(&buf, p.Label)
		}
		return buf.String()
	}

	// 有标号(标号独占一个指令)
	if p.Label != "" {
		ds := &DataStore{
			Name:    p.Label,
			Comment: p.Comment,
		}
		return ds.String()
	}

	// op name
	opName := p.Op.String()
	if len(opName) < 8 {
		opName = (opName + "        ")[:8]
	}
	opName = "       " + opName
	adrName := fmt.Sprintf("%d", p.ADR)
	if p.ADRLabel != "" {
		adrName = p.ADRLabel
	}

	// syscall
	if p.Op == SYSCALL {
		fmt.Fprintf(&buf, "%s %v", opName, SyscallId(p.ADR))
	} else if p.Op.UseGR() {
		if p.Op.Size() == 2 {
			if p.XR != 0 {
				// OpName GR0, ADR, GR1
				fmt.Fprintf(&buf, "%s GR%d, %s+GR%d", opName, p.GR, adrName, p.XR)
			} else {
				// OpName GR0, ADR
				fmt.Fprintf(&buf, "%s GR%d, %s", opName, p.GR, adrName)
			}
		} else {
			// OpName GR0
			fmt.Fprintf(&buf, "%v GR%d", opName, p.GR)
		}
	} else {
		if p.Op.Size() == 2 {
			if p.XR != 0 {
				fmt.Fprintf(&buf, "%v %s+GR%d", opName, adrName, p.XR)
			} else {
				fmt.Fprintf(&buf, "%v %s", opName, adrName)
			}
			// OpName ADR
		} else {
			// OpName
			fmt.Fprintf(&buf, "%v", opName)
		}
	}

	var s = buf.String()
	if p.Comment != "" {
		s = (s + strings.Repeat(" ", 24+7))[:24+7] + p.Comment
	}
	return s
}
