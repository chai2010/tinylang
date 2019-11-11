// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strconv"

// 记号类型
type Token int

const (
	// 特殊记号
	ILLEGAL Token = iota // 无效记号
	EOL                  // 行结尾
	EOF                  // 文件结尾
	COMMENT              // 注释, ; 行注释

	// 普通记号
	ID     // 标识符, main
	NUM    // 数字, 12345
	STRING // 字符串, "abc", 不支持原始的 'abc' 格式字符串
	COMMA  // 逗号

	keyword_beg
	// {{ 关键字开始

	// 伪指令
	START // 程序开始
	END   // 程序结束
	DC    // 定义常量
	DS    // 定义字符串

	// 内置的系统调用指令
	IN    // 输入
	OUT   // 输出
	EXIT  // 退出
	READ  // 新增, 读
	WRITE // 新增, 写

	// 寄存器
	GR0
	GR1
	GR2
	GR3
	GR4

	// 机器指令
	// 具体函数参考COMET文档
	HALT
	LD
	ST
	LEA
	ADD
	SUB
	MUL
	DIV
	MOD
	AND
	OR
	EOR
	SLA
	SRA
	SLL
	SRL
	CPA
	CPL
	JMP
	JPZ
	JMI
	JNZ
	JZE
	PUSH
	POP
	CALL
	RET

	// 系统调用
	SYSCALL

	// }} 关键字结束
	keyword_end
)

// 记号对应的字符串
var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOL:     "EOL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	ID:     "ID",
	NUM:    "NUM",
	STRING: "STRING",
	COMMA:  "COMMA",

	START: "START",
	END:   "END",
	DC:    "DC",
	DS:    "DS",

	IN:    "IN",
	OUT:   "OUT",
	EXIT:  "EXIT",
	READ:  "READ",
	WRITE: "WRITE",

	HALT: "ALT",
	LD:   "LD",
	ST:   "ST",
	LEA:  "LEA",
	ADD:  "ADD",
	SUB:  "SUB",
	MUL:  "MUL",
	DIV:  "DIV",
	MOD:  "MOD",
	AND:  "AND",
	OR:   "OR",
	EOR:  "EOR",
	SLA:  "SLA",
	SRA:  "SRA",
	SLL:  "SLL",
	SRL:  "SRL",
	CPA:  "CPA",
	CPL:  "CPL",
	JMP:  "JMP",
	JPZ:  "JPZ",
	JMI:  "JMI",
	JNZ:  "JNZ",
	JZE:  "JZE",
	PUSH: "USH",
	POP:  "POP",
	CALL: "ALL",
	RET:  "RET",

	SYSCALL: "SYSCALL",
}

// 是否为伪指令
func (tok Token) IsMACRO() bool {
	return tok == START || tok == END || tok == DC || tok == DS
}

// 是否为系统调用宏
func (tok Token) IsMACRO_SYSCALL() bool {
	return tok == IN || tok == OUT || tok == EXIT || tok == READ || tok == WRITE
}

// 是否为寄存器
func (tok Token) IsGR() bool {
	return tok == GR0 || tok == GR1 || tok == GR2 || tok == GR3 || tok == GR4
}

// 是否为关键字
func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}

// 字符串形式
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// 查找ID或对应的关键字
func Lookup(id string) Token {
	if tok, is_keyword := keywords[id]; is_keyword {
		return tok
	}
	return ID
}

// 判断名字是否为关键字
func IsKeyword(name string) bool {
	return Lookup(name).IsKeyword()
}
