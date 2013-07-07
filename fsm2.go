package main

import (
    //"fmt"
    "os"
    "encoding/json"
    "io/ioutil"
)


type Field struct {
    name string
    kind string // int, string, record
    mode string // optional, repeated
    fields []Field
}

type Schema struct {
  fields []Field
}

func main() {

  fmt
  os
  encoding.json
}
