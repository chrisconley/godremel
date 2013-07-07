package main

import (
    "fmt"
    "strings"
    "os"
    "encoding/json"
    "io/ioutil"
)


type Field struct {
    Name string
    Kind string // int, string, record
    Mode string // optional, repeated
    Fields []Field
}

type ProcessedField struct {
    Name string // get rid of this?
    Path string
    Parent *ProcessedField
    // will have to add Mode on here, then we can write functions to get max Repetition and Definition level
    // Do we need Kind?
}

type Schema struct {
  Fields []Field
}

func processFields(fields []Field, processedFields []ProcessedField, parent ProcessedField) []ProcessedField {
  for _, field := range fields {

    path := ""
    if parent.Path != "" {
      path = strings.Join([]string{parent.Path, field.Name}, ".")
    } else {
      path = field.Name
    }

    processedField := ProcessedField{field.Name, path, &parent}
    processedFields = append(processedFields, processedField)
    if field.Fields != nil {
      processedFields = processFields(field.Fields, processedFields, processedField)
    }
  }
  return processedFields
}

func main() {

  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      fmt.Printf("File error: %v\n", e)
      os.Exit(1)
  }

  var schema Schema
  err := json.Unmarshal(file, &schema)
  if err != nil {
    fmt.Println("error:", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  for _, pField := range pFields {
    fmt.Printf("%v\n", pField)
  }
}
