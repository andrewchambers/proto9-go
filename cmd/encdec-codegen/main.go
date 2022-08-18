package main

import (
	"fmt"
	"go/types"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func isByteSlice(t types.Type) bool {
	switch t := t.(type) {
	case *types.Slice:
		switch etype := t.Elem().(type) {
		case *types.Basic:
			switch etype.Kind() {
			case types.Byte:
				return true
			}
		}
	}
	return false
}

func isStringSlice(t types.Type) bool {
	switch t := t.(type) {
	case *types.Slice:
		switch etype := t.Elem().(type) {
		case *types.Basic:
			switch etype.Kind() {
			case types.String:
				return true
			}
		}
	}
	return false
}

func isQidSlice(t types.Type) bool {
	switch t := t.(type) {
	case *types.Slice:
		switch etype := t.Elem().(type) {
		case *types.Named:
			return etype.Obj().Name() == "Qid"
		}
	}
	return false
}

func isDirEntSlice(t types.Type) bool {
	switch t := t.(type) {
	case *types.Slice:
		switch etype := t.Elem().(type) {
		case *types.Named:
			return etype.Obj().Name() == "DirEnt"
		}
	}
	return false
}

func isNumberType(t types.Type) bool {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return true
		}
	}
	return false
}

func isStringType(t types.Type) bool {
	switch t := t.(type) {
	case *types.Basic:
		return t.Kind() == types.String
	}
	return false
}

func isContainerType(t types.Type) bool {
	switch t.(type) {
	case *types.Named:
		return true
	}
	return false
}

func upperCaseFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func outPrelude(out io.Writer) {
	fmt.Fprintln(out, "package proto9\nimport (\n\t\"bytes\"\n)")
}

func outEncodedSize(topLevelName string, t *types.Struct, out io.Writer) {
	fmt.Fprintf(out, "func (v *%s) EncodedSize() uint64 {\n", topLevelName)
	fmt.Fprintln(out, "sz := uint64(0)")
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		if isNumberType(f.Type()) {
			switch f.Type().(*types.Basic).Kind() {
			case types.Uint8:
				fmt.Fprintf(out, "sz += 1 // %s\n", f.Name())
			case types.Uint16:
				fmt.Fprintf(out, "sz += 2 // %s\n", f.Name())
			case types.Uint32:
				fmt.Fprintf(out, "sz += 4 // %s\n", f.Name())
			case types.Uint64:
				fmt.Fprintf(out, "sz += 8 // %s\n", f.Name())
			default:
				fatalErr(
					fmt.Errorf("don't know how to get size of type field %s", f.Type()),
				)
			}
		} else if isStringSlice(f.Type()) {
			fmt.Fprintf(out, "// %s\n", f.Name())
			fmt.Fprintf(out, "sz += 2\n")
			fmt.Fprintf(out, "for _, s := range v.%s {\n", f.Name())
			fmt.Fprintf(out, "sz += 2 + uint64(len(s))\n")
			fmt.Fprintf(out, "}\n")
		} else if isDirEntSlice(f.Type()) {
			fmt.Fprintf(out, "// %s\n", f.Name())
			fmt.Fprintf(out, "sz += 4\n")
			fmt.Fprintf(out, "for i := range v.%s {\n", f.Name())
			fmt.Fprintf(out, "sz += v.%s[i].EncodedSize()\n", f.Name())
			fmt.Fprintf(out, "}\n")
		} else if isByteSlice(f.Type()) {
			fmt.Fprintf(out, "sz += 4 + uint64(len(v.%s))\n", f.Name())
		} else if isStringType(f.Type()) {
			fmt.Fprintf(out, "sz += 2 + uint64(len(v.%s))\n", f.Name())
		} else if isQidSlice(f.Type()) {
			fmt.Fprintf(out, "sz += 2 + uint64(len(v.%s))*13\n", f.Name())
		} else if isContainerType(f.Type()) {
			fmt.Fprintf(out, "sz += v.%s.EncodedSize()\n", f.Name())
		} else {
			fatalErr(
				fmt.Errorf("don't know how to get size of type field %s", f.Type()),
			)
		}
	}
	fmt.Fprintf(out, "return sz\n}\n\n")
}

