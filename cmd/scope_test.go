package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSome(label string, ids ...int) []interface{} {
	items := []interface{}{}
	for id := range ids {
		key := fmt.Sprintf("%s%d", label, id)
		items = append(items, map[string]interface{}{"Id": key, "Name": key + "-Name"})
	}
	return items
}

type Namespace map[string]string

type NamesCatalog struct {
	categories map[string]Namespace
}

func toNamespace(items interface{}) (Namespace, error) {

	var ns Namespace

	list, ok := items.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected '[]interface{}' got '%T'", items)
	}

	for _, item := range list {
		n, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Expected 'map[string]interface{}' got '%T'", item)
		}

		itemID, ok := n["Id"].(string)
		if !ok {
			return nil, fmt.Errorf("Expected 'string' got '%T'", n["Id"])
		}

		itemValue, ok := n["Value"].(string)
		if !ok {
			return nil, fmt.Errorf("Expected 'string' got '%T'", n["Value"])
		}

		ns[itemID] = itemValue
	}
	return ns, nil
}

// AddScopeValues imports the ScopeValues section from Octopus results to the names catalog
func (c NamesCatalog) AddScopeValues(scopeValues map[string]interface{}) error {

	for categoryName, items := range scopeValues {
		ns, err := toNamespace(items)
		if err != nil {
			return err
		}
		c.categories[categoryName] = ns
	}
	return nil
}

// GetName translates an Id to a name given a category or the provided Id if no name found
func (c NamesCatalog) GetName(category, id string) string {
	ns, ok := c.categories[category]
	if !ok {
		return id
	}
	return ns[id]
}

func TestAddSimpleScopeValues(t *testing.T) {

	testCases := []struct {
		itemCategory   string
		itemID         string
		itemName       string
		lookupCategory string
		expectedFound  bool
	}{
		{"Category", "SomeItemId", "Some Item Name", "Category", true},
		//{"Category", "SomeItemId", "Some Item Name", "AnotherCategory", false},
		//{"Environments", "Environment-1", "QA", "Environments", true},
	}

	for _, tc := range testCases {

		// assume minimal scope values
		scopeValues := map[string]interface{}{
			tc.itemCategory: []interface{}{
				map[string]interface{}{"Id": tc.itemID, "Name": tc.itemName},
			},
		}

		// execute: import the scope values to a new catalog
		sut := NamesCatalog{}
		sut.AddScopeValues(scopeValues)

		// verify: use catalog to translate name and verify result
		actualName := sut.GetName(tc.lookupCategory, tc.itemID)
		if tc.expectedFound {
			assert.Equal(t, tc.itemName, actualName)
		} else {
			assert.Equal(t, tc.itemID, actualName)
		}
	}
}

/*

func TestScopes(t *testing.T) {

	scopeValues1 := map[string]interface{}{
		"Actions":      getSome("Action", 1, 2, 3),
		"Channels":     []interface{}{},
		"Environments": getSome("Environment", 1, 2, 3),
		"Machines":     getSome("Machine", 1, 2, 3),
		"Roles":        getSome("Role", 1, 2, 3),
	}

	scopeValues2 := map[string]interface{}{
		"Actions":      getSome("Action", 1, 2, 3),
		"Channels":     []interface{}{},
		"Environments": getSome("Environment", 1, 2, 3),
		"Machines":     getSome("Machine", 1, 2, 3),
		"Roles":        getSome("Role", 1, 2, 3),
	}

	unifiedScopeValues := map[string]interface{}{}
	unifiedScopeValues := mergeScopeValues(unifiedScopeValues, scopeValues1)

	output, err := json.MarshalIndent(scopeValues, "", "  ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(output))

	   ScopeValues := map[string]interface {}{
	   	"Actions":[]interface {}{
	   	},
	   	"Channels":[]interface {}{
	   	},
	   	"Environments":[]interface {}{
	   		map[string]interface {}{"Id":"Environments-1", "Name":"QA"},
	   		map[string]interface {}{"Id":"Environments-41", "Name":"SIMSAT"},
	   		map[string]interface {}{"Id":"Environments-3","Name":"Staging"},
	   		map[string]interface {}{"Id":"Environments-4","Name":"Piloto"},
	   		map[string]interface {}{"Id":"Environments-21", "Name":"LAB"},
	   		map[string]interface {}{"Id":"Environments-2", "Name":"PRODAWS"},
	   		map[string]interface {}{"Id":"Environments-5", "Name":"Intranet"}
	   	},
	   	"Machines":[]interface {}{
	   		map[string]interface {}{"Id":"Machines-167", "Name":"ARCHIVING-IS"},
	   		map[string]interface {}{"Id":"Machines-221", "Name":"ARCHIVING-NFS-WORKER"},
	   		map[string]interface {}{"Id":"Machines-168", "Name":"ARCHIVING-WEB"},
	   		map[string]interface {}{"Id":"Machines-164", "Name":"arch-linux-webapi"},
	   		map[string]interface {}{"Id":"Machines-197", "Name":"BLOBSTORAGE01"},
	   		map[string]interface {}{"Id":"Machines-187", "Name":"blobstorage01.reachcore.io"}
	   	}
	   }

*/
