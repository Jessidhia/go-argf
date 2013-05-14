package argf

import (
	"io"
	"os"
)

type argf struct {
	fh   *os.File
	args []string
}

// Returns an io.ReadCloser that behaves like ruby's ARGF, but without modifying os.Args.
// Use by giving (os.Args[1:]...) as argument, or (flag.Args()...).
//
// The io.ReadCloser behaves as a concatenation of all named files, opening each file in
// sequence as needed to feed the Read calls. A filename of "-" means reading from os.Stdin.
// Calling New with no arguments also reads from os.Stdin.
//
// Read() only returns EOF once all files are read. Any error other than EOF may only
// apply to the specific failed file; calling Read() again will skip the bad file and open
// the next one.
//
// Calling Close() will close the currently open underlying file, or do nothing if it's
// already closed.
func New(args ...string) io.ReadCloser {
	if len(args) == 0 {
		args = []string{"-"}
	}

	return &argf{args: args}
}

func (argf *argf) Read(p []byte) (n int, err error) {
	if argf.fh == nil {
		if len(argf.args) == 0 {
			return 0, io.EOF
		}

		// "pop" path from argf.args
		path := argf.args[0]
		argf.args = argf.args[1:]

		if path == "-" {
			argf.fh = os.Stdin
		} else {
			argf.fh, err = os.Open(path)
			if err != nil {
				return 0, err
			}
		}
	}
	n, err = argf.fh.Read(p)

	switch err {
	case io.EOF:
		// let the next Read call handle the next file, or the actual EOF
		err = nil
		argf.fh.Close()
		argf.fh = nil

		if n == 0 {
			// recurse once, to avoid returning (0, nil)
			return argf.Read(p)
		}
	case nil: // AOK
	default:
		// skip this file
		argf.fh.Close()
		argf.fh = nil
	}

	return
}

func (argf *argf) Close() (err error) {
	if argf.fh != nil {
		err = argf.fh.Close()
		argf.fh = nil
	}
	return
}
