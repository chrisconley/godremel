package go_dremel


type Record map[interface{}]interface{}
type RepetitionLevel int

type Row struct {
  Value string
  RepetitionLevel RepetitionLevel
  D int
}

type Reader struct {
  Field ProcessedField
  Rows []Row
  CurrentRowIndex int
}
func (reader *Reader) HasData() bool {
  return reader.CurrentRow().Value != ""
}

func (reader *Reader) FetchNextRow() Row {
  reader.CurrentRowIndex += 1
  return reader.CurrentRow()
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

func (reader *Reader) NextRepetionLevel() RepetitionLevel {
  nextRow := reader.NextRow()
  return nextRow.RepetitionLevel
}


func (reader *Reader) NextReader() Reader {
    //field := fsm[FsmState{reader.field, reader.NextRepetionLevel()}
    //somehow find reader by field
    //or change fsm so that states are nested and Reader has its own transitions
    return Reader{}
}

func AssembleRecord(readers []Reader) Record {
  record := Record{}

  // this isn't right, but I'm not sure what the "root" field reader is
  // Maybe readers[0] is supposed to be for "id", and lastReader is for a "" or "root" reader
  //lastReader := readers[0]

  reader := readers[0]

  for reader.HasData() {
    row := reader.FetchNextRow()
    if row.Value != "" {
      //lastReader = moveToLevel(reader.TreeLevel(), reader, lastReader)
      //appendValue(record, reader)
    } else {
      //lastReader = moveToLevel(reader.FullDefinitionLevel(), reader, lastReader)
    }
    //reader = reader.NextReader()
    //lastReader = returnToLevel(reader.TreeLevel(), reader, lastReader)
  }
  //lastReader = returnToLevel(0)
  return record
}
