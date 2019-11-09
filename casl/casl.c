// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "casl.h"
#include "label.h"


off_t mem[MEMSIZE];	/* 64k 内存 */

char pgmName[32];	/* 汇编程序 */
char codName[32];	/* 机器代码 */

FILE *source, *code;	/* 文件指针 */

int pc = pc_start;	/* 指令地址 */
int line = 0;		/* 行计数器 */

int state = 0;		/* 状态标志 */
int Error = 0;		/* 错误标志 */

char buf[LINESIZE+2];	/* 行缓冲区	*/
int  len;		/* 字符串长度	*/
int  pos;		/* 行下标	*/

char tokStr[LINESIZE];	/* 标号字符串	*/
TokenType token;	/* 标号类型	*/

extern void getToken(void);

void
getLine(void)
{
	int len; line++;
	if(!fgets(buf, NELEMS(buf)-1, source)) {
		token = ENDLINE;
		return;
	}
	len = strlen(buf);
	if(buf[len-1] == '\n') buf[--len] = '\0';
	if(len > LINESIZE) QUIT("该行超出72个字符");
	pos = 0; getToken();
}

TokenType
tok_lookup(char *s)
{
	const static char *grs[] = {
		"GR0", "GR1", "GR2", "GR3", "GR4" };
	const static char *str[] = { 
		"HALT", "LD", "ST", "LEA",
		"ADD", "SUB", "MUL", "DIV", "MOD",
		"AND", "OR", "EOR",
		"SLA", "SRA", "SLL", "SRL",
		"CPA", "CPL",
		"JMP", "JPZ", "JMI", "JNZ", "JZE",
		"PUSH", "POP", "CALL", "RET",
		"READ", "WRITE", "IN", "OUT", "EXIT",
		"START", "END", "DC", "DS" };
	const static TokenType tok[] = {
		HALT, LD, ST, LEA,
		ADD, SUB, MUL, DIV, MOD,
		AND, OR, EOR,
		SLA, SRA, SLL, SRL,
		CPA, CPL,
		JMP, JPZ, JMI, JNZ, JZE,
		PUSH, POP, CALL, RET,
		READ, WRITE, IN, OUT, EXIT,
		START, END, DC, DS };
	int i;
	for(i = 0; i < NELEMS(grs); ++i) {
		if(!strcmp(s, grs[i])) {
			sprintf(tokStr, "%d", i);
			return NUM;
		}
	}
	for(i = 0; i < NELEMS(str); ++i) {
		if(!strcmp(s, str[i])) return tok[i];
	}
	return ID;
}

void
getToken(void)
{
	static int start = 0;
	static int string = 1;
	static int num = 2;
	static int id = 3;
	static int done = 4;
	
	int flag = start;
	int tokIdx = 0;
	
	while(flag != done) {
		int c = buf[pos++];
		if(flag == start) {
			if(c == '\0' || c == ';') {
				buf[--pos] = '\0';
				tokStr[0] = '\0';
				token = ENDLINE;
				flag = done;
			}else if(c == ',') {
				token = COMMA;
				flag = done;
			}else if(isupper(c)) {
				tokStr[tokIdx++] = (char)c;
				flag = id;
			}else if(c == '\'') {
				flag = string;
				len = 0;
			}else if(c == '+' || c == '-'
				|| c == '#' || isdigit(c)) {
				if(c == '#') {
					tokStr[tokIdx++] = '0';
					tokStr[tokIdx++] = 'x';
				}else {
					tokStr[tokIdx++] = (char)c;
				}
				flag = num;
			}else if(!isspace(c)) QUIT("未知记号");
		}else if(flag == id) {
			if(!isdigit(c) && !isupper(c)) {
				tokStr[tokIdx] = '\0';
				token = tok_lookup(tokStr);
				flag = done;
				pos--;
			}else tokStr[tokIdx++] = (char)c;
		}else if(flag == num) {
			if(!isdigit(c)) {
				tokStr[tokIdx] = '\0';
				token = NUM;
				flag = done;
				pos--;
			}else tokStr[tokIdx++] = (char)c;
		}else if(flag == string) {
			if(c == '\'') {
				token = STRING;
				flag = done;
			}else if(c == '\\') {
				c = buf[pos++];
				if(c == '\0') QUIT("字符串定义错误");
				if(c == '\'') c = '\'';
				else if(c == '0') c = '\0';
				else if(c == 't') c = '\t';
				else if(c == 'n') c = '\n';
				tokStr[tokIdx++] = (char)c;
				len++;
			}else if(c == '\0') {
				QUIT("字符串定义错误");
			}else {
				tokStr[tokIdx++] = (char)c;
				len++;
			}	/* string 判断 */
		}
		if(flag != string &&
			tokIdx > WORDSIZE) {
			QUIT("记号多于6个字符");
		}
	}	/* while 循环结束 */
}

