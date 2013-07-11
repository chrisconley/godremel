package go_dremel

import (
    "fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func TestAssembly(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  fields := findFields(pFields, "id", "names.languages.country")
  fsm := ConstructFSM(fields)

  //fmt.Printf("FIELDS: %v\n", fields)

  columns := map[string][]Row{}
  columns["id"] = []Row{
    Row{"1", 0, 0},
    Row{"2", 0, 0},
  }
  columns["names.languages.country"] = []Row{
    Row{"us", 0, 3}, Row{"", 2, 2}, Row{"", 1, 1},
    Row{"gb", 1, 3}, Row{"", 0, 1},
  }

  readers := MakeReaders(columns, fields, fsm)
  //fmt.Printf("Readers: %s\n", readers)
  //for _, reader := range readers {
     //fmt.Printf("**reader field: %s\n", reader.Field)
  //}
  record := AssembleRecord(readers)
  fmt.Printf("record: %v\n", *record)
  if 1 != 1 {
    t.Errorf("hi")
  }

  record.ToMap()

}
