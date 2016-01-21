// Copyright 2014 Brett Slatkin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	generatedTemplate = template.Must(template.New("render").Parse(`// generated by collection-wrapper -- DO NOT EDIT
package {{.Package}}

import (
  "fmt"
  "encoding/json"
)

{{range .Types}}
type {{.Name}}Collection struct {
  s []*{{.Name}}
}

func New{{.Name}}Collection() *{{.Name}}Collection {
  return &{{.Name}}Collection{}
}

func (v *{{.Name}}Collection) Clear() {
  v.s = v.s[:0]
}

func (v *{{.Name}}Collection) Equal(rhs *{{.Name}}Collection) bool {
  if rhs == nil {
    return false
  }

  if len(v.s) != len(rhs.s) {
    return false
  }

  for i := range v.s {
    if !v.s[i].Equal(rhs.s[i]) {
      return false
    }
  }

  return true
}

func (v *{{.Name}}Collection) MarshalJSON() ([]byte, error) {
  return json.Marshal(v.s)
}

func (v *{{.Name}}Collection) UnmarshalJSON(data []byte) error {
  return json.Unmarshal(data, &v.s)
}

func (v *{{.Name}}Collection) Copy(rhs *{{.Name}}Collection) {
  v.s = make([]*{{.Name}}, len(rhs.s))
  copy(v.s, rhs.s)
}

func (v *{{.Name}}Collection) Clone() *{{.Name}}Collection {
  return &{{.Name}}Collection{
    s: v.s[:],
  }
}

func (v *{{.Name}}Collection) Index(rhs *{{.Name}}) int {
  for i, lhs := range v.s {
    if lhs == rhs {
      return i
    }
  }
  return -1
}

func (v *{{.Name}}Collection) Insert(i int, n *{{.Name}}) {
  if i < 0 || i > len(v.s) {
    fmt.Printf("Vapi::{{.Name}}Collection field_values.go error trying to insert at index %d\n", i)
    return
  }
  v.s = append(v.s, nil)
  copy(v.s[i+1:], v.s[i:])
  v.s[i] = n
}

func (v *{{.Name}}Collection) Remove(i int) {
  if i < 0 || i >= len(v.s) {
    fmt.Printf("Vapi::{{.Name}}Collection field_values.go error trying to remove bad index %d\n", i)
    return
  }
  copy(v.s[i:], v.s[i+1:])
  v.s[len(v.s)-1] = nil
  v.s = v.s[:len(v.s)-1]
}

func (v *{{.Name}}Collection) Count() int {
  return len(v.s)
}

func (v *{{.Name}}Collection) At(i int) *{{.Name}} {
  if i < 0 || i >= len(v.s) {
    fmt.Printf("Vapi::{{.Name}}Collection field_values.go invalid index %d\n", i)
  }
  return v.s[i]
}

{{end}}`))
)

type GeneratedType struct {
	Name string
}

func getRenderedPath(inputPath string) (string, error) {
	if !strings.HasSuffix(inputPath, ".go") {
		return "", fmt.Errorf("Input path %s doesn't have .go extension", inputPath)
	}
	trimmed := strings.TrimSuffix(inputPath, ".go")
	dir, file := filepath.Split(trimmed)
	return filepath.Join(dir, fmt.Sprintf("%s_collection.go", file)), nil
}

type generateTemplateData struct {
	Package string
	Types   []GeneratedType
}

func render(w io.Writer, packageName string, types []GeneratedType) error {
	return generatedTemplate.Execute(w, generateTemplateData{packageName, types})
}