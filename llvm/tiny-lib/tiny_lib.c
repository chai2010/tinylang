#include "tiny_lib.h"

#include <stdio.h>

int __tiny_read() {
	int x;
	printf("READ: ");
	scanf("%d", &x);
	return x;
}

void __tiny_write(int x) {
	printf("%d\n", x);
}
