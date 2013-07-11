package go_dremel

import (
    //"fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func TestAncestors(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  //for _, pField := range pFields {
    //fmt.Printf("%v\n", pField)
  //}

  ancestors := pFields[6].Ancestors()

  if len(ancestors) != 2 {
    t.Errorf("Ancestor error: %v\n", ancestors)
  }

  if ancestors[0].Path != "names.languages" || ancestors[1].Path != "names" {
    t.Errorf("Ancestor error: %v\n", ancestors)
  }
}

func TestFindField(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  field := findField("names.languages.code", pFields)
  if field.Path != "names.languages.code" {
    t.Errorf("Field error: %v\n", field)
  }
}

func TestGetCommonRepetitionLevel(t * testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})

  codeField := findField("names.languages.code", pFields)
  countryField := findField("names.languages.country", pFields)
  urlField := findField("names.url", pFields)

  rLevel := 0

  rLevel = GetCommonRepetitionLevel(codeField, countryField)
  if rLevel != 2 {
    t.Errorf("Repetition level error: %v\n", rLevel)
  }

  rLevel = GetCommonRepetitionLevel(codeField, urlField)
  if rLevel != 1 {
    t.Errorf("Repetition level error: %v\n", rLevel)
  }

}


// Trying generic collect function: http://play.golang.org/p/5Osey24Itz
