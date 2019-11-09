// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef UTIL_H
#define UTIL_H

#include "tiny.h"

// 字符串hash化
int hash(char *key, int M);

// 生成新的字符串原子
char *str_new(char *s);


// 释放所有原子
void str_free(void);

// 生成新的CASL标号

char* new_label(void);

// 分配一个语句类型的结点
TreeNode* newTreeNode(NodeKind nodeKind, int kind);

// 打印记号

void printToken(TokenType token, const char *str);

// 打印语法树

void printTree(TreeNode *tree);

#endif // UTIL_H
