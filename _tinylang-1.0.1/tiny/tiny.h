// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef TINY_H
#define TINY_H

#include <stdio.h>
#include <ctype.h>
#include <string.h>
#include <stdlib.h>

#define		TRUE	1
#define		FALSE	0

/* 计算数组的元素个数 */

#define NELEMS(x) (sizeof(x) / sizeof((x)[0]))

/* 抛出异常 */

#define RAISE(e) do {	\
	fprintf(stderr, "未被捕获的异常 %s ", (e));	\
	fprintf(stderr, "引发位置 %s:%d\n", __FILE__, __LINE__);	\
	fprintf(stderr, "终止...\n");		\
	fflush(stderr);		\
	abort();	\
}while(0)

/* TINY语言的各种记号 */

typedef enum { ENDFILE, ERROR, ID, NUM,
	IF, THEN, ELSE, END, REPEAT, UNTIL, READ, WRITE,
	ASSIGN, EQ, LT, PLUS, MINUS, TIMES, OVER, LPAREN, RPAREN, SEMI
} TokenType;

/* 各种类型 */

typedef enum { StmtK, ExpK } NodeKind;
typedef enum { IfK, RepeatK, AssignK, ReadK, WriteK } StmtKind;
typedef enum { OpK, ConstK, IdK } ExpKind;
typedef enum { Void, Integer, Boolean } ExpType;

/* 语法树结构 */

typedef struct treeNode {
	struct treeNode * child[3];
	struct treeNode * sibling;
	int line;
	NodeKind nodekind;
	union { StmtKind stmt; ExpKind exp;} kind;
	union { TokenType op; char *val, *name; } attr;
	ExpType type;
} TreeNode;

/* 全局参数 */

extern char pgmName[30];
extern char codName[30];
extern char lstName[30];

extern FILE* source;
extern FILE* listing;
extern FILE* code;

/* 标志信息 */

extern int line;
extern int Error;

// 入口函数
extern int tinyMain(int argc, char *argv[]);

#endif
