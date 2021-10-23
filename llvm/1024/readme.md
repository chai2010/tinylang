# LLVM 版本 2014 彩蛋

Tiny 的 `write` 函数只能输出整数, 通过以下或者增加输出字符能力:

```c
void __tiny_write(int x) {
	if(x > 1024*1024) {
		printf("%c", x-1024*1024);
		return;
	}
	printf("%d\n", x);
}
```

然后通过以下 Go 程序产生 Tiny 代码:

```go
// go run a.go 2> 1024.tiny

package main

func main() {
	for _, c := range []byte(s[1:]) {
		println(`write 1024 * 1024 +`, c, ";")
	}
}

const s = `...`
```

然后将生成的 Tiny 程序再编译为本地程序:

```
$ go run .. 1024.tiny > a.out.ll
$ clang -Wno-override-module -o a.out.exe a.out.ll
$ ./a.out.exe
+---+    +---+
| o |    | o |
|   +----+   |
|            |
|    1024    |
|            |
+------------+
```
