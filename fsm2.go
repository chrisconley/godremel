package main

import (
    "fmt"
    "os"
    "encoding/json"
    "io/ioutil"
)


type Field struct {
    Name string
    Kind string // int, string, record
    Mode string // optional, repeated
    Fields []Field
    Parent *Field
}

type Schema struct {
  Fields []Field
}

func processFields(fields []Field, processedFields []Field) []Field {
  for _, field := range fields {
    fmt.Printf("%s\n", field.Name)
    processedFields = append(processedFields, field)
    if field.Fields != nil {
      processFields(field.Fields, processedFields)
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
  //fmt.Printf("%s\n", schema.Fields)
  fmt.Printf("%s\n", processFields(schema.Fields, []Field{}))
}
