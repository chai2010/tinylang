// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import "fmt"

// COMET指令类型
type OpType byte

// COMET机器指令
//
// 新增的指令: MUL, DIV, MOD, HALT, SYSCALL
const (
	HALT OpType = 0x00 // 停机
	LD   OpType = 0x01 // 取数, GR = (E)
	ST   OpType = 0x02 // 存数, E = (GR)
	LEA  OpType = 0x03 // 取地址, GR = E

	ADD OpType = 0x04 // 相加, GR = (GR)+(E)
	SUB OpType = 0x05 // 相减, GR = (GR)-(E)
	MUL OpType = 0x06 // 相乘, GR = (GR)*(E)
	DIV OpType = 0x07 // 相除, GR = (GR)/(E)
	MOD OpType = 0x08 // 取模, GR = (GR)%(E)

	AND OpType = 0x09 // 与, GR = (GR)&(E)
	OR  OpType = 0x0A // 或, GR = (GR)|(E)
	EOR OpType = 0x0B // 异或, GR = (GR)^(E)

	SLA OpType = 0x0C // 算术左移, GR = GR<<(E), 空出的的位置补0
	SRA OpType = 0x0D // 算术右移, GR = GR>>(E), 空出的的位置被置成第0位的值
	SLL OpType = 0x0E // 逻辑左移, GR = GR<<(E), 空出的的位置补0
	SRL OpType = 0x0F // 逻辑右移, GR = GR>>(E), 空出的的位置被置0

	CPA OpType = 0x10 // 算术比较, (GR)-(E), 有符号数, 设置FR
	CPL OpType = 0x11 // 逻辑比较, (GR)-(E), 无符号数, 设置FR

	JMP OpType = 0x12 // 无条件跳转, PC = E
	JPZ OpType = 0x13 // FR不小于跳转, PC = E
	JMI OpType = 0x14 // FR小于跳转, PC = E
	JNZ OpType = 0x15 // FR不等于0, PC = E
	JZE OpType = 0x16 // FR等于0跳转, PC = E

	PUSH OpType = 0x17 // 进栈, SP = (SP)-1, (SP) = E
	POP  OpType = 0x18 // 出栈, GR = ((SP)), SP = (SP)+1

	CALL OpType = 0x19 // 调用, SP = (SP)-1，(SP) = (PC)+2，PC = E
	RET  OpType = 0x1A // 返回, SP = (SP)+1

	SYSCALL OpType = 0xFF // 系统调用, GR0~GR3可用于交换数据
)

func (op OpType) Valid() bool {
	return int(op) < len(OpTab) && OpTab[op].Name != ""
}

func (op OpType) Size() int {
	if int(op) > len(OpTab) {
		return 0
	}
	return OpTab[op].Len
}

func (op OpType) String() string {
	if int(op) > len(OpTab) {
		return fmt.Sprintf("OpType(%d)", int(op))
	}
	if OpTab[op].Name == "" {
		return fmt.Sprintf("OpType(%d)", int(op))
	}

	return OpTab[op].Name
}

// COMET机器指令长度和名字
var OpTab = [...]struct {
	Op   OpType
	Name string
	Len  int
}{
	HALT: {HALT, "HALT", 1},

	LD:  {LD, "LD", 2},
	ST:  {ST, "ST", 2},
	LEA: {LEA, "LEA", 2},

	ADD: {ADD, "ADD", 2},
	SUB: {SUB, "SUB", 2},
	MUL: {MUL, "MUL", 2},
	DIV: {DIV, "DIV", 2},
	MOD: {MOD, "MOD", 2},

	AND: {AND, "AND", 2},
	OR:  {OR, "OR", 2},
	EOR: {EOR, "EOR", 2},

	SLA: {SLA, "SLA", 2},
	SRA: {SRA, "SRA", 2},
	SLL: {SLL, "SLL", 2},
	SRL: {SRL, "SRL", 2},

	CPA: {CPA, "CPA", 2},
	CPL: {CPL, "CPL", 2},

	JMP: {JMP, "JMP", 2},
	JPZ: {JPZ, "JPZ", 2},
	JMI: {JMI, "JMI", 2},
	JNZ: {JNZ, "JNZ", 2},
	JZE: {JZE, "JZE", 2},

	PUSH: {PUSH, "PUSH", 2},
	POP:  {POP, "POP", 1},
	CALL: {CALL, "CALL", 2},
	RET:  {RET, "RET", 1},

	SYSCALL: {SYSCALL, "SYSCALL", 2},
}