func outEncoder(topLevelName string, t *types.Struct, out io.Writer) {
	fmt.Fprintf(out, "func (v *%s) Encode(b *bytes.Buffer) error {\n", topLevelName)
	fmt.Fprintln(out, "var err error")
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		if isNumberType(f.Type()) {
			fmt.Fprintf(out, "err = encode%s(b, v.%s)\n", upperCaseFirst(f.Type().String()), f.Name())
		} else if isByteSlice(f.Type()) {
			fmt.Fprintf(out, "err = encodeByteSlice(b, v.%s)\n", f.Name())
		} else if isStringSlice(f.Type()) {
			fmt.Fprintf(out, "err = encodeStringSlice(b, v.%s)\n", f.Name())
		} else if isDirEntSlice(f.Type()) {
			fmt.Fprintf(out, "err = encodeDirEntSlice(b, v.%s)\n", f.Name())
		} else if isStringType(f.Type()) {
			fmt.Fprintf(out, "err = encodeString(b, v.%s)\n", f.Name())
		} else if isQidSlice(f.Type()) {
			fmt.Fprintf(out, "err = encodeQids(b, v.%s)\n", f.Name())
		} else if isContainerType(f.Type()) {
			fmt.Fprintf(out, "err = v.%s.Encode(b)\n", f.Name())
		} else {
			fatalErr(
				fmt.Errorf("don't know how to encode type field %s", f.Type()),
			)
		}
		fmt.Fprintln(out, "if err != nil {\nreturn err\n}")
	}
	fmt.Fprintf(out, "return nil\n}\n\n")
}

func outDecoder(topLevelName string, t *types.Struct, out io.Writer) {
	fmt.Fprintf(out, "func (v *%s) Decode(b *bytes.Buffer) error {\n", topLevelName)
	fmt.Fprintln(out, "var err error")
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		if isNumberType(f.Type()) {
			fmt.Fprintf(out, "v.%s, err = decode%s(b)\n", f.Name(), upperCaseFirst(f.Type().String()))
		} else if isByteSlice(f.Type()) {
			fmt.Fprintf(out, "v.%s, err = decodeByteSlice(b)\n", f.Name())
		} else if isDirEntSlice(f.Type()) {
			fmt.Fprintf(out, "v.%s, err = decodeDirEntSlice(b)\n", f.Name())
		} else if isStringSlice(f.Type()) {
			fmt.Fprintf(out, "v.%s, err = decodeStringSlice(b)\n", f.Name())
		} else if isStringType(f.Type()) {
			fmt.Fprintf(out, "v.%s, err = decodeString(b)\n", f.Name())
		} else if isQidSlice(f.Type()) {
			fmt.Fprintf(out, "v.%s, err = decodeQids(b)\n", f.Name())
		} else if isContainerType(f.Type()) {
			fmt.Fprintf(out, "err = v.%s.Decode(b)\n", f.Name())
		} else {
			fatalErr(
				fmt.Errorf("don't know how to encode type field %s", f.Type()),
			)
		}
		fmt.Fprintln(out, "if err != nil {\nreturn err\n}")
	}
	fmt.Fprintf(out, "return nil\n}\n\n")
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func main() {
	// Handle arguments to command
	if len(os.Args) != 3 {
		fatalErr(fmt.Errorf("expected: PACKAGE OUTFILE"))
	}

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, os.Args[1])
	if err != nil {
		fatalErr(fmt.Errorf("loading packages for inspection: %v", err))
	}

	if packages.PrintErrors(pkgs) > 0 {
		fmt.Fprintf(os.Stderr, "continuing despite errors...\n")
	}

	out := os.Stdout

	if os.Args[2] != "-" {
		out, err = os.Create(os.Args[2])
		if err != nil {
			fatalErr(err)
		}
	}

	outPrelude(out)
	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if _, ok := obj.(*types.TypeName); !ok {
				continue
			}

			t, ok := obj.Type().Underlying().(*types.Struct)
			if !ok {
				continue
			}

			fileName := pkg.Fset.PositionFor(obj.Pos(), false).Filename
			if !strings.HasSuffix(fileName, "/fcalltypes.go") {
				continue
			}

			outEncodedSize(name, t, out)
			outEncoder(name, t, out)
			outDecoder(name, t, out)
		}
	}

	err = out.Close()
	if err != nil {
		fatalErr(err)
	}
}
