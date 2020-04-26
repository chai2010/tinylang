// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// 加载指令, 返回结果可用于 NewComet 执行
func LoadProgram(path string, src interface{}) (prog *Program, err error) {
	if src == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		src = data
	}

	var r io.Reader
	switch src := src.(type) {
	case io.Reader:
		r = src
	case []byte:
		r = bytes.NewReader(src)
	case string:
		r = strings.NewReader(src)
	default:
		return nil, fmt.Errorf("unnown src type: %T", src)
	}

	var hdr struct {
		PC  uint16
		Len uint16
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}

	data := make([]uint16, int(hdr.Len))
	if err = binary.Read(r, binary.LittleEndian, &data); err != nil {
		return nil, err
	}

	prog = &Program{
		PC:  hdr.PC,
		Bin: data,
	}
	return
}
