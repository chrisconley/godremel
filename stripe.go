package go_dremel

import (
  "fmt"
)

type Decoder struct {
  Field Field
  //record map[string]interface{}
  //records map[string][]interface{}
  Record interface{}
}

type FieldValue struct {
  Field Field
  Value interface{}
}

func (decoder *Decoder) ReadValues() chan FieldValue  {
  c := make(chan FieldValue)
  go func() {
    for _, f := range decoder.Field.Fields {
      //if f.Mode == "repeated" {
      switch vv := decoder.Record.(type) {
      case map[string]interface{}:
        value := decoder.Record[decoder.Field.Name]
        c <- FieldValue{decoder.Field, value}
      case map[string][]interface{}:
        for _, value := range decoder.Record[decoder.Field.Name] {
          c <- FieldValue{decoder.Field, value}
        }
      }
    }
    close(c)
  }()
  return c
}

func StripeRecord(field Field, record interface{}, datastore *DataStore) {
  seenFields := []Field{}
  decoder := Decoder{field, record}
  for fieldValue := range decoder.ReadValues() {
    fmt.Printf("Field: %v, Value: %v\n", fieldValue.Field.Name, fieldValue.Value)

  }
}
