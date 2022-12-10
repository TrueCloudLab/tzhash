package main

import (
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/TrueCloudLab/tzhash/tz"
	"golang.org/x/sys/cpu"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	filename   = flag.String("name", "-", "file to use")
	hashimpl   = flag.String("impl", "", "implementation to use")
)

func main() {
	var (
		f   io.Reader
		err error
	)

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *filename != "-" {
		if f, err = os.Open(*filename); err != nil {
			log.Fatal("could not open file: ", err)
		}
	} else {
		f = os.Stdin
	}

	// Override CPU feature flags to make sure a proper backend is used.
	var h hash.Hash
	switch *hashimpl {
	case "avx":
		cpu.X86.HasAVX = true
		cpu.X86.HasAVX2 = false
		h = tz.New()
	case "avx2":
		cpu.X86.HasAVX = true
		cpu.X86.HasAVX2 = true
		h = tz.New()
	case "generic":
		cpu.X86.HasAVX = false
		cpu.X86.HasAVX2 = false
		h = tz.New()
	default:
		log.Fatalf("Invalid backend: %s", *hashimpl)
	}

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal("error while reading file: ", err)
	}
	fmt.Printf("%x\t%s\n", h.Sum(nil), *filename)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
