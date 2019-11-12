// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bytes"
	"fmt"
)

// 完整指令
type Instruction struct {
	Label     string // 标号
	Op        OpType // 指令码
	GR        uint16 // 通用寄存器
	XR        uint16 // 寻址寄存器
	ADR       uint16 // 地址
	SyscallId uint8  // 系统调用号
	Comment   string // 注释
}

// 解码指令
func (p *CPU) ParseInstruction(pc uint16) (ins *Instruction, ok bool) {
	ins = &Instruction{
		Op:        OpType(p.Mem[pc] / 0x100),
		GR:        p.Mem[pc] % 0x100 / 0x10,
		XR:        p.Mem[pc] % 0x10,
		ADR:       uint16(p.Mem[pc+1]),
		SyscallId: uint8(p.Mem[pc] % 0x100),
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

	// 有标号
	if p.Label != "" {
		fmt.Fprint(&buf, p.Label+" ")
		return buf.String()
	}

	// 系统调用单独处理
	// 系统调用号为十六进制格式 [??]
	if p.Op == SYSCALL {
		return fmt.Sprintf("%v [%02x]", p.Op, p.SyscallId)
	}

	// 包含GR参数
	if p.Op.UseGR() {
		if p.Op.Size() == 2 {
			if p.XR != 0 {
				// OpName GR0, ADR, GR1
				fmt.Fprintf(&buf, ", %04x, GR%d", p.ADR, p.XR)
			} else {
				// OpName GR0, ADR
				fmt.Fprintf(&buf, "%v GR%d", p.Op, p.GR)
			}
		} else {
			// OpName GR0
			fmt.Fprintf(&buf, "%v GR%d", p.Op, p.GR)
		}
	} else {
		if p.Op.Size() == 2 {
			// OpName ADR
			fmt.Fprintf(&buf, "%v, %04x", p.Op, p.ADR)
		} else {
			// OpName
			fmt.Fprintf(&buf, "%v", p.Op)
		}
	}

	return buf.String()
}
