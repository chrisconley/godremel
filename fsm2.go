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
}

type Schema struct {
  Fields []Field
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
  fmt.Printf("%s\n", schema)
}
