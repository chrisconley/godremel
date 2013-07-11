package go_dremel

type Column string

type Row struct {
  Value interface{}
  RepetitionLevel int
  D int
}

type DataStore interface {
  ReadColumn(column string) []Row
  WriteRow(column string, row Row)
}

type MemStore struct {
  Data map[string][]Row
}

func (memStore *MemStore) ReadColumn(column string) []Row  {
  if _, present := memStore.Data[column]; present {
    return memStore.Data[column]
  } else {
    return []Row{}
  }
}

func (memStore *MemStore) WriteRow(column string, row Row) {
  rows := memStore.ReadColumn(column)
  rows = append(rows, row)
  memStore.WriteRows(column, rows)
}

func (memStore *MemStore) WriteRows(column string, rows []Row) {
  memStore.Data[column] = rows
}
