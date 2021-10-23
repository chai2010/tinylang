# Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

default:
	-@go fmt ./...
	-@go vet ./...

	go run main.go ./hello.tiny > a.out.ll
	clang -Wno-override-module -o a.out.exe a.out.ll
	./a.out.exe

clean:
