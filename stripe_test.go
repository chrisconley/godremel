package go_dremel

import (
    //"fmt"
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

  StripeRecord(schema, record, &MemStore{})
  if 1 != 2 {
    t.Errorf("hi")
  }
}
