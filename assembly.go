package go_dremel


type Record map[interface{}]interface{}

type Row struct {
  Value string
  R int
  D int
}

type Reader struct {
  Field ProcessedField
  data []Row
}

func AssembleRecord(readers []Reader) Record {
  record := Record{}
  return record
}
