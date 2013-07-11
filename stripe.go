package go_dremel

import (
  "strings"
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

type Writer struct {
  Name string
  Field Field
  Value interface{}
  Parent *Writer
}

var RootWriter = Writer{"__root__", Field{}, nil, nil}

func (writer *Writer) RepeatedFieldDepth() int {
  depth := 0
  if writer.Field.Mode == "repeated" {
    depth++
  }
  parent := writer.Parent
  for parent.Name != RootWriter.Name {
    if parent.Field.Mode == "repeated" {
      depth++
    }
    parent = parent.Parent
  }
  return depth
}

func (writer *Writer) DefinitionLevel() int {
  depth := 0
  if writer.Field.Mode != "required" && writer.Value != nil {
    depth++
  }
  parent := writer.Parent
  for parent.Name != RootWriter.Name {
    if parent.Field.Mode != "required" && parent.Value != nil {
      depth++
    }
    parent = parent.Parent
  }
  return depth
}

func (writer *Writer) Path() string {
  path := ""
  if writer.Parent != nil && writer.Parent.Path() != "" {
    path = strings.Join([]string{writer.Parent.Path(), writer.Field.Name}, ".")
  } else {
    path = writer.Field.Name
  }
  return path
}

func StripeRecord(field Field, record interface{}, datastore DataStore, writer Writer, rLevel int) {
  seenFields := map[string]bool{}
  decoder := Decoder{field, record}
  for fieldValue := range decoder.ReadValues() {
    childWriter := Writer{fieldValue.Field.Name, fieldValue.Field, fieldValue.Value, &writer}
    childRepetitionLevel := rLevel

    // if we've seen this field already
    if _, present := seenFields[fieldValue.Field.Name]; present {
       childRepetitionLevel = childWriter.RepeatedFieldDepth()
    } else {
      seenFields[fieldValue.Field.Name] = true
    }

    if fieldValue.Field.Kind == "record" {
      StripeRecord(fieldValue.Field, fieldValue.Value, datastore, childWriter, childRepetitionLevel)
    } else {
      row := Row{childWriter.Value, childRepetitionLevel, childWriter.DefinitionLevel()}
      datastore.WriteRow(childWriter.Path(), row)
    }

  }
}
