package go_dremel

import (
    "fmt"
    "testing"
    "io/ioutil"
    "encoding/json"
    )

func readJson(t *testing.T, filePath string, o interface{}) {
  file, e := ioutil.ReadFile(filePath)
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  err := json.Unmarshal(file, o)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
}

func TestStripeRecord(t *testing.T) {
  var schema Field
  readJson(t, "./docs.json", &schema)

  memstore := MemStore{map[string][]Row{}}

  var record interface{}
  readJson(t, "./record1.json", &record)

  var record2 interface{}
  readJson(t, "./record2.json", &record2)

  StripeRecord(schema, record, &memstore, RootWriter, 0)
  StripeRecord(schema, record2, &memstore, RootWriter, 0)

  for c, rows := range memstore.Data {
    fmt.Printf("%v\n", c)
    for _, r := range rows {
      fmt.Printf("%v\n", r)
    }
  }
  countryRows := []Row{
    Row{"us", 0, 3},
    Row{nil, 2, 2},
    Row{nil, 1, 1},
    Row{"gb", 1, 3},
    Row{nil, 0, 1},
  }
  for i := 0; i < len(memstore.Data["names.languages.country"]); i++ {
    if memstore.Data["names.languages.country"][i] != countryRows[i] {
      t.Errorf("hi")
    }
  }
}
