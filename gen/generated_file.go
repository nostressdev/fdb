package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

type GeneratedFile struct {
	buf bytes.Buffer
}

func (f *GeneratedFile) Import(imp string) {
	panic("Implemented me")
}

func (f *GeneratedFile) Printf(str string, args ...string) {
	_, err := fmt.Fprintf(&f.buf, str, args)
	if err != nil {
		panic(err)
	}
}

func (f *GeneratedFile) Println(str string) {
	_, err := fmt.Fprintln(&f.buf, str)
	if err != nil {
		panic(err)
	}
}

func (f *GeneratedFile) Print(str string) {
	_, err := fmt.Fprint(&f.buf, str)
	if err != nil {
		panic(err)
	}
}

func (f *GeneratedFile) Write(out io.Writer) {
	// TODO Imports
	bufBytes, err := format.Source(f.buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = out.Write(bufBytes)
	if err != nil {
		panic(err)
	}
}