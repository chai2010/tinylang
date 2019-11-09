// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "comet.h"

/* 装载comet程序 */

void
comet_load(void)
{
	off_t n, tmp[2];
	fseek(source, 0, SEEK_SET);
	n = fread(tmp, sizeof(off_t), 2, source);
	if(n != 2) {
		printf("文件 %s 格式错误\n", pgmName);
		exit(1);
	}
	n = fread(&cmt.mem[tmp[0]], sizeof(off_t), tmp[1], source);
	if(n != tmp[1]) {
		printf("文件 %s 格式错误\n", pgmName);
		exit(1);
	}
	cmt.gr[4] = (off_t)sp_start;
	cmt.pc = tmp[0];
}

/* 处理comet计算机的IO设备 */

void
comet_io(void)
{
	int i;
	off_t addr;
	short count, type, fio;
	char *fmt;
	count = cmt.mem[IO_FLAG] & IO_MAX;
	if(count == 0) return;
	fio = cmt.mem[IO_FLAG] & IO_FIO;
	type = cmt.mem[IO_FLAG] & IO_TYPE;
	addr = cmt.mem[IO_ADDR];
	if(type == IO_CHR) fmt = "%c";
	else if(type == IO_OCT) fmt = "%o";
	else if(type == IO_DEC) fmt = "%d";
	else if(type == IO_HEX) fmt = "%x";
	else {
		cmt.mem[IO_FLAG] &= (!IO_MAX);
		cmt.mem[IO_FLAG] |= IO_ERROR;
		return;
	}
	for(i = 0; i < count; ++i) {
		if(fio == IO_IN)
			scanf(fmt, &cmt.mem[addr++]);
		else if(fio == IO_OUT)
			printf(fmt, cmt.mem[addr++]);
	}
	cmt.mem[IO_FLAG] &= (!IO_MAX);
}

/* 单步运行comet计算机 */

int
comet_step(void)
{
	off_t temp;
	off_t adr, x1, x2;
	short op, gr, xr;
	
	op = (off_t)cmt.mem[cmt.pc] / 0x100;
	gr = (off_t)cmt.mem[cmt.pc] % 0x100 / 0x10;
	xr = (off_t)cmt.mem[cmt.pc] % 0x10;
	adr = (off_t)cmt.mem[cmt.pc + 1];
	
	if(gr < 0 || gr > 4) {
		temp = cmt.mem[cmt.pc];
		printf("非法指令：mem[%x] = %x\n", cmt.pc, temp);
		return 0;
	}
	if(xr < 0 || xr > 4) {
		temp = cmt.mem[cmt.pc];
		printf("非法指令：mem[%x] = %x\n", cmt.pc, temp);
		return 0;
	}
	if(xr != 0) adr += cmt.gr[xr];

	comet_io();
	
	switch(op) {
		case HALT:
			cmt.pc += 1;
			return 0;
		case LD:
			cmt.pc += 2;
			cmt.gr[gr] = cmt.mem[adr];
			return 1;
		case ST:
			cmt.pc += 2;
			cmt.mem[adr] = cmt.gr[gr];
			return 1;
			
		case LEA:
			cmt.pc += 2;
			cmt.gr[gr] = adr;
			cmt.fr = cmt.gr[gr];
			return 1;
		case ADD:
			cmt.pc += 2;
			cmt.gr[gr] += cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case SUB:
			cmt.pc += 2;
			cmt.gr[gr] -= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case MUL:
			cmt.pc += 2;
			cmt.gr[gr] *= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case DIV:
			cmt.pc += 2;
			cmt.gr[gr] /= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case MOD:
			cmt.pc += 2;
			cmt.gr[gr] %= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case AND:
			cmt.pc += 2;
			cmt.gr[gr] &= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case OR :
			cmt.pc += 2;
			cmt.gr[gr] |= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case EOR:
			cmt.pc += 2;
			cmt.gr[gr] ^= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		
		case SLA:
			cmt.pc += 2;
			cmt.gr[gr] <<= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		case SRA:
			cmt.pc += 2;
			cmt.gr[gr] >>= cmt.mem[adr];
			cmt.fr = cmt.gr[gr];
			return 1;
		
		case SLL:
			cmt.pc += 2;
			x1 = cmt.mem[gr];
			x1 <<= cmt.mem[adr];
			cmt.mem[gr] = cmt.fr = x1;
			return 1;
		case SRL:
			cmt.pc += 2;
			x1 = cmt.mem[gr];
			x1 >>= cmt.mem[adr];
			cmt.mem[gr] = cmt.fr = x1;
			return 1;
		
		case CPA:
			cmt.pc += 2;
			cmt.fr = cmt.gr[gr] - cmt.mem[adr];
			return 1;
		case CPL:
			cmt.pc += 2;
			x1 = cmt.gr[gr];
			x2 = cmt.mem[adr];
			cmt.fr = x1 - x2;
			return 1;
		
		case JMP:
			cmt.pc += 2;
			cmt.pc = adr;
			return 1;
		case JPZ:
			cmt.pc += 2;
			if(cmt.fr >= 0) cmt.pc = adr;
			return 1;
		case JMI:
			cmt.pc += 2;
			if(cmt.fr < 0) cmt.pc = adr;
			return 1;
		case JNZ:
			cmt.pc += 2;
			if(cmt.fr != 0) cmt.pc = adr;
			return 1;
		case JZE:
			cmt.pc += 2;
			if(cmt.fr == 0) cmt.pc = adr;
			return 1;
		
		case PUSH:
			cmt.pc += 2;
			x1 = --cmt.gr[4];
			cmt.mem[x1] = cmt.mem[adr];
			return 1;
		case POP:
			cmt.pc += 1;
			x1 = cmt.gr[4]++;
			cmt.gr[gr] = cmt.mem[x1];
			return 1;
		
		case CALL:
			cmt.pc += 2;
			x1 = --cmt.gr[4];
			cmt.mem[x1] = cmt.pc;
			cmt.pc = cmt.mem[adr];
			return 1;
		case RET:
			cmt.pc += 1;
			x1 = cmt.gr[4]++;
			cmt.pc = cmt.mem[x1];
			return 1;
			
		default : 
			temp = cmt.mem[cmt.pc];
			printf("非法指令：mem[%x] = %x\n", cmt.pc, temp);
			return 0;
	}
}

