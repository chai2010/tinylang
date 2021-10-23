# github.com/sbinet/wasm/tree/master/cmd/wasm-dump

default:
	clang tiny_lib.c -emit-llvm -S -c -o tiny_lib.ll
	clang -Wno-override-module -o a.out.exe hello.c tiny_lib.ll
	./a.out.exe

clean:
