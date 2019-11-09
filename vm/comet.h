// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef COMET_H
#define COMET_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* 计算数组元素个数 */

#define		NELEMS(x)	(sizeof(x) / sizeof((x)[0]))

#define		MEMSIZE		0x10000		/* 内存大小 */
#define		pc_max		0xFC00		/* 最大地址 */
#define		sp_start	0xFC00		/* 栈地址	*/

#define		IO_ADDR		0xFD10		/* IO地址	*/
#define		IO_FLAG		0xFD11		/* IO状态	*/

#define		IO_FIO		0x0100		/* 标志		*/
#define		IO_TYPE		0x1C00		/* 类型		*/
#define		IO_MAX		0x00FF		/* 传输个数	*/
#define		IO_ERROR	0x0200		/* 错误		*/

#define		IO_IN		0x0000		/* 输入		*/
#define		IO_OUT		0x0100		/* 输出		*/
#define		IO_CHR		0x0400		/* 字符		*/
#define		IO_OCT		0x0800		/* 八进制	*/
#define		IO_DEC		0x0C00		/* 十进制	*/
#define		IO_HEX		0x1000		/* 十六进制	*/

/* comet计算机地址类型 */

#define		off_t		unsigned short

/* comet机器指令 */

typedef enum {
	HALT, LD, ST, LEA,
	ADD, SUB, MUL, DIV, MOD,
	AND, OR, EOR,
	SLA, SRA, SLL, SRL,
	CPA, CPL,
	JMP, JPZ, JMI, JNZ, JZE,
	PUSH, POP, CALL, RET,
} OpType;

/* comet机器指令 */

struct { OpType op; char *str; int len; } opTab[] = {
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
	{CALL, "CALL", 2}, {RET, "RET", 1}
};

/* 调试指令 */

typedef enum {
	HELP, GO, STEP, JUMP, REGS, IMEM, DMEM,
	ALTER, TRACE, PRINT, CLEAR, QUIT } DbType;

/* 调试指令 */

struct { DbType db; char *s1, *s2; } dbTab[] = {
	{HELP , "help" , "h"},
	{GO   , "go"   , "g"},
	{STEP , "step" , "s"},
	{JUMP , "jump" , "j"},
	{REGS , "regs" , "r"},
	{IMEM , "iMem" , "i"},
	{DMEM , "dMem" , "d"},
	{ALTER, "alter", "a"},
	{TRACE, "trace", "t"},
	{PRINT, "print", "p"},
	{CLEAR, "clear", "c"},
	{QUIT , "quit" , "q"}
};

/* comet计算机结构 */

struct comet {
	off_t pc;
	short fr;
	short gr[5];
	short mem[MEMSIZE];
} cmt;

/* 程序名称 */

char pgmName[32];

/* 程序文件指针 */

FILE * source;

/* 调试开关 */

int debug = 0;

#endif
