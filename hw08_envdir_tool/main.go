package main

import (
	"log"
	"os"
)

func main() {
	returnCode := 111

	if len(os.Args) < 3 {
		log.Printf("go-envdir: not enough arguments < 2")
		os.Exit(returnCode)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Printf("go-envdir: %v", err)
		os.Exit(returnCode)
	}

	returnCode, err = RunCmd(os.Args[2:len(os.Args)], env)
	if err != nil {
		log.Printf("go-envdir: %v", err)
		os.Exit(returnCode)
	}

	os.Exit(returnCode)
}
