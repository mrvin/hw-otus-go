package main

import (
	"flag"
	"log"
)

var (
	from, to      string
	limit, offset int64
	isQuiet       bool
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.BoolVar(&isQuiet, "quiet", false, "run without output")
}

func main() {
	flag.Parse()

	if from == "" {
		log.Fatalf("go-cp: file to read from is empty")
	}
	if to == "" {
		log.Fatalf("go-cp: file to write to is empty")
	}

	if err := Copy(from, to, offset, limit, isQuiet); err != nil {
		log.Fatalf("go-cp: %v", err)
	}
}
