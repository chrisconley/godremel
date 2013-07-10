package go_dremel

import (
    "strings"
)


type Field struct {
    Name string
    Kind string // int, string, record
    Mode string // optional, repeated
    Fields []Field
}

type ProcessedField struct {
    Name string // get rid of this?
    Path string
    Mode string
    Parent *ProcessedField
    // will have to add Mode on here, then we can write functions to get max Repetition and Definition level
    // Do we need Kind?
}

func (processedField *ProcessedField) Ancestors() []ProcessedField {
    parent := processedField.Parent
    if parent != nil && parent.Name != "" {
      ancestors := []ProcessedField{*parent}
      return append(ancestors, parent.Ancestors()...)
    } else {
      return []ProcessedField{}
    }
}

func (field *ProcessedField) MaxRepetitionLevel() int {
    maxRepetitionLevel := 0
    if field.Mode == "repeated" {
      maxRepetitionLevel += 1
    }
    for _, a := range field.Ancestors() {
      if a.Mode == "repeated" {
        maxRepetitionLevel += 1
      }
    }
    return maxRepetitionLevel
}

func findField (path string, fields []ProcessedField) ProcessedField {
  returnField := ProcessedField{}
  for _, field := range fields {
    if field.Path == path {
      returnField = field
    }
  }
  return returnField
}

func makeStringSet (strings []string) map[string]bool {
  set := map[string]bool{}
  for _, s := range strings {
    set[s] = true
  }
  return set
}

func findFields (fields []ProcessedField, paths ...string) []ProcessedField {
  pathSet := makeStringSet(paths)
  returnFields := []ProcessedField{}
  for _, field := range fields {
    if pathSet[field.Path] {
      returnFields = append(returnFields, field)
    }
  }
  return returnFields
}

type Schema struct {
  Fields []Field
}

func GetCommonRepetitionLevel(f1 ProcessedField, f2 ProcessedField) int {
  commonAncestors := []ProcessedField{}
  for _, a1 := range f1.Ancestors() {
    a2 := findField(a1.Path, f2.Ancestors())
    if a2.Path != "" {
      commonAncestors = append(commonAncestors, a2)
    }

  }

  maxRepetitionLevel := 0
  for _, a := range commonAncestors {
    if a.MaxRepetitionLevel() > maxRepetitionLevel {
      maxRepetitionLevel = a.MaxRepetitionLevel()
    }
  }
  return maxRepetitionLevel
}

func processFields(fields []Field, processedFields []ProcessedField, parent ProcessedField) []ProcessedField {
  for _, field := range fields {

    path := ""
    if parent.Path != "" {
      path = strings.Join([]string{parent.Path, field.Name}, ".")
    } else {
      path = field.Name
    }

    processedField := ProcessedField{field.Name, path, field.Mode, &parent}
    processedFields = append(processedFields, processedField)
    if field.Fields != nil {
      processedFields = processFields(field.Fields, processedFields, processedField)
    }
  }
  return processedFields
}
