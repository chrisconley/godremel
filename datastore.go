package go_dremel

type Column string

type Row struct {
  Value string
  RepetitionLevel int
  D int
}

type DataStore interface {
  ReadColumn(column string) []Row
  WriteRow(column string, row Row)
}

type MemStore struct {
  data map[string][]Row
}

func (memStore *MemStore) ReadColumn(column string) []Row  {
  return []Row{}
}

func (memStore *MemStore) WriteRow(column string, row Row) {
}
