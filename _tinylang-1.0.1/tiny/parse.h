// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef PARSE_H
#define PARSE_H

#include "tiny.h"
#include "util.h"

char *str_new(char *str);
char *lab_lookup(char *key);

/* TINY程序符号表 */

struct label {
	struct label *link;
	char *key, *val;
	int line;
} *lab_buckets[211];

/* 加入新的标号 */

void
lab_insert(char * key, char *val, int line)
{
	struct label *p;
	int h;
	h = hash(key, NELEMS(lab_buckets));
	p = lab_buckets[h];
	while(p != NULL) {
		if(!strcmp(key, p->key)) break;
		p = p->link;
	}
	if(p == NULL) {
		p = malloc(sizeof(*p));
		if(p == NULL) RAISE("内存溢出");
		p->key = str_new(key);
		p->val = str_new(val);
		p->line = line;
		p->link = lab_buckets[h];
		lab_buckets[h] = p;
	}
}

/* 查询标号表 */

char *
lab_lookup(char *key)
{
	struct label *p;
	int h;
	h = hash(key, NELEMS(lab_buckets));
	p = lab_buckets[h];
	while(p != NULL) {
		if(!strcmp(key, p->key)) break;
		p = p->link;
	}
	return p ? p->val: NULL;
}

/* 释放标号表 */

void
lab_free(void)
{
	int i;
	for(i = 0; i < NELEMS(lab_buckets); ++i) {
		struct label * p, * q;
		q = p = lab_buckets[i];
		for(; p != NULL; p = q) {
			q = p->link; free(p);
		}
		lab_buckets[i] = NULL;
	}
}

/* 打印标号表 */

void
printSymTab(FILE * listing)
{
	int i;
	fprintf(listing,"变量名称  对应标号  初始行号\n");
	fprintf(listing,"--------  --------  --------\n");
	for(i = 0; i < NELEMS(lab_buckets); ++i) {
		struct label *p;
		for(p = lab_buckets[i]; p; p = p->link) {
			fprintf(listing, "%-8s  ", p->key);
			fprintf(listing, "%-8s  ", p->val);
			fprintf(listing, "%-8d\n",p->line);
		}
	}
}

/* 遍历语法树 */

void
traverse(TreeNode *t, void pre(TreeNode *), void post(TreeNode *))
{
	TreeNode *p;
	int i;
	if(t == NULL) return;
	if(pre != NULL) pre(t);
	for (i = 0; i < NELEMS(t->child); i++)
		traverse(t->child[i], pre, post);
	p = t->sibling;
	if(post != NULL) post(t);
	traverse(p, pre, post);
}

/* 插入新的结点 */

void
nod_insert(TreeNode *t)
{
	if(t->nodekind == StmtK) {
		StmtKind k = t->kind.stmt;
		if(k == AssignK || k == ReadK) {
			if(lab_lookup(t->attr.name) == NULL) {
				char *lab = new_label();
				lab_insert(t->attr.name, lab, t->line);
			}
		}
	}else if(t->nodekind == ExpK) {
		if(t->kind.exp == IdK) {
			if(lab_lookup(t->attr.name) == NULL) {
				char *lab = new_label();
				lab_insert(t->attr.name, lab, t->line);
			}
		}
	}
}

/* 类型错误 */
void
typeError(TreeNode *t, char *msg)
{
	fprintf(listing, "\n类型错误: %d 行 %s\n\n", t->line, msg);
	Error = TRUE;
}

/* 结点检查 */
void
nod_check(TreeNode * t)
{
	if(t->nodekind == ExpK) {
		ExpKind exp = t->kind.exp;
		if(exp == ConstK || exp == IdK) {
			t->type = Integer;
		}else if(exp == OpK) {
			if((t->child[0]->type != Integer) ||
				(t->child[1]->type != Integer))
				typeError(t, "操作符未使用整数变量");
			if ((t->attr.op == EQ) || (t->attr.op == LT))
				t->type = Boolean;
			else
				t->type = Integer;
		}
	}else if(t->nodekind == StmtK) {
		switch(t->kind.stmt) {
			case IfK:
				if(t->child[0]->type == Integer)
					typeError(t->child[0], "if判断语句未使用布尔变量");
				break;
			case AssignK:
				if(t->child[0]->type != Integer)
					typeError(t->child[0], "assig语句未使用整数变量");
				break;
			case WriteK:
				if(t->child[0]->type != Integer)
					typeError(t->child[0], "write语句未使用整数变量");
				break;
			case RepeatK:
				if(t->child[1]->type == Integer)
					typeError(t->child[1], "repeat判断语句未使用布尔变量");
				break;
			default: break;
		}

	}
}

/* 分析语法树 */

void
parseTree(TreeNode *tree)
{
	if(Error) return;
	traverse(tree, nod_insert, nod_check);
	fprintf(listing, "\n符号表：\n\n");
	printSymTab(listing);
}

#endif
