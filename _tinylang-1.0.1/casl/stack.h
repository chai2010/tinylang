// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef STACK_H
#define STACK_H

#include "casl.h"

typedef struct Stack_T {
	int size, i; int elem[];
} *Stack_T;

Stack_T
Stack_new(void)
{
	const static int size = 8;
	Stack_T stk = malloc(sizeof(*stk) + sizeof(int) * size);
	if(stk == NULL) RAISE("Stack_new内存分配失败");
	stk->size = size; stk->i = 0;
	return stk;
}

void
Stack_free(Stack_T stk)
{
	if(stk) free(stk);
}

Stack_T
Stack_push(Stack_T stk, int e)
{
	const static int inc = 4;
	if(stk == NULL) RAISE("Stack_push操作失败，空指针");
	stk->elem[stk->i++] = e;
	if(stk->i >= stk->size) {
		stk->size = stk->size + inc;
		stk = realloc(stk, sizeof(*stk) + sizeof(int)*stk->size);
		if(stk == NULL) RAISE("Stack_push内存分配失败");
	}
	return stk;
}

int
Stack_pop(Stack_T stk)
{
	if(stk == NULL)
		RAISE("Stack_pop操作失败，空指针");
	if(stk->i <= 0)
		RAISE("Stack_pop操作失败，栈已空");
	return stk->elem[--stk->i];
}

int
Stack_empty(Stack_T stk)
{
	if(stk == NULL)
		RAISE("Stack_empty操作失败，空指针");
	return (stk->i == 0);
}

#endif
