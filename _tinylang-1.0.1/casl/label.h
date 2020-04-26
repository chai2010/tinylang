// Copyright 2005 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef LABEL_H
#define LABEL_H

#include "casl.h"
#include "stack.h"

typedef struct label {
	struct label *link;
	char *key;
	int   addr;
	Stack_T datoff;
} *Label_T;

Label_T buckets[211];

int
hash(const char *key)
{
	int h = 0, a = 127, i;
	for(i = 0; key[i]; ++i)
		h = (a * h + key[i]) % NELEMS(buckets);
	return h;
}

void
lab_define(const char *key, int addr)
{
	Label_T p;
	int h = hash(key);
	
	p = buckets[h];
	while(p != NULL) {
		if(!strcmp(key, p->key)) break;
		p = p->link;
	}
	if(p != NULL) {
		if(p->datoff == NULL) QUIT("重复定义标号");
		while(!Stack_empty(p->datoff)) {
			mem[Stack_pop(p->datoff)] = addr;
		}
		free(p->datoff);
		p->datoff = NULL;
		p->addr = addr;
	}else {
		int len = strlen(key);
		p = malloc(sizeof(*p) + len + 1);
		if(p == NULL) QUIT("lab_define内存分配失败");
		p->key = (char *)(p + 1);
		strcpy(p->key, key);
		p->datoff = NULL;
		p->addr = addr;
		p->link = buckets[h];
		buckets[h] = p;
	}
}

int
lab_get(const char *key, int off)
{
	Label_T p;
	int addr = 0;
	int h = hash(key);
	
	p = buckets[h];
	while(p != NULL) {
		if(!strcmp(key, p->key)) break;
		p = p->link;
	}
	if(p != NULL) {
		if(p->datoff == NULL) addr = p->addr;
		else p->datoff = Stack_push(p->datoff, off);
	}else {
		int len = strlen(key);
		p = malloc(sizeof(*p) + len + 1);
		if(p == NULL) QUIT("lab_get内存分配失败");
		p->key = (char *)(p + 1);
		strcpy(p->key, key);
		p->datoff = Stack_new();
		Stack_push(p->datoff, off);
		p->addr = line;
		p->link = buckets[h];
		buckets[h] = p;
	}
	return addr;
}

void
lab_map(void map(Label_T p, void *cl), void *cl)
{
	int i;
	for(i = 0; i < NELEMS(buckets); ++i) {
		Label_T lab = buckets[i];
		while(lab != NULL) {
			map(lab, cl);
			lab = lab->link;
		}
	} /* end for */
}

void
lab_free(void)
{
	int i;
	for(i = 0; i < NELEMS(buckets); ++i) {
		Label_T lab = buckets[i];
		Label_T p = lab;
		while(lab != NULL) {
			p = lab->link;
			free(lab);
			lab = p;
		}
	} /* end for */
}

#endif
