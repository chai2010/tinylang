// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef UTIL_H
#define UTIL_H

#include "tiny.h"

/* 原子化的字符串 */

static struct string {
	struct string *link;
	char *str;
} *str_buckets[211];

/* 字符串hash化 */

int
hash(char *key, int M)
{
	int i = 0, h = 0, a = 17;
	for(; key[i] != '\0'; i++) {
		h = (h * a + key[i]) % M;
	}
	return h;
}

/* 生成新的字符串原子 */

char *
str_new(char *s)
{
	struct string *p;
	int h;
	h = hash(s, NELEMS(str_buckets));
	for(p = str_buckets[h]; p; p = p->link)
		if(!strcmp(s, p->str)) return p->str;
	p = malloc(sizeof(*p) + strlen(s) + 1);
	if(p == NULL) RAISE("内存溢出");
	p->str = (char *)(p + 1);
	strcpy(p->str, s);
	p->link = str_buckets[h];
	str_buckets[h] = p;
	return p->str;
}

/* 释放所有原子 */

void
str_free(void)
{
	struct string *p, *q;
	int i;
	for (i = 0; i < NELEMS(str_buckets); i++) {
		for(p = str_buckets[i]; p; p = q) {
			q = p->link; free(p);
		}
		str_buckets[i] = NULL;
	}
}

/* 生成新的CASL标号 */

char *
new_label(void)
{
	static char label[6+1] = "AAAAAA";
	static int i = 0;
	if(label[0] > 'Z') RAISE("标号溢出");
	label[i = (i+1) % 6]++;
	return str_new(label);
}

/* 分配一个语句类型的结点 */

TreeNode *
newTreeNode(NodeKind nodeKind, int kind)
{
	int i;
	TreeNode *t = malloc(sizeof(*t));
	if(t == NULL) RAISE("内存溢出");
	for (i = 0; i < NELEMS(t->child); i++)
		t->child[i] = NULL;
	t->sibling = NULL;
	t->nodekind = nodeKind;
	t->kind.stmt = kind;
	t->line = line;
	t->type = Void;
	return t;
}

/* 打印记号 */

void
printToken(TokenType token, const char *str)
{ 
	switch(token) {
		case IF:
		case THEN:
		case ELSE:
		case END:
		case REPEAT:
		case UNTIL:
		case READ:
		case WRITE:
			fprintf(listing, "关键字: %s\n", str);
			break;
		case ASSIGN: fprintf(listing,":=\n"); break;
		case LT: fprintf(listing,"<\n"); break;
		case EQ: fprintf(listing,"=\n"); break;
		case LPAREN: fprintf(listing,"(\n"); break;
		case RPAREN: fprintf(listing,")\n"); break;
		case SEMI: fprintf(listing,";\n"); break;
		case PLUS: fprintf(listing,"+\n"); break;
		case MINUS: fprintf(listing,"-\n"); break;
		case TIMES: fprintf(listing,"*\n"); break;
		case OVER: fprintf(listing,"/\n"); break;
		case ENDFILE:
			fprintf(listing, "文件结束\n");
			break;
		case NUM:
			fprintf(listing, "数值, 值= %s\n", str);
			break;
		case ID:
			fprintf(listing, "变量, 名称= %s\n", str);
			break;
		case ERROR:
			fprintf(listing, "错误: %s\n", str);
			break;
		default:
			fprintf(listing, "未知符号: %d\n", token);
			break;
	}
}

/* 打印语法树 */

void
printTree(TreeNode *tree)
{
	static int indentno = 0;
	int i;
	
	indentno++;
	
	while (tree != NULL) {
		for(i = 1; i < indentno; i++)
			fprintf(listing,"\t");
		if(tree->nodekind == StmtK) {
			switch(tree->kind.stmt) {
				case IfK:
					fprintf(listing,"If判断\n");
					break;
				case RepeatK:
					fprintf(listing,"Repeat循环\n");
					break;
				case AssignK:
					fprintf(listing,"Assign赋值: %s\n",tree->attr.name);
					break;
				case ReadK:
					fprintf(listing,"Read读: %s\n",tree->attr.name);
					break;
				case WriteK:
					fprintf(listing,"Write写\n");
					break;
				default:
					fprintf(listing,"未知语句类型\n");
					break;
			}
		}else if(tree->nodekind == ExpK) {
			switch (tree->kind.exp) {
				case OpK:
					fprintf(listing,"运算符: ");
					printToken(tree->attr.op, "\0");
					break;
				case ConstK:
					fprintf(listing,"常数: %s\n",tree->attr.val);
					break;
				case IdK:
					fprintf(listing,"标号: %s\n",tree->attr.name);
					break;
				default:
					fprintf(listing,"未知表达式类型\n");
					break;
			}
		}else {
			fprintf(listing,"未知结点类型\n");
		}
		for(i = 0; i < NELEMS(tree->child); i++)
			printTree(tree->child[i]);
		tree = tree->sibling;
	}
	indentno--;
}

#endif