/* 输出指令，用于调试 */

void
writeIns(off_t pc, off_t n)
{
	off_t op, gr, adr, xr;
	off_t i;

	for(i = 0; i < n; ++i) {
		op = cmt.mem[pc] / 0x100;
		gr = cmt.mem[pc] % 0x100 / 0x10;
		xr = cmt.mem[pc] % 0x10;
		adr = cmt.mem[pc + 1];
		if(op > RET) {
			printf("mem[%-4x]: 未知\n", pc);
			return;
		}
		if(gr < 0 || gr > 4) {
			printf("mem[%-4x]: 未知\n", pc);
			return;
		}
		if(xr < 0 || xr > 4) {
			printf("mem[%-4x]: 未知\n", pc);
			return;
		}
		printf("mem[%-4x]: %s\t", pc, opTab[op].str);
		if(op == HALT || op == RET) {
			printf("\n");
			pc += 1;
			continue;
		}else if(op == POP) {
			printf("GR%d\n", gr);
			pc += 1;
			continue;
		}if(op < CPL) {
			printf("GR%d, %x", gr, adr);
			pc += 2;
		}else {
			printf("%x", adr);
			pc += 2;
		}
		if(xr != 0) printf(", GR%d", xr);
		printf("\n");
	}
}

/* 调试comet计算机程序 */

