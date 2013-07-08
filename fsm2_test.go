package go_dremel

import (
    "fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func TestSqrt(t *testing.T) {
  fmt.Printf("hi")
  field1 := ProcessedField{}
  field2 := ProcessedField{}
  result := getCommonRepetitionLevel(field1, field2)
	if result != 2 {
		t.Errorf("Expected %v to be %v", result, 1)
	}
}

func TestAncestors(t *testing.T) {
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
  for _, pField := range pFields {
    fmt.Printf("%v\n", pField)
  }

  fmt.Printf("%v\n", pFields[6].Ancestors())
}
