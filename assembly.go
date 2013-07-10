package go_dremel

import (
  "fmt"
  "strings"
)

type RecordChildren map[interface{}]*Record

type Record struct {
  Name string
  Parent *Record
  Children RecordChildren
  Value interface{}
  Values []interface{}
}
type RepetitionLevel int

type Row struct {
  Value string
  RepetitionLevel int
  D int
}

func MakeReaders(columns map[string][]Row, fields []ProcessedField, fsm FSM2) []*Reader {
  readers := []*Reader{}

  for _, field := range fields {
    reader := Reader{field, columns[field.Path], fsm[field], 0}
    readers = append(readers, &reader)
  }

  return readers
}

type FieldFSM map[int]ProcessedField

type Reader struct {
  Field ProcessedField
  Rows []Row
  FSM FieldRepetitionLevelTransitions
  CurrentRowIndex int
}

var EmptyReader = &Reader{}

func (reader *Reader) HasData() bool {
  return reader.Field != ProcessedField{}
}

func (reader *Reader) FetchNextRow() Row {
  row := reader.CurrentRow()

  reader.CurrentRowIndex += 1
  return row
}

func (reader *Reader) CurrentRow() Row {
  if reader.CurrentRowIndex < len(reader.Rows) {
    return reader.Rows[reader.CurrentRowIndex]
  } else {
    return Row{}
  }
}

func (reader *Reader) NextRow() Row {
  nextIndex := reader.CurrentRowIndex + 1
  if nextIndex < len(reader.Rows) {
    return reader.Rows[nextIndex]
  } else {
    return Row{}
  }
}

func (reader *Reader) NextRepetionLevel() int {
  nextRow := reader.NextRow()
  fmt.Printf("NextRow: %v\n", reader.NextRow())
  return nextRow.RepetitionLevel
}

func findReaderByField(field ProcessedField, readers []*Reader) *Reader {
  reader := EmptyReader
  for i := 0; i < len(readers); i++ {
    r := readers[i]
    if r.Field == field {
      reader = r
    }
  }
  return reader
}


func (reader *Reader) NextReader(readers []*Reader) *Reader {
    destinationField := reader.FSM[reader.NextRepetionLevel()]
    fmt.Printf("NextRepetionLevel: %v\n", reader.NextRepetionLevel())
    destinationReader := findReaderByField(destinationField, readers)
    fmt.Printf("destinationReader: %v\n", destinationReader)
    return destinationReader
}

func countNonEmptyStrings(strings []string) int {
  count := 0
  for i := 0; i < len(strings); i++ {
    if strings[i] != "" {
      count++
    }
  }
  return count
}

func moveToLevel(record *Record, nextReader *Reader, lastReader *Reader, lowestCommonAncestor *Reader) (*Record) {
    commonPath := lowestCommonAncestor.Field.Path
    nextPath := nextReader.Field.Path
    lastPath := lastReader.Field.Path
    commonPaths := strings.Split(commonPath, ".")
    nextPaths := strings.Split(nextPath, ".")
    lastPaths := strings.Split(lastPath, ".")
    fmt.Printf("MOVING\n")
    fmt.Printf("RECORD: %v\n", record.Name)
    // end nested records up to lowest common ancestor of next and last reader
    for index := countNonEmptyStrings(lastPaths); index > countNonEmptyStrings(commonPaths); index-- {
      fmt.Printf("ENDING\n")
      fmt.Printf("Len commonPaths: %v\n", countNonEmptyStrings(commonPaths))
      fmt.Printf("LastPaths: %v, Len lastPaths: %v\n", lastPaths, countNonEmptyStrings(lastPaths))
      fmt.Printf("Index: %v\n", index)
      record = record.Parent
      fmt.Printf("RECORD: %v\n", record.Name)
    }

    // start nested records up from lowest common ancestor to nextReader.Path
    for index := countNonEmptyStrings(commonPaths); index < countNonEmptyStrings(nextPaths); index++ {
      fmt.Printf("STARTING\n")
      fmt.Printf("Len commonPaths: %v\n", countNonEmptyStrings(commonPaths))
      fmt.Printf("Len nextPaths: %v\n", countNonEmptyStrings(nextPaths))
      fmt.Printf("Index: %v\n", index)
      name := nextPaths[index]
      record.Children[name] = &Record{Name:name, Children:RecordChildren{}, Parent: record}
      record = record.Children[name]
      fmt.Printf("RECORD: %v\n", record.Name)
    }

    // set lastReader to one at newLevel
    lastReader = nextReader

    return record
}

