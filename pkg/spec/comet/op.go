// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

// CASL机器指令
const (
	HALT = 0X0 // 停机
	LD   = 0X1 // 取数, GR = (E)
	ST   = 0X2 // 存数, E = (GR)
	LEA  = 0X3 // 取地址, GR = E

	ADD = 0X4 // 相加, GR = (GR)+(E)
	SUB = 0X5 // 相减, GR = (GR)-(E)
	MUL = 0X6 // 相乘, GR = (GR)*(E)
	DIV = 0X7 // 相除, GR = (GR)/(E)
	MOD = 0X8 // 取模, GR = (GR)%(E)

	AND = 0X9 // 与, GR = (GR)&(E)
	OR  = 0XA // 或, GR = (GR)|(E)
	EOR = 0XB // 异或, GR = (GR)^(E)

	CPA = 0XC // 算术比较, (GR)-(E), 有符号数, 设置FR
	CPL = 0XD // 逻辑比较, (GR)-(E), 无符号数, 设置FR

	SLA = 0XE  // 算术左移, 空出的的位置补0
	SRA = 0XF  // 算术右移, 空出的的位置被置成第0位的值
	SLL = 0X10 // 逻辑左移, 空出的的位置补0
	SRL = 0X11 // 逻辑右移, 空出的的位置被置0

	JMP = 0X12 // 无条件跳转, PC = E
	JPZ = 0X13 // 不小于跳转, PC = E
	JMI = 0X14 // 小于跳转, PC = E
	JNE = 0X15 // 不等于跳转, PC = E
	JZE = 0X16 // 等于跳转, PC = E

	PUSH = 0X17 // 进栈, SP = (SP)-1, (SP) = E
	POP  = 0X18 // 出栈, GR = ((SP)), SP = (SP)+1

	CALL = 0X19 // 调用, SP = (SP)-1，(SP) = (PC)+2，PC = E
	RET  = 0X1A // 返回, SP = (SP)+1
)
