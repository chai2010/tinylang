// +build ignore

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/chai2010/tinylang/pkg/spec/comet"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	bin, pc := loadBin("sum.comet")
	fmt.Println(pc, len(bin), bin)
	vm := comet.NewComent(nil, bin, pc)
	vm.DebugRun()
}

func loadBin(path string) (bin []uint16, pc int) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var hdr struct {
		PC  uint16
		Len uint16
	}

	if err = binary.Read(f, binary.LittleEndian, &hdr); err != nil {
		log.Fatal(err)
	}

	data := make([]uint16, int(hdr.Len))
	if err = binary.Read(f, binary.LittleEndian, &data); err != nil {
		log.Fatal(err)
	}

	return data, int(hdr.PC)
}
