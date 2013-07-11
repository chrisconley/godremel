package go_dremel

import (
  "fmt"
  //"reflect"
)

type FieldValue struct {
  Field Field
  Value interface{}
}

type Decoder struct {
  Field Field
  Record interface{}
}

func (decoder *Decoder) ReadValues() chan FieldValue  {
  c := make(chan FieldValue)
  go func() {
    for _, f := range decoder.Field.Fields {
      recordValue := decoder.getValue(f.Name)
      if f.Mode == "repeated" && recordValue != nil {
        for _, value := range recordValue.([]interface{}) {
          c <- FieldValue{f, value}
        }
      } else {
        c <- FieldValue{f, recordValue}
      }
      //switch rType := recordValue.(type) {
      //case []interface{}:
        //for _, value := range recordValue.([]interface{}) {
          //c <- FieldValue{f, value}
        //}
      //case interface{}:
        //c <- FieldValue{f, recordValue}
      ////case nil:
        ////c <- FieldValue{f, ""}
      //default:
        //fmt.Printf("Mystery Type: %v\n", rType)
      //}
    }
    close(c)
  }()
  return c
}

func (decoder *Decoder) getValue(fieldName string) interface{} {
  if decoder.Record != nil {
    return decoder.Record.(map[string]interface{})[fieldName]
  } else {
    return nil
  }
}

func StripeRecord(field Field, record interface{}, datastore DataStore) {
  //seenFields := []Field{}
  //fmt.Printf("RECORD: %v\n", record)
  decoder := Decoder{field, record}
  for fieldValue := range decoder.ReadValues() {
    //fmt.Printf("Field: %v, Value: %v\n", fieldValue.Field, fieldValue.Value)
    if fieldValue.Field.Kind == "record" {
      StripeRecord(fieldValue.Field, fieldValue.Value, datastore)
    } else {
      fmt.Printf("Field: %v, Value: %v\n", fieldValue.Field, fieldValue.Value)
    }

  }
}
