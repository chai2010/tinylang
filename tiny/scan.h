// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef SCAN_H
#define SCAN_H

#include "tiny.h"
#include "util.h"

TreeNode *stmt_sequence(void);
TreeNode *statement(void);
TreeNode *if_stmt(void);
TreeNode *repeat_stmt(void);
TreeNode *assign_stmt(void);
TreeNode *read_stmt(void);
TreeNode *write_stmt(void);
TreeNode *exp(void);
TreeNode *simple_exp(void);
TreeNode *term(void);
TreeNode *factor(void);

/* 记号状态 */
typedef enum {
	START, INASSIGN, INCOMMENT, INNUM, INID, DONE
} StateType;

/* 记号类型 */

char tokStr[32];
TokenType token;

/* 文件结束标志 */

int EOF_flag = FALSE;

/* 前进或倒退一个字符 */

int
nextChar(int flag)
{
	static int pos, size;
	static char buf[128];
	
	if(!flag) {
		if(!EOF_flag) pos--;
		return 0;
	};
	if(pos < size) return buf[pos++];
	if(!fgets(buf, NELEMS(buf), source)) {
		EOF_flag = TRUE;
		return EOF;
	}
	fprintf(listing, "%4d: %s", ++line, buf);
	size = strlen(buf);
	pos = 0;
	return buf[pos++];
}

/* 检验是否是保留字，返回记号类型 */

TokenType
tok_lookup(char * s)
{
	static char *str[] = { "if", "then", "else", "end",
		"repeat", "until", "read", "write" };
	static TokenType tok[] = { IF, THEN, ELSE, END,
		REPEAT, UNTIL, READ, WRITE };
	int i;
	for(i = 0; i < NELEMS(str); i++)
		if(!strcmp(s, str[i])) return tok[i];
	return ID;
}

/* 获取新的标号 */

TokenType
getToken(void)
{
	int tokIdx = 0;
	TokenType curTok;
	StateType state = START;
	int save;
	while (state != DONE) {
		int c = nextChar(TRUE);
		save = TRUE;
		switch(state) {
			case START:
				if(isdigit(c))
					state = INNUM;
				else if (isalpha(c))
					state = INID;
				else if (c == ':')
					state = INASSIGN;
				else if ((c == ' ') || (c == '\t') || (c == '\n'))
					save = FALSE;
				else if (c == '{') {
					save = FALSE;
					state = INCOMMENT;
				}else {
					state = DONE;
					switch(c) {
						case EOF:
							save = FALSE;
							curTok = ENDFILE;
							break;
						case '=':
							curTok = EQ;
							break;
						case '<':
							curTok = LT;
							break;
						case '+':
							curTok = PLUS;
							break;
						case '-':
							curTok = MINUS;
							break;
						case '*':
							curTok = TIMES;
							break;
						case '/':
							curTok = OVER;
							break;
						case '(':
							curTok = LPAREN;
							break;
						case ')':
							curTok = RPAREN;
							break;
						case ';':
							curTok = SEMI;
							break;
						default:
							curTok = ERROR;
							break;
					}
				}
				break;
			case INCOMMENT:
				save = FALSE;
				if(c == EOF) {
					state = DONE;
					curTok = ENDFILE;
				}else if(c == '}') state = START;
				break;
			case INASSIGN:
				state = DONE;
				if(c == '=')
					curTok = ASSIGN;
				else {
					nextChar(FALSE);
					save = FALSE;
					curTok = ERROR;
				}
				break;
			case INNUM:
				if (!isdigit(c)) {
					nextChar(FALSE);
					save = FALSE;
					state = DONE;
					curTok = NUM;
				}
				break;
			case INID:
				if(!isalpha(c)) {
					nextChar(FALSE);
					save = FALSE;
					state = DONE;
					curTok = ID;
				}
				break;
			case DONE:
			default:
				fprintf(listing,"扫描错误: 状态= %d\n",state);
				state = DONE;
				curTok = ERROR;
				break;
		}
		if(save && tokIdx <= NELEMS(tokStr))
			tokStr[tokIdx++] = (char)c;
		if (state == DONE) { tokStr[tokIdx] = '\0';
			if (curTok == ID) curTok = tok_lookup(tokStr);
		}
	}
	fprintf(listing, "\t%d: ", line);
	printToken(curTok, tokStr);
	return curTok;
}

/* 处理错误状态 */

void
syntaxError(char *msg)
{
	fprintf(listing, "\n语法错误: %d 行 %s", line, msg);
	Error = TRUE;
}

/* 若匹配则跳过，否则错误 */

void
match(TokenType expected)
{
	if (token != expected) {
		syntaxError("错误符号 -> ");
		printToken(token,tokStr);
		fprintf(listing, "\n");
	}else token = getToken();
}

