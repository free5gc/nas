//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/free5gc/nas/internal/tools/generator"
)

// Generate table of types in nasType package
func main() {
	dirs, err := os.ReadDir("nasType")
	if err != nil {
		panic(err)
	}

	fOut := generator.NewOutputFile("internal/tools/generator/types.go", "generator", []string{
		"\"reflect\"",
		"",
		"\"github.com/free5gc/nas/nasType\"",
	})

	fmt.Fprintln(fOut, "var nasTypeTable map[string]reflect.Type = map[string]reflect.Type{")
	for _, dir := range dirs {
		name := dir.Name()
		// Assume one type by one file
		if strings.HasPrefix(name, "NAS_") && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			name := strings.TrimPrefix(name, "NAS_")
			name = strings.TrimSuffix(name, ".go")
			fmt.Fprintf(fOut, "\"%s\": reflect.TypeOf(nasType.%s{}),\n", name, name)
		}
	}
	fmt.Fprintln(fOut, "}")

	if err := fOut.Close(); err != nil {
		panic(err)
	}
}
