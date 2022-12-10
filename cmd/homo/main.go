package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/TrueCloudLab/tzhash/tz"
)

var (
	concat   = flag.Bool("concat", false, "Concatenate hashes")
	filename = flag.String("file", "", "File to read from")
)

func main() {
	var (
		err   error
		file  = os.Stdin
		lines = make([]string, 0, 10)
	)

	flag.Parse()
	if *filename != "" {
		if file, err = os.Open(*filename); err != nil {
			fatal("error while opening file: %v", err)
		}
	}

	for f := bufio.NewScanner(file); f.Scan(); {
		lines = append(lines, f.Text())
	}

	if *concat {
		var (
			h      []byte
			hashes = make([][]byte, len(lines))
		)
		for i := range lines {
			if hashes[i], err = hex.DecodeString(lines[i]); err != nil {
				fatal("error while decoding hex-string: %v", err)
			}
		}
		h, err := tz.Concat(hashes)
		if err != nil {
			fatal("error while concatenating hashes: %v", err)
		}
		fmt.Println(hex.EncodeToString(h))
		return
	}

	for i := range lines {
		h := tz.Sum([]byte(lines[i]))
		fmt.Println(hex.EncodeToString(h[:]))
	}
}

func fatal(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}
