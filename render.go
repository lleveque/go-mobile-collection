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
	generatedTemplate = template.Must(template.New("render").Parse(`// generated by gomobile-collection -- DO NOT EDIT{{$errorsOn := .ErrorsOn}}
// to re-generate, run the following command :
// gomobile-collection --use-errors={{$errorsOn}} {{.Path}}
package {{.Package}}
import (
{{if $errorsOn}}  "errors"
{{end}}  "fmt"
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

{{if $errorsOn}}func (v *{{.Name}}Collection) MarshalJSON() ([]byte, error) {
  return json.Marshal(v.s)
{{else}}func (v *{{.Name}}Collection) MarshalJSON() (res []byte) {
  res, err := json.Marshal(v.s)
  if err != nil {
    fmt.Printf("{{.Name}}Collection : error marshalling JSON")
    return nil
  }
{{end}}}

{{if $errorsOn}}func (v *{{.Name}}Collection) UnmarshalJSON(data []byte) error {
  return json.Unmarshal(data, &v.s)
{{else}}func (v *{{.Name}}Collection) UnmarshalJSON(data []byte) {
  err := json.Unmarshal(data, &v.s)
  if err != nil {
    fmt.Printf("{{.Name}}Collection : error unmarshalling JSON")
  }
{{end}}}

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

{{if $errorsOn}}func (v *{{.Name}}Collection) Insert(i int, n *{{.Name}}) error{{else}}func (v *{{.Name}}Collection) Insert(i int, n *{{.Name}}){{end}} {
  if i < 0 || i > len(v.s) {
{{if $errorsOn}}    return errors.New(fmt.Sprintf("{{.Name}}Collection : error trying to insert at index %d\n", i)){{else}}    fmt.Printf("{{.Name}}Collection : error trying to insert at index %d\n", i)
    return{{end}}
  }
  v.s = append(v.s, nil)
  copy(v.s[i+1:], v.s[i:])
  v.s[i] = n
{{if $errorsOn}}  return nil
{{end}}}

{{if $errorsOn}}func (v *{{.Name}}Collection) Remove(i int) error{{else}}func (v *{{.Name}}Collection) Remove(i int){{end}} {
  if i < 0 || i >= len(v.s) {
{{if $errorsOn}}    return errors.New(fmt.Sprintf("{{.Name}}Collection : error trying to remove bad index %d\n", i)){{else}}    fmt.Printf("{{.Name}}Collection : error trying to remove bad index %d\n", i)
    return{{end}}
  }
  copy(v.s[i:], v.s[i+1:])
  v.s[len(v.s)-1] = nil
  v.s = v.s[:len(v.s)-1]
{{if $errorsOn}}  return nil
{{end}}}

func (v *{{.Name}}Collection) Count() int {
  return len(v.s)
}

{{if $errorsOn}}func (v *{{.Name}}Collection) At(i int) (*{{.Name}}, error){{else}}func (v *{{.Name}}Collection) At(i int) *{{.Name}}{{end}} {
  if i < 0 || i >= len(v.s) {
{{if $errorsOn}}    return nil, errors.New(fmt.Sprintf("{{.Name}}Collection : invalid index %d\n", i)){{else}}    fmt.Printf("{{.Name}}Collection : invalid index %d\n", i){{end}}
  }
{{if $errorsOn}}  return v.s[i], nil{{else}}  return v.s[i]{{end}}
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
  Path string
	Package string
	Types   []GeneratedType
  ErrorsOn bool
}

func render(w io.Writer, inputPath string, packageName string, types []GeneratedType, errorsOn bool) error {
	return generatedTemplate.Execute(w, generateTemplateData{inputPath, packageName, types, errorsOn})
}