void
skipLabel(void)
{
	if(isspace(buf[0])) return;
	if(token != ID) QUIT("非法标号");
	lab_define(tokStr, pc);
	getToken();
}

void
skipOp(void)
{
	short op = token;
	if(op < HALT || op > DS) QUIT("未知指令");
	mem[pc] = op << 8;
	getToken();
}

void
skipGR(void)
{
	short gr;
	if(token != NUM) QUIT("缺少GR");
	gr = (short)atoi(tokStr);
	if(gr < 0 || gr > 4) QUIT("GR错误");
	mem[pc] |= (gr<<4);
	getToken();
}

void
skipADR()
{
	short adr;
	if(token == NUM) {
		adr = (short)atoi(tokStr);
	}else if(token == ID) {
		adr = lab_get(tokStr, pc+1);
	}else QUIT("ADR错误");
	mem[pc+1] = adr;
	getToken();
}

void
skipXR()
{
	short xr;
	if(token != NUM) QUIT("XR错误");
	xr = (short)atoi(tokStr);
	if(xr < 1 || xr > 4) QUIT("XR错误");
	mem[pc] |= xr;
	getToken();
}

void
skipCOMMA(void)
{
	if(token != COMMA) QUIT("逗号不匹配");
	getToken();
}

void
macro_in(void)
{
	static off_t cmd[] = {
		ST   << 8, ac_comet,
		PUSH << 8, ac_comet,
		LEA  << 8, 0,
		ST   << 8, IO_ADDR,
		LEA  << 8, IO_MAX,
		ST   << 8, ac_comet,
		LD   << 8, 0,
		AND  << 8, ac_comet,
		ST   << 8, ac_comet,
		PUSH << 8, ac_comet,
		LEA  << 8, IO_CHR | IO_IN,
		ST   << 8, ac_comet,
		POP  << 8,
		OR   << 8, ac_comet,
		ST   << 8, IO_FLAG,
		POP  << 8 };
	const short ai = 5, ni = 13;
	int i;
	if(token == ID)
		cmd[ai] = lab_get(tokStr, pc+ai);
	else if(token == NUM)
		cmd[ai] = (short)atoi(tokStr);
	else
		QUIT("READ参数错误");
	getToken(); skipCOMMA();
	if(token == ID)
		cmd[ni] = lab_get(tokStr, pc+ai);
	else if(token == NUM)
		cmd[ni] = (short)atoi(tokStr);
	else
		QUIT("READ参数错误");
	for(i = 0; i < NELEMS(cmd); ++i)
		mem[pc++] = cmd[i];
	getToken();
}

void
macro_out(void)
{
	static off_t cmd[] = {
		ST   << 8, ac_comet,
		PUSH << 8, ac_comet,
		LEA  << 8, 0,
		ST   << 8, IO_ADDR,
		LEA  << 8, IO_MAX,
		ST   << 8, ac_comet,
		LD   << 8, 0,
		AND  << 8, ac_comet,
		ST   << 8, ac_comet,
		PUSH << 8, ac_comet,
		LEA  << 8, IO_CHR | IO_OUT,
		ST   << 8, ac_comet,
		POP  << 8,
		OR   << 8, ac_comet,
		ST   << 8, IO_FLAG,
		POP  << 8 };
	const short ai = 5, ni = 13;
	int i;
	if(token == ID)
		cmd[ai] = lab_get(tokStr, pc+ai);
	else if(token == NUM)
		cmd[ai] = (short)atoi(tokStr);
	else
		QUIT("READ参数错误");
	getToken(); skipCOMMA();
	if(token == ID)
		cmd[ni] = lab_get(tokStr, pc+ai);
	else if(token == NUM)
		cmd[ni] = (short)atoi(tokStr);
	else
		QUIT("READ参数错误");
	for(i = 0; i < NELEMS(cmd); ++i)
		mem[pc++] = cmd[i];
	getToken();
}

void
macro_exit(void)
{
	mem[pc++] = (short)(HALT << 8);
}

void
macro_read(void)
{
	static off_t cmd[] = {
		ST   << 8, ac_comet,
		PUSH << 8, ac_comet,
		LEA  << 8, 0,
		ST   << 8, IO_ADDR,
		LEA  << 8, (1 & IO_MAX) | IO_DEC | IO_IN,
		ST   << 8, IO_FLAG,
		POP  << 8 };
	const short ai = 5;
	int i;

	if(token == ID)
		cmd[ai] = lab_get(tokStr, pc+ai);
	else if(token == NUM)
		cmd[ai] = (short)atoi(tokStr);
	else
		QUIT("READ参数错误");
	for(i = 0; i < NELEMS(cmd); ++i)
		mem[pc++] = cmd[i];
	getToken();
}

void
macro_write(void)
{
	static off_t cmd[] = {
		ST   << 8, ac_comet,
		PUSH << 8, ac_comet,
		LEA  << 8, 0,
		ST   << 8, IO_ADDR,
		LEA  << 8, (1 & IO_MAX) | IO_DEC | IO_OUT,
		ST   << 8, IO_FLAG,
		POP  << 8 };
	const short ai = 5;
	int i;
	if(token == ID)
		cmd[ai] = lab_get(tokStr, pc+ai);
	else if(token == NUM)
		cmd[ai] = (short)atoi(tokStr);
	else
		QUIT("READ参数错误");
	for(i = 0; i < NELEMS(cmd); ++i)
		mem[pc++] = cmd[i];
	getToken();
}

void
macro_start(void)
{
	if(state != 0) QUIT("START指令重复");
	mem[pc] = (short)(JMP << 8);
	if(token != ENDLINE) skipADR();
	else mem[pc + 1] = pc + 2;
	state = START;
	pc += 2;
}

void
macro_end(void)
{
	if(state != START) QUIT("缺少START指令");
	state = END;
}

void
macro_dc(void)
{
	short adr;
	if(token == NUM) {
		adr = (short)atoi(tokStr);
		mem[pc++] = adr;
	}else if(token == ID) {
		adr = lab_get(tokStr, pc);
		mem[pc++] = adr;
	}else if(token == STRING) {
		int i = 0;
		while(i < len) {
			mem[pc++] = tokStr[i++];
		}
	}else QUIT("DC 参数错误");
	getToken();
}

void
macro_ds(void)
{
	int num;
	if(token != NUM) QUIT("DS 参数错误");
	num = atoi(tokStr);
	if(num < 0 || num > pc_max) QUIT("DS 参数错误");
	getToken();
	pc += num;
}

void
buildCode(void)
{
	while(!feof(source)) {
		int op;
		getLine();
		if(token == ENDLINE) continue;
		if(pc > pc_max) QUIT("程序太大");
		skipLabel();
		op = token;
		skipOp();
		switch(op) {
			/* 两个字长的指令 */
			case LD: case ST: case LEA:
			case ADD: case SUB: case MUL: case DIV: case MOD:
			case AND: case OR: case EOR:
			case SLA: case SRA: case SLL: case SRL:
			case CPA: case CPL:
				skipGR(); skipCOMMA(); skipADR();
				if(token != ENDLINE) {
					skipCOMMA(); skipXR();
				}
				pc += 2; break;
			
			case JMP: case JPZ: case JMI: case JNZ: case JZE:
			case PUSH: case CALL:
				skipADR();
				if(token != ENDLINE) {
					skipCOMMA(); skipXR();
				}
				pc += 2; break;
			
			case POP:case RET: case HALT:
				if(op == POP) skipGR();
				pc += 1; break;
				
			case IN: macro_in(); break;
			case OUT: macro_out(); break;
			case READ: macro_read(); break;
			case WRITE: macro_write(); break;
			case EXIT: macro_exit(); break;
			case DC: macro_dc(); break;
			case DS: macro_ds(); break;
			case START: macro_start(); break;
			case END: macro_end(); break;
		
			default: QUIT("未知操作");
		}
		if(token != ENDLINE) QUIT("指令错误");
	} // while 循环结束
	if(state != END) QUIT("缺少END指令");
}

void
init(int n, char *v[])
{
	int len;
	char *s;

	printf("==================\n");
	printf("CASL汇编语言编译器\n");
	printf("==================\n\n");
	if(n != 2) QUIT("命令 <文件>");
	len = strlen(v[1]);
	if(len > 16) QUIT("文件名太长");
	strcpy(pgmName, v[1]);
	s = strchr(pgmName, '.');
	if(s != NULL) {
		if(strcmp(s, ".casl")) QUIT("不是casl程序");
		else *s = '\0';
	}
	strcpy(codName, pgmName);
	strcat(pgmName, ".casl");
	strcat(codName, ".comet");
	source = fopen(pgmName, "r");
	if(source == NULL) QUIT("CASL程序不能打开");
}

void
chkLab(Label_T lab, void *cl)
{
	char msg[32];
	if(lab == NULL || lab->datoff == NULL) return;
	sprintf(msg, "错误信息：%-4d行 %s 标号没有定义\n", lab->addr, lab->key);
	QUIT(msg);
}

void
casl_free(void)
{
	off_t tmp[2];
	lab_map(chkLab, NULL);
	code = fopen(codName, "wb");
	if(code == NULL) QUIT("目标文件不能打开");
	tmp[0] = (off_t)pc_start;
	tmp[1] = (off_t)(pc -  pc_start);
	fwrite(tmp, sizeof(off_t), NELEMS(tmp), code);
	fwrite(&mem[tmp[0]], sizeof(off_t), tmp[1], code);
	printf("输入文件 %s\n", pgmName);
	printf("输出文件 %s\n", codName);
	fclose(source);
	fclose(code);
	lab_free();
}

int
caslMain(int n, char *v[])
{
	init(n, v);
	buildCode();
	casl_free();
	return 0;
}

