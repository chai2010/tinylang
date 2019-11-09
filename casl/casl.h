// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CASL_H
#define CASL_H

#include <stdio.h>
#include <ctype.h>
#include <string.h>
#include <stdlib.h>

/* 计算数组元素个数 */

#define NELEMS(x) (sizeof(x) / sizeof((x)[0]))

/* 抛出异常，e 异常名称为字符串*/

#define RAISE(e) do {	\
	fprintf(stderr, "未被捕获的异常 %s ", (e));	\
	fprintf(stderr, "引发位置 %s:%d\n", __FILE__, __LINE__);	\
	fprintf(stderr, "终止...\n\n");		\
	fflush(stderr);		\
	abort();	\
}while(0)

/* casl汇编程序的错误处理宏， msg为消息字符串 */

#define QUIT(msg) do {	\
	printf("系统参数: %s 文件 %d 行\n", __FILE__, __LINE__);	\
	printf("错误信息：%s 文件 %d 行 %s\n", pgmName, line, (msg));	\
	printf("退出...\n\n");	\
	exit(1);	\
}while(0)

#define		MEMSIZE		0x10000		/* 内存大小	*/
#define		pc_max		0xFB00		/* 最大地址	*/
#define		sp_start	0xFB00		/* 栈地址	*/

#define		ac_comet	0xFE00		/* 临时变量	*/

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

#define		MEMSIZE		0x10000		/* 内存大小	*/
#define		WORDSIZE	6		/* 标号长度	*/
#define		LINESIZE	72		/* 指令长度	*/
#define		pc_start	0x0000		/* 起始地址	*/
#define		pc_max		0xFB00		/* 最大地址	*/

/* 无符号 short 型宏 */

#define		off_t		unsigned short

/* 记号类型，包含全部的指令 */

typedef enum {
	
	/* COMET机器指令
	   新增加4个指令:
	   MUL, DIV, MOD, HALT */
	HALT, LD, ST, LEA,
	ADD, SUB, MUL, DIV, MOD,
	AND, OR, EOR,
	SLA, SRA, SLL, SRL,
	CPA, CPL,
	JMP, JPZ, JMI, JNZ, JZE,
	PUSH, POP, CALL, RET,
	
	/* 宏指令（共5个）
	   新增READ, WRITE */
	READ, WRITE, IN, OUT, EXIT,
	
	/* 汇编伪命令 */
	START, END, DC, DS,
	
	/* 其他的标号 */
	ID, NUM, STRING, COMMA, ENDLINE
} TokenType;

extern off_t mem[MEMSIZE];	/* 64k 内存 */

extern char pgmName[32];	/* 汇编程序 */
extern char codName[32];	/* 机器代码 */

extern FILE *source, *code;	/* 文件指针 */

extern int pc;	/* 指令地址 */
extern int line;		/* 行计数器 */

extern int state;		/* 状态标志 */
extern int Error;		/* 错误标志 */


extern int caslMain(int n, char *v[]);

#endif
