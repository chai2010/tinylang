# Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

default:
	go run 0_gen.go > zz_utils.go
	make goyacc

	-@go fmt
	go test

goyacc:
	goyacc -o tiny.y.go -p "yy" tiny.y

clean:
