package octopus

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Variable :
type Variable struct {
	ID              string
	Name            string
	Scope           map[string][]string
	Value           string
	VariableSetID   string
	VariableSetName string
}

// Entry :
type Entry struct {
	ID   string
	Name string
}

// Namespace :
type Namespace []Entry

// VariableSet :
type VariableSet struct {
	VariableSetName string
	ScopeValues     map[string]Namespace
	Variables       []Variable
	OwnerID         string
}

// Model :
type Model struct {
	VariableSets      []VariableSet
	ScopeDictionaries map[string]map[string]string // (category/id/name)
}

// NewModel :
func NewModel() *Model {
	m := &Model{}
	m.ScopeDictionaries = make(map[string]map[string]string)
	return m
}

// AddVariableSet :
func (m *Model) AddVariableSet(variableSetName string, jsonString string) error {
	var vs VariableSet
	err := json.Unmarshal([]byte(jsonString), &vs)
	if err != nil {
		return fmt.Errorf("error unmarshaling json as VariableSet %w", err)
	}
	vs.VariableSetName = variableSetName
	m.VariableSets = append(m.VariableSets, vs)
	for category, ns := range vs.ScopeValues {
		if m.ScopeDictionaries[category] == nil {
			m.ScopeDictionaries[category] = map[string]string{}
		}
		for _, entry := range ns {
			m.ScopeDictionaries[category][entry.ID] = entry.Name
		}
	}
	return nil
}

// FindVariables :
func (m *Model) FindVariables(searchTerm string) ([]Variable, error) {
	results := []Variable{}

	for _, vs := range m.VariableSets {
		for _, v := range vs.Variables {
			v.VariableSetID = vs.OwnerID
			v.VariableSetName = vs.VariableSetName
			if strings.Contains(v.Value, searchTerm) {
				results = append(results, m.translateScope(v))
			}
		}
	}

	return results, nil
}

func (m *Model) translateScope(variable Variable) Variable {
	for category, ids := range variable.Scope {
		variable.Scope[category] = translateWords(ids, m.ScopeDictionaries[pluralized(category)])
	}
	return variable
}

func pluralized(singular string) string {
	return fmt.Sprintf("%ss", singular)
}

func translateWords(words []string, dictionary map[string]string) []string {
	result := make([]string, len(words))
	for i, word := range words {
		translation, found := dictionary[word]
		if found {
			result[i] = translation
		} else {
			result[i] = word
		}
	}
	return result
}
