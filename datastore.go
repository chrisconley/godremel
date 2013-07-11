package go_dremel

type Column string

type Row struct {
  Value string
  RepetitionLevel int
  D int
}

type DataStore interface {
  ReadColumn(column Column) []Row
  WriteRow(column Column, row Row)
}

type MemStore struct {
  data map[Column][]Row
}

func (memStore *MemStore) ReadColumn(column Column) []Row  {
  return []Row{}
}

func (memStore *MemStore) WriteRow(column Column, row Row) {
}
