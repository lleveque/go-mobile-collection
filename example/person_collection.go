// generated by gomobile-collection -- DO NOT EDIT
// to re-generate, run the following command :
// gomobile-collection --use-errors=false example/person.go
package society
import (
  "fmt"
  "encoding/json"
)

type PersonCollection struct {
  s []*Person
}

func NewPersonCollection() *PersonCollection {
  return &PersonCollection{}
}

func (v *PersonCollection) Clear() {
  v.s = v.s[:0]
}

func (v *PersonCollection) Equal(rhs *PersonCollection) bool {
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

func (v *PersonCollection) MarshalJSON() (res []byte) {
  res, err := json.Marshal(v.s)
  if err != nil {
    fmt.Printf("PersonCollection : error marshalling JSON")
    return nil
  }
  return res
}

func (v *PersonCollection) UnmarshalJSON(data []byte) {
  err := json.Unmarshal(data, &v.s)
  if err != nil {
    fmt.Printf("PersonCollection : error unmarshalling JSON")
  }
}

func (v *PersonCollection) Copy(rhs *PersonCollection) {
  v.s = make([]*Person, len(rhs.s))
  copy(v.s, rhs.s)
}

func (v *PersonCollection) Clone() *PersonCollection {
  return &PersonCollection{
    s: v.s[:],
  }
}

func (v *PersonCollection) Index(rhs *Person) int {
  for i, lhs := range v.s {
    if lhs == rhs {
      return i
    }
  }
  return -1
}

func (v *PersonCollection) Insert(i int, n *Person) {
  if i < 0 || i > len(v.s) {
    fmt.Printf("PersonCollection : error trying to insert at index %d\n", i)
    return
  }
  v.s = append(v.s, nil)
  copy(v.s[i+1:], v.s[i:])
  v.s[i] = n
}

func (v *PersonCollection) Remove(i int) {
  if i < 0 || i >= len(v.s) {
    fmt.Printf("PersonCollection : error trying to remove bad index %d\n", i)
    return
  }
  copy(v.s[i:], v.s[i+1:])
  v.s[len(v.s)-1] = nil
  v.s = v.s[:len(v.s)-1]
}

func (v *PersonCollection) Count() int {
  return len(v.s)
}

func (v *PersonCollection) At(i int) *Person {
  if i < 0 || i >= len(v.s) {
    fmt.Printf("PersonCollection : invalid index %d\n", i)
    return
  }
  return v.s[i]
}
