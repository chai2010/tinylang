// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

// CASL机器指令
const (
	HALT = 0x00 // 停机
	LD   = 0x01 // 取数, GR = (E)
	ST   = 0x02 // 存数, E = (GR)
	LEA  = 0x03 // 取地址, GR = E

	ADD = 0x04 // 相加, GR = (GR)+(E)
	SUB = 0x05 // 相减, GR = (GR)-(E)
	MUL = 0x06 // 相乘, GR = (GR)*(E)
	DIV = 0x07 // 相除, GR = (GR)/(E)
	MOD = 0x08 // 取模, GR = (GR)%(E)

	AND = 0x09 // 与, GR = (GR)&(E)
	OR  = 0x0A // 或, GR = (GR)|(E)
	EOR = 0x0B // 异或, GR = (GR)^(E)

	CPA = 0x0C // 算术比较, (GR)-(E), 有符号数, 设置FR
	CPL = 0x0D // 逻辑比较, (GR)-(E), 无符号数, 设置FR

	SLA = 0x0E // 算术左移, 空出的的位置补0
	SRA = 0x0F // 算术右移, 空出的的位置被置成第0位的值
	SLL = 0x10 // 逻辑左移, 空出的的位置补0
	SRL = 0x11 // 逻辑右移, 空出的的位置被置0

	JMP = 0x12 // 无条件跳转, PC = E
	JPZ = 0x13 // 不小于跳转, PC = E
	JMI = 0x14 // 小于跳转, PC = E
	JNZ = 0x15 // 不等于0, PC = E
	JZE = 0x16 // 等于0跳转, PC = E

	PUSH = 0x17 // 进栈, SP = (SP)-1, (SP) = E
	POP  = 0x18 // 出栈, GR = ((SP)), SP = (SP)+1

	CALL = 0x19 // 调用, SP = (SP)-1，(SP) = (PC)+2，PC = E
	RET  = 0x1A // 返回, SP = (SP)+1

	SYSCALL = 0xFF // 系统调用, GR0~GR3可用于交换数据
)

// COMET机器指令长度和名字
var OpTab = []struct {
	Op   int
	Name string
	Len  int
}{
	{HALT, "HALT", 1},

	{LD, "LD", 2}, {ST, "ST", 2}, {LEA, "LEA", 2},

	{ADD, "ADD", 2}, {SUB, "SUB", 2},
	{MUL, "MUL", 2}, {DIV, "DIV", 2}, {MOD, "MOD", 2},
	{AND, "AND", 2}, {OR, "OR", 2}, {EOR, "EOR", 2},

	{SLA, "SLA", 2}, {SRA, "SRA", 2}, {SLL, "SLL", 2}, {SRL, "SRL", 2},

	{CPA, "CPA", 2}, {CPL, "CPL", 2},

	{JMP, "JMP", 2},
	{JPZ, "JPZ", 2}, {JMI, "JMI", 2}, {JNZ, "JNZ", 2}, {JZE, "JZE", 2},

	{PUSH, "PUSH", 2}, {POP, "POP", 1},
	{CALL, "CALL", 2}, {RET, "RET", 1},

	{SYSCALL, "SYSCALL", 1},
}
