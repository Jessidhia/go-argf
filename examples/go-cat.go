package main

import (
	"fmt"
	"github.com/Kovensky/go-argf"
	"io"
	"os"
)

func main() {
	argf := argf.New(os.Args[1:]...)
	err := io.ErrNoProgress // can't be nil, will be overwritten

	// io.Copy returns nil on EOF
	for err != nil {
		_, err = io.Copy(os.Stdout, argf)
		if err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0]+":", err)
		}
	}
}
