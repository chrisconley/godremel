package go_dremel

import (
    "fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func TestConstructFSM(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Schema
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  fields := findFields(pFields, "id", "names.languages.code", "names.languages.country")
  fsm := ConstructFSM(fields)
  fmt.Printf("%v\n", fsm)
  if 1 != 2 {
    t.Errorf("hi")
  }
}
