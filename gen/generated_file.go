package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"sort"
	"strconv"
)

type GeneratedFile struct {
	buf     bytes.Buffer
	imports [][2]string
}

func (f *GeneratedFile) Import(impPath string, impName ...string) {
	if len(impName) == 0 {
		f.imports = append(f.imports, [2]string{"", impPath})
	} else {
		f.imports = append(f.imports, [2]string{impName[0], impPath})
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
	// Код практически полностью скопирован
	// из "google.golang.org/protobuf/compiler/protogen#GeneratedFile.Content()"
	// однако только НУЖНАЯ часть, то есть метод Write() НЕ является копией Content()

	original := f.buf.Bytes()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		panic(fmt.Errorf("unparsable Go source: %v\n%v", err, src.String()))
	}
	sort.Slice(f.imports, func(i, j int) bool {
		return f.imports[i][1] < f.imports[j][1]
	})
	if len(f.imports) > 0 {
		// Insert block after package statement or
		// possible comment attached to the end of the package statement.
		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}

		// Construct the import block.
		impDecl := &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: pos,
			Lparen: pos,
			Rparen: pos,
		}
		for _, importPath := range f.imports {
			impDecl.Specs = append(impDecl.Specs, &ast.ImportSpec{
				Name: &ast.Ident{
					Name:    importPath[0],
					NamePos: pos,
				},
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    strconv.Quote(importPath[1]),
					ValuePos: pos,
				},
				EndPos: pos,
			})
		}
		file.Decls = append([]ast.Decl{impDecl}, file.Decls...)
	}

	var res bytes.Buffer
	if err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(&res, fset, file); err != nil {
		panic(fmt.Errorf("can not reformat Go source: %v", err))
	}

	bufBytes, err := format.Source(res.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = out.Write(bufBytes)
	if err != nil {
		panic(err)
	}
}