func returnToLevel(record *Record, nextReader *Reader, lastReader *Reader, lowestCommonAncestor *Reader) (*Record) {
    commonPath := lowestCommonAncestor.Field.Path
    nextPath := nextReader.Field.Path
    lastPath := lastReader.Field.Path
    commonPaths := strings.Split(commonPath, ".")
    nextPaths := strings.Split(nextPath, ".")
    lastPaths := strings.Split(lastPath, ".")
    fmt.Printf("LowestCommonAncestor Path: %v\n", commonPaths)
    fmt.Printf("Next Path: %v\n", nextPaths)
    fmt.Printf("Last Path: %v\n", lastPaths)
    fmt.Printf("RECORD: %v\n", record.Name)

    // end nested records up to lowest common ancestor of next and last reader
    for index := countNonEmptyStrings(lastPaths); index > countNonEmptyStrings(commonPaths); index-- {
      fmt.Printf("ENDING\n")
      fmt.Printf("Len commonPaths: %v\n", countNonEmptyStrings(commonPaths))
      fmt.Printf("Len lastPaths: %v\n", countNonEmptyStrings(nextPaths))
      fmt.Printf("Index: %v\n", index)
      record = record.Parent
      fmt.Printf("RECORD: %v\n", record.Name)
    }

    // set lastReader to one at newLevel
    lastReader = nextReader

    return record
}

func getLowestCommonReaderAncestor(r1 *Reader, r2 *Reader, readers []*Reader) *Reader {
  commonFieldAncestor := GetLowestCommonAncestor(r1.Field, r2.Field)
  return findReaderByField(commonFieldAncestor, readers)
}

func appendValue(record *Record, reader *Reader, value string) {
  if reader.Field.Mode == "repeated" {
    if record.Values == nil {
      record.Values = make([]interface{}, 0, 100)
      record.Values[0] = value
    } else {
      record.Values = append(record.Values, value)
    }
  } else {
    record.Value = value
  }
}

func AssembleRecord(readers []*Reader) *Record {
  record := &Record{Name: "root", Children:RecordChildren{}}

  rootReader := EmptyReader
  lastReader := rootReader

  counter := 0
  reader := readers[0]

  for reader.HasData() && counter < 20 {
    counter++
    row := reader.FetchNextRow()
    lowestCommonAncestor := getLowestCommonReaderAncestor(reader, lastReader, readers)
    if row.Value != "" {
      record = moveToLevel(record, reader, lastReader, lowestCommonAncestor)
      appendValue(record, reader, row.Value)
    } else {
      record = moveToLevel(record, reader, lastReader, lowestCommonAncestor)
    }
    reader = reader.NextReader(readers)
    if (reader != EmptyReader) {
      lowestCommonAncestor = getLowestCommonReaderAncestor(reader, lastReader, readers)
      fmt.Printf("ALMOST FINAL RETURN\n")
      fmt.Printf("reader%v\n", reader)
      record = returnToLevel(record, reader, lastReader, lowestCommonAncestor)
    }
  }
  fmt.Printf("FINAL RETURN\n")
  fmt.Printf("~~~lastreader: %v\n", lastReader)
  record = returnToLevel(record, rootReader, lastReader, EmptyReader)
  return record
}
