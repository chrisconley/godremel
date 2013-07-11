package go_dremel

import (
    "fmt"
    "testing"
    "io/ioutil"
    "encoding/json"
    )

func TestStripeRecord(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }

  file, e = ioutil.ReadFile("./record1.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var record interface{}
  err = json.Unmarshal(file, &record)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }

  memstore := MemStore{map[string][]Row{}}
  StripeRecord(schema, record, &memstore, RootWriter, 0)
  for c, rows := range memstore.Data {
    fmt.Printf("%v\n", c)
    for _, r := range rows {
      fmt.Printf("%v\n", r)
    }
  }
  if 1 != 2 {
    t.Errorf("hi")
  }
}
