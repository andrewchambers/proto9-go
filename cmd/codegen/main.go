package main

import (
	"fmt"
	"go/types"
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

func generateEncoder(topLevelName string, t *types.Struct) {
	fmt.Printf("func (m *%s) Encode(b *bytes.Buffer) {\n", topLevelName)
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		if isNumberType(f.Type()) {
			fmt.Printf("encode%s(b, m.%s)\n", upperCaseFirst(f.Type().String()), f.Name())
		} else if isByteSlice(f.Type()) {
			fmt.Printf("encodeByteSlice(b, m.%s)\n", f.Name())
		} else if isStringType(f.Type()) {
			fmt.Printf("encodeString(b, m.%s)\n", f.Name())
		} else if isQidSlice(f.Type()) {
			fmt.Printf("encodeQids(b, m.%s)\n", f.Name())
		} else if isContainerType(f.Type()) {
			fmt.Printf("m.%s.Encode(b)\n", f.Name())
		} else {
			fatalErr(
				fmt.Errorf("don't know how to encode type field %s", f.Type()),
			)
		}
	}
	fmt.Printf("}\n\n")
}

func generateDecoder(topLevelName string, t *types.Struct) {
	fmt.Printf("func (m *%s) Decode(b *bytes.Buffer) {\n", topLevelName)
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		if isNumberType(f.Type()) {
			fmt.Printf("m.%s = decode%s(b)\n", f.Name(), upperCaseFirst(f.Type().String()))
		} else if isByteSlice(f.Type()) {
			fmt.Printf("m.%s = decodeByteSlice(b)\n", f.Name())
		} else if isStringType(f.Type()) {
			fmt.Printf("m.%s = decodeString(b)\n", f.Name())
		} else if isQidSlice(f.Type()) {
			fmt.Printf("m.%s = decodeQids(b)\n", f.Name())
		} else if isContainerType(f.Type()) {
			fmt.Printf("m.%s.Decode(b)\n", f.Name())
		} else {
			fatalErr(
				fmt.Errorf("don't know how to encode type field %s", f.Type()),
			)
		}
	}
	fmt.Printf("}\n\n")
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func main() {
	// Handle arguments to command
	if len(os.Args) != 2 {
		fatalErr(fmt.Errorf("expected exactly one argument: <source type>"))
	}

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, os.Args[1])
	if err != nil {
		fatalErr(fmt.Errorf("loading packages for inspection: %v", err))
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	fmt.Println(`
		package proto9

		import (
			"bytes"
		)		
	`)

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
			generateEncoder(name, t)
			generateDecoder(name, t)
		}
	}
}