/* 处理tiny序列 */

TreeNode *
stmt_sequence(void)
{
	TreeNode * t = statement();
	TreeNode * p = t;
	while(token != END && token != ELSE &&
		token != UNTIL && token != ENDFILE) {
		TreeNode * q;
		match(SEMI);
		q = statement();
		if(q != NULL) {
			if (t == NULL) { t = p = q; }
			else { p->sibling = q; p = q; }
		}
	}
	return t;
}

/* 五种类型的语句 */

TreeNode *
statement(void)
{
	TreeNode * t = NULL;
	switch (token) {
		case IF : t = if_stmt(); break;
		case REPEAT : t = repeat_stmt(); break;
		case ID : t = assign_stmt(); break;
		case READ : t = read_stmt(); break;
		case WRITE : t = write_stmt(); break;
		default : syntaxError("错误符号 -> ");
			printToken(token, tokStr);
			token = getToken();
			break;
	}
	return t;
}

TreeNode *
if_stmt(void)
{
	TreeNode * t = newTreeNode(StmtK, IfK);
	match(IF);
	if(t != NULL) t->child[0] = exp();
	match(THEN);
	if(t != NULL) t->child[1] = stmt_sequence();
	if(token == ELSE) {
		match(ELSE);
		if(t != NULL) t->child[2] = stmt_sequence();
	}
	match(END);
	return t;
}

TreeNode *
repeat_stmt(void)
{
	TreeNode * t = newTreeNode(StmtK, RepeatK);
	match(REPEAT);
	if(t != NULL) t->child[0] = stmt_sequence();
	match(UNTIL);
	if(t != NULL) t->child[1] = exp();
	return t;
}

TreeNode *
assign_stmt(void)
{
	TreeNode * t = newTreeNode(StmtK, AssignK);
	if((t!=NULL) && (token==ID))
		t->attr.name = str_new(tokStr);
	match(ID);
	match(ASSIGN);
	if(t!=NULL) t->child[0] = exp();
	return t;
}

TreeNode *
read_stmt(void)
{
	TreeNode * t = newTreeNode(StmtK, ReadK);
	match(READ);
	if((t!=NULL) && (token==ID))
		t->attr.name = str_new(tokStr);
	match(ID);
	return t;
}

TreeNode *
write_stmt(void)
{
	TreeNode * t = newTreeNode(StmtK, WriteK);
	match(WRITE);
	if(t!=NULL) t->child[0] = exp();
	return t;
}

/* 表达式类型 */

TreeNode *
exp(void)
{
	TreeNode * t = simple_exp();
	if((token == LT) || (token == EQ)) {
		TreeNode * p = newTreeNode(ExpK, OpK);
		if (p!=NULL) {
			p->child[0] = t;
			p->attr.op = token;
			t = p;
		}
		match(token);
		if(t!=NULL) t->child[1] = simple_exp();
	}
	return t;
}

TreeNode *
simple_exp(void)
{
	TreeNode * t = term();
	while((token == PLUS) || (token == MINUS)) {
		TreeNode * p = newTreeNode(ExpK, OpK);
		if(p !=NULL ) {
			p->child[0] = t;
			p->attr.op = token;
			t = p;
			match(token);
			t->child[1] = term();
		}
	}
	return t;
}

TreeNode *
term(void)
{
	TreeNode * t = factor();
	while ((token == TIMES)||(token == OVER)) {
		TreeNode * p = newTreeNode(ExpK, OpK);
		if(p != NULL) {
			p->child[0] = t;
			p->attr.op = token;
			t = p;
			match(token);
			p->child[1] = factor();
		}
	}
	return t;
}

TreeNode *
factor(void)
{
	TreeNode * t = NULL;
	switch(token) {
		case NUM :
			t = newTreeNode(ExpK, ConstK);
			if((t!=NULL) && (token==NUM))
				t->attr.val = str_new(tokStr);
			match(NUM);
			break;
		case ID :
			t = newTreeNode(ExpK, IdK);
			if ((t!=NULL) && (token==ID))
				t->attr.name = str_new(tokStr);
			match(ID);
			break;
		case LPAREN :
			match(LPAREN);
			t = exp();
			match(RPAREN);
			break;
		default:
			syntaxError("错误符号 -> ");
			printToken(token, tokStr);
			token = getToken();
			break;
	}
	return t;
}

/* 建立语法树 */

TreeNode *
buildTree(void)
{
	TreeNode * t;
	token = getToken();
	t = stmt_sequence();
	if(token != ENDFILE)
		syntaxError("代码结束于文件之前\n");
	fprintf(listing, "\n语法树:\n\n");
	printTree(t);
	return t;
}

#endif