void
comet_debug(void)
{
	off_t stepcnt = 0;
	off_t pntflag = 0;
	off_t traflag = 0;
	
	char buf[32], s[32];
	off_t x1, x2;
	int i, n, cmd;
	
	printf("调试 （帮助输入 help）...\n\n");
	
LOOP:
	
	do {
		fflush(stdin);
		printf ("输入命令: ");
		fgets(buf, NELEMS(buf), stdin);
		n = sscanf(buf, "%s %x %x", s, &x1, &x2);
		for(i = 0, cmd = -1; i < NELEMS(dbTab); ++i) {
			if(!strcmp(dbTab[i].s1, s) ||
				!strcmp(dbTab[i].s2, s)) {
				cmd = dbTab[i].db;
				break;
			}
		}
	}while(n <= 0);

	switch(cmd) {
		case HELP:
			printf("命令列表:\n");
			printf("  h)elp           显示本命令列表\n");
			printf("  g)o             运行程序直到停止\n");
			printf("  s)tep  <n>      执行 n 条指令 （默认为 1 ）\n");
			printf("  j)ump  <b>      跳转到 b 地址 （默认为当前地址）\n");
			printf("  r)egs           显示寄存器内容\n");
			printf("  i)Mem  <b <n>>  显示从 b 开始 n 个内存数据\n");
			printf("  d)Mem  <b <n>>  显示从 b 开始 n 个内存指令\n");
			printf("  a(lter <b <v>>  修改 b 位置的内存数据为 v 值\n");
			printf("  t)race          开关指令显示功能\n");
			printf("  p)rint          开关指令计数功能\n");
			printf("  c)lear          重置模拟器内容\n");
			printf("  q)uit           终止模拟器\n");
			break;
			
		case GO:
			stepcnt = 0;
			do { stepcnt++;
				if(traflag) writeIns(cmt.pc, 1);
			}while(comet_step());
			if(pntflag)
				printf("执行指令数目 = %d\n", stepcnt);
			break;
			
		case STEP:
			if(n >= 2) stepcnt = x1;
			else stepcnt = 1;
			for(i = 0; i < stepcnt; ++i) {
				if(traflag) writeIns(cmt.pc, 1);
				if(!comet_step()) break;
			}
			if(pntflag)
				printf("执行指令数目 = %d\n", i);
			break;

		case JUMP:
			if(n < 2) x1 = cmt.pc;
			cmt.pc = x1;
			printf("指令跳转到 %x\n", x1);
			break;

		case REGS:
			printf("显示寄存器数据\n");
			printf("GR[0] = %4x\tPC = %4x\n",
				(off_t)cmt.gr[0], cmt.pc);
			printf("GR[1] = %4x\tSP = %4x\n",
				(off_t)cmt.gr[1], (off_t)cmt.gr[4]);
			printf("GR[2] = %4x\t", (off_t)cmt.gr[2]);
			if(cmt.fr > 0) printf("FR =   00\n");
			else if(cmt.fr < 0) printf("FR =   10\n");
			else printf("FR =   01\n");
			printf("GR[3] = %4x\n", (off_t)cmt.gr[3]);
			break;

		case IMEM:
			printf("显示内存指令\n");
			if(n < 2) x1 = cmt.pc;
			if(n < 3) x2 = 1;
			writeIns(x1, x2);
			break;


		case DMEM:
			printf("显示内存数据\n");
			if(n < 2) x1 = cmt.pc;
			if(n < 3) x2 = 1;
			if(x2 < 0) printf("参数错误\n");
			for(i = 0; i < x2; ++i) {
				off_t temp;
				temp = cmt.mem[x1];
				printf("mem[%-4x] = %x\n", x1, temp);
				x1++;
			}
			break;

		case ALTER:
			printf("修改内存数据 ");
			if(n == 3) {
				printf(" mem[%x] = %x\n", x1, x2);
				cmt.mem[x1] = x2;
			}else printf("失败！\n");
			break;

		case TRACE:
			traflag = !traflag;
			printf("指令显示功能");
			if(traflag) printf("打开\n");
			else printf("关闭\n");
			break;

		case PRINT:
			pntflag = !pntflag;
			printf("指令计数功能");
			if(pntflag) printf("打开\n");
			else printf("关闭\n");
			break;

		case CLEAR:
			printf("程序重新载入内存\n");
			comet_load();
			break;

		case QUIT:
			printf("退出调试...\n");
			return;
			
		default :
			printf("未知命令 %s\n", s);
			break;
	}
	goto LOOP;
}

/* 初始化相关参数 */

void
init(int n, char*v[])
{
	int len;
	char *s;

	printf("===============\n");
	printf("COMET虚拟计算机\n");
	printf("===============\n\n");
	if(n != 2 && n != 3) {
		printf("用法: %s [-d(ebug] <文件名>\n", v[0]);
		exit(1);
	}
	len = strlen(v[n-1]);
	if(len > 16) {
		printf("文件名太长");
		exit(1);
	}
	strcpy(pgmName, v[n-1]);
	s = strchr(pgmName, '.');
	if(s != NULL) {
		if(strcmp(s, ".comet")) {
			printf("%s 不是comet程序\n", pgmName);
			exit(1);
		} else {
			*s = '\0';
		}
	}
	strcat(pgmName, ".comet");
	if(n == 3) {
		if(strcmp(v[n-2], "-d") &&
			strcmp(v[n-2], "-debug")) {
			printf("用法: %s [-d(ebug] <文件名>\n", v[0]);
			exit(1);
		}
		debug = 1;
	}
	source = fopen(pgmName, "rb");
	if(source == NULL) {
		printf("%s 文件不能打开\n", pgmName);
		exit(1);
	}
	comet_load();
}

/* 主函数 */

int
main(int argc, char *argv[])
{
	init(argc, argv);
	if(debug) comet_debug();
	else while(comet_step());
	fclose(source);
	return 0;
}
