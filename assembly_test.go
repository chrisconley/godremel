package go_dremel

import (
    //"fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func TestAssembly(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Schema
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
}
