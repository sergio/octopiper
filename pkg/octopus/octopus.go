package octopus

// Variable
type Variable struct {
	ID    string
	Name  string
	Scope map[string][]string
	Value string
}

type Entry struct {
	ID   string
	Name string
}

type Namespace []Entry

type VariableSet struct {
	ScopeValues map[string]Namespace
	Variables   []Variable
	OwnerID     string
}
