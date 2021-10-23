# Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

default:
	go run gen1024.go 2> 1024.tiny
	go run .. 1024.tiny > a.out.ll
	clang -Wno-override-module -o a.out.exe a.out.ll
	./a.out.exe

clean:
