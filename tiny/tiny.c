// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "tiny.h"
#include "scan.h"
#include "parse.h"
#include "code.h"
#include "util.h"

char pgmName[30];	/* TINY源文件名称	*/
char lstName[30];	/* 输出的信息文件名称	*/
char codName[30];	/* 产生的汇编代码文件	*/

FILE * source;		/* TINY源代码文件指针	*/
FILE * listing;		/* 信息文件指针		*/
FILE * code;		/* 汇编代码文件指针	*/

int line  = 0;		/* 对应源文件行号	*/
int Error = FALSE;	/* 错误标志，一般未使用	*/

/* 初始化TINY编译器的相关信息 */

void
init(int n, char *v[])
{
	char *s;
	int len;

	printf("=====================\n");
	printf("TINY编译器 到CASL语言\n");
	printf("=====================\n\n");

	if(n != 2) {
		fprintf(stderr, "用法: %s <文件名>\n", v[0]);
		exit(1);
	}
	len = strlen(v[1]);
	if(len > 16) {
		fprintf(stderr, "文件 %s 名称太长\n", v[1]);
		exit(1);
	}
	strcpy(pgmName, v[1]);
	s = strchr(pgmName, '.');
	if(s != NULL) {
		if(strcmp(s, ".tiny") != 0) {
			fprintf(stderr, "文件 %s 不是TINY程序\n", pgmName);
			exit(1);
		} else {
			*s = '\0';
		}
	}
	strcpy(lstName, pgmName);
	strcpy(codName, pgmName);
	strcat(pgmName, ".tiny");
	strcat(lstName, ".list");
	strcat(codName, ".casl");
	
	source = fopen(pgmName, "r");
	if(source == NULL) {
		fprintf(stderr, "无法打开文件 %s\n", pgmName);
		fcloseall();
		exit(1);
	}
	code = fopen(codName, "w");
	if(code == NULL) {
		fprintf(stderr, "无法建立文件 %s\n", codName);
		fcloseall();
		exit(1);
	}
	listing = fopen(lstName, "w");
	if(listing == NULL) {
		fprintf(stderr, "无法建立文件 %s\n", lstName);
		fcloseall();
		exit(1);
	}
	fprintf(listing, "\nTINY编译器: %s\n\n", pgmName);
	
	printf("编译文件 %s\n\n", pgmName);
	printf("编译中...\n\n");
}

/* 释放系统申请的资源 */

void
tiny_free(TreeNode *tree)
{
	traverse(tree, NULL, free);
	lab_free(); str_free();
	if(Error) printf("未知错误 :(\n");
	else printf("分析结束 :)\n");
	fcloseall();
}

/* 主函数 */

int
main(int argc, char *argv[])
{
	TreeNode *tree;		/* 语法树 */
	
	init(argc, argv);	/* 初始化相关信息	*/
	tree = buildTree();	/* 建立TINY程序的语法树	*/
	parseTree(tree);	/* 类型分析和生成符号表	*/
	buildCode(tree);	/* 生成CASL汇编代码	*/
	
	tiny_free(tree);	/* 释放系统申请的资源	*/
	return 0;
}

