// +build ignore

package main

import (
	"encoding/binary"
	"flag"
	"log"
	"os"

	"github.com/chai2010/tinylang/pkg/spec/comet"
)

var (
	flagFile  = flag.String("f", "sum.comet", "comet app file")
	flagDebug = flag.Bool("d", false, "debug mode")
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	flag.Parse()

	bin, pc := loadBin(*flagFile)
	vm := comet.NewComent(nil, bin, pc)

	if *flagDebug {
		vm.DebugRun()
	} else {
		vm.Run()
	}
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
