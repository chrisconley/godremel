package go_dremel

import (
  "fmt"
)

type Record map[interface{}]interface{}
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
func (reader *Reader) HasData() bool {
  return reader.CurrentRow().Value != ""
}

func (reader *Reader) FetchNextRow() Row {
  row := reader.CurrentRow()

  //reader.CurrentRowIndex += 1
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
  reader := &Reader{}
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

func AssembleRecord(readers []*Reader) Record {
  record := Record{}

  //for i := 0; i < len(readers); i++ {
    //r := readers[i]
    //fmt.Printf("READERS: %v\n", r)
    //fmt.Printf("READERS: %p\n", &r)
  //}

  // this isn't right, but I'm not sure what the "root" field reader is
  // Maybe readers[0] is supposed to be for "id", and lastReader is for a "" or "root" reader
  //lastReader := readers[0]

  counter := 0
  reader := readers[0]
  fmt.Printf("~~~initreader: %v\n", reader)
  fmt.Printf("&&&initreader: %p\n", &reader)

  for reader.HasData() && counter < 20 {
    counter++
    fmt.Printf("~~~reader: %v\n", reader)
    fmt.Printf("&&&reader: %p\n", &reader)
    row := reader.FetchNextRow()
    reader.CurrentRowIndex += 1
    fmt.Printf("@@CurrentRowIndex: %v\n", reader.CurrentRowIndex)
    if row.Value != "" {
      //lastReader = moveToLevel(reader.TreeLevel(), reader, lastReader)
      //appendValue(record, reader)
    } else {
      //lastReader = moveToLevel(reader.FullDefinitionLevel(), reader, lastReader)
    }
    reader = reader.NextReader(readers)
    //lastReader = returnToLevel(reader.TreeLevel(), reader, lastReader)
  }
  //lastReader = returnToLevel(0)
  return record
}
