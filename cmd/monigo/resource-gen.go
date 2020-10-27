//+build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const blob = "monigo/resources/blob.go"

var packageTemplate = template.Must(template.New("").Funcs(map[string]interface{}{"conv": FormatByteSlice}).Parse(`// This file is was generated using 'go run cmd/monigo/resource-gen.go'
// DO NOT EDIT

package resources

func init(){
 {{- range $name, $file := . }}
     Add("{{ $name }}", []byte{ {{ conv $file }} })
 {{- end }}
}
`))

func FormatByteSlice(sl []byte) string {
	builder := strings.Builder{}
	for _, v := range sl {
		builder.WriteString(fmt.Sprintf("%d,", int(v)))
	}
	return builder.String()
}

func main() {
	log.Println("Baking box... \U0001F4E6")

	if _, err := os.Stat("monigo/resources/data"); os.IsNotExist(err) {
		log.Fatal("Resources directory does not exists")
	}

	resources := make(map[string][]byte)
	err := filepath.Walk("monigo/resources/data", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error :", err)
			return err
		}
		relativePath := filepath.ToSlash(strings.TrimPrefix(path, "monigo/resources/data"))
		if info.IsDir() {
			log.Println(path, "is a directory, skipping... \U0001F47B")
			return nil
		} else {
			log.Println(path, "is a file, baking in... \U0001F31F")
			b, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading %s: %s", path, err)
				return err
			}
			resources[relativePath] = b
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error walking through box directory:",err)
	}

	f, err := os.Create(blob)
	if err != nil {
		log.Fatal("Error creating blob file:", err)
	}
	defer f.Close()

	builder := &bytes.Buffer{}

	err = packageTemplate.Execute(builder, resources)
	if err != nil {
		log.Fatal("Error executing template", err)
	}

	data, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatal("Error formatting generated code", err)
	}
	err= ioutil.WriteFile(blob, data, os.ModePerm)
	if err != nil {
		log.Fatal("Error writing blob file", err)
	}

	log.Println("Baking box done... \U0001F680")
}