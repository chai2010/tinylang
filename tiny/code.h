// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CODE_H
#define CODE_H

#include "tiny.h"

void codeGen(TreeNode *);

/* 生成单行注释 */

void
emitComment(char *msg)
{
	fprintf(code, "; %s", msg);
}

/* 生成CASL指令 */

void
emitCASL(char *id, char *op, char *p1, char *p2, char *p3)
{
	fprintf(code, "\n");
	if(id) fprintf(code, "%s", id);
	fprintf(code, "\t");
	if(op) fprintf(code, "%s", op);
	fprintf(code, "\t");
	if(p1) fprintf(code, "%s", p1);
	if(p2) fprintf(code, ",\t%s", p2);
	else fprintf(code, "\t");
	if(p3) fprintf(code, ",\t%s", p3);
	else fprintf(code, "\t");
	fprintf(code, "\t");
}

/* 五种tiny控制语句 */

void
genStmt(TreeNode * tree)
{
	TreeNode *p1, *p2, *p3;
	char *jmp, *loc1, *loc2;
	
	switch(tree->kind.stmt) {
		case IfK :
			emitCASL(0, 0, 0, 0, 0);
			emitComment("if 判断语句开始");
			p1 = tree->child[0];
			p2 = tree->child[1];
			p3 = tree->child[2];
			jmp = (p1->attr.op == LT)? "JPZ" : "JNZ";
			loc1 = new_label(); loc2 = new_label();
			codeGen(p1);
			emitCASL(0, jmp, loc1, 0, 0);
			emitComment("if 语句，跳转到else部分");
			codeGen(p2);
			emitCASL(0, "JMP", loc2, 0, 0);
			emitComment("if 语句，跳转到end部分");
			emitCASL(loc1, "DS", "0", 0, 0);
			emitComment("对应 if 语句的 else 地址");
			codeGen(p3);
			emitCASL(loc2, "DS", "0", 0, 0);
			emitComment("对应 if 语句的 end 地址");
			emitCASL(0, 0, 0, 0, 0);
			emitComment("if 语句结束");
			break;
		
		case RepeatK:
			emitCASL(0, 0, 0, 0, 0);
			emitComment("repeat 循环语句开始");
			p1 = tree->child[0];
			p2 = tree->child[1];
			jmp = (p2->attr.op == LT)? "JPZ" : "JNZ";
			loc1 = new_label();
			emitCASL(loc1, "DS", "0", 0, 0);
			emitComment("对应 repeat 语句开始地址");
			codeGen(p1);
			codeGen(p2);
			emitCASL(0, jmp, loc1, 0, 0);
			emitComment("跳转到 repeat 循环语句开始");
			emitCASL(0, 0, 0, 0, 0);
			emitComment("repeat 循环语句结束");
			break;
		
		case AssignK:
			codeGen(tree->child[0]);
			loc1 = lab_lookup(tree->attr.name);
			emitCASL(0, "ST", "GR0", loc1, 0);
			emitComment("assign 语句");
			break;
		
		case ReadK:
			loc1 = lab_lookup(tree->attr.name);
			emitCASL(0, "READ", loc1, 0, 0);
			emitComment("read 语句");
			break;
		
		case WriteK:
			codeGen(tree->child[0]);
			emitCASL(0, "ST", "GR0", "AC", 0);
			emitComment("保存 write 值在 AC 中");
			emitCASL(0, "WRITE", "AC", 0, 0);
			emitComment("write 语句，输出 AC 值");
			break;
		
		default: break;
	}
}

/* 三种表达式 */
void
genExp(TreeNode * tree)
{
	char *loc;
	TreeNode *p1, *p2;
	switch(tree->kind.exp) {
		case ConstK :
			emitCASL(0, "LEA", "GR0", tree->attr.val, 0);
			emitComment("对应常量");
			break;
		
		case IdK :
			loc = lab_lookup(tree->attr.name);
			emitCASL(0, "LD", "GR0", loc, 0);
			emitComment("对应变量");
			break;
		
		case OpK :
			p1 = tree->child[0];
			p2 = tree->child[1];
			codeGen(p2);
			emitComment("运算符右边的值");
			emitCASL(0, "ST", "GR0", "AC", 0);
			emitComment("保存运算符右边的值");
			codeGen(p1);
			emitComment("运算符左边的值");
			switch(tree->attr.op) {
				case PLUS :
					emitCASL(0, "ADD", "GR0", "AC", 0);
					break;
				case MINUS :
					emitCASL(0, "SUB", "GR0", "AC", 0);
					break;
				case TIMES :
					emitCASL(0, "MUL", "GR0", "AC", 0);
					break;
				case OVER :
					emitCASL(0, "DIV", "GR0", "AC", 0);
					break;
				case LT :
				case EQ :
					emitCASL(0, "CPA", "GR0", "AC", 0);
					break;
				default: break;
			}
			emitComment("计算表达式的值");
			break;
		
		default: break;
	}
}

/* 生成机器码 */

void
codeGen(TreeNode *tree)
{
	if(tree == NULL) return;
	if(tree->nodekind == StmtK) genStmt(tree);
	else if(tree->nodekind == ExpK) genExp(tree);
	codeGen(tree->sibling);
}

/* 生成CASL程序 */

void
buildCode(TreeNode *tree)
{
	int i;
	if(Error) return;
	emitComment("============\n");
	emitComment("CASL汇编程序\n");
	emitComment("============\n");
	emitCASL(0, "START", "CASL00", 0, 0);
	emitComment("程序入口");
	emitCASL("AC    ", "DS", "1", 0, 0);
	emitComment("AC 存放临时变量");
	for(i = 0; i < NELEMS(lab_buckets); ++i) {
		struct label *p = lab_buckets[i];
		while(p != NULL) {
			char msg[20];
			emitCASL(p->val, "DS", "1", 0, 0);
			sprintf(msg, "对应变量 %s", p->key);
			emitComment(msg);
			p = p->link;
		}
	}
	emitCASL("CASL00", "DS", "0", 0, 0);
	emitComment("CASL00 为程序的启动地址");
	codeGen(tree);
	emitCASL(0, "HALT", 0, 0, 0);
	emitComment("停机");
	emitCASL(0, "END", 0, 0, 0);
	emitComment("程序结束");
}

#endif
