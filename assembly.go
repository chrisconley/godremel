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

func AssembleRecord(readers []Reader) Record {
  record := Record{}

  // this isn't right, but I'm not sure what the "root" field reader is
  // Maybe readers[0] is supposed to be for "id", and lastReader is for a "" or "root" reader
  lastReader := readers[0]

  reader := readers[0]

  for reader.HasData() {
    row := reader.FetchNextRow()
    if row.Value != "" {
      //moveToLevel(reader.TreeLevel(), reader)
      //appendValue(record, reader)
    } else {
      //moveToLevel(reader.FullDefinitionLevel(), reader)
    }



  }
  return record
}
