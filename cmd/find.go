package cmd

import (
	"encoding/json"
	"fmt"
	"octopiper/pkg/cli"
	"octopiper/pkg/jsonfilter"
	"octopiper/pkg/octopus"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find octopus variables by value in any variable set.",
	Long:  `Finds and ouptupts octopus variables containing the specified text.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		searchTerm := args[0]
		jsondata, err := find(searchTerm)
		jsondata, err = jsonfilter.Query(query, jsondata)
		if err != nil {
			return err
		}

		cli.WriteOutput(jsondata)
		return err

	},
}

type variableSetSummary struct {
	VariableSetID string `json:"VariableSetId"`
	Name          string `json:"Name"`
}

func find(searchTerm string) (interface{}, error) {
	// In-memory octopus information model
	octopusModel := octopus.NewModel()

	var variableSets []variableSetSummary

	// Get list of project variable sets
	jsontext, err := octopusClient().Get("projects/all")
	if err != nil {
		return nil, err
	}

	var summaries []variableSetSummary
	err = json.Unmarshal([]byte(jsontext), &summaries)
	if err != nil {
		return nil, err
	}

	variableSets = append(variableSets, summaries...)

	// Get list of library variable sets
	jsontext, err = octopusClient().Get("libraryvariablesets/all")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsontext), &summaries)
	if err != nil {
		return nil, err
	}

	variableSets = append(variableSets, summaries...)

	// Download all variable sets
	var searchResults []octopus.Variable
	for _, summary := range variableSets {
		jsonString, err := octopusClient().Get(fmt.Sprintf("variables/%s", summary.VariableSetID))
		if err != nil {
			return nil, err
		}

		// add to in-memory information model
		octopusModel.AddVariableSet(summary.Name, jsonString)
	}

	// perform search for term
	searchResults, err = octopusModel.FindVariables(searchTerm)
	if err != nil {
		return nil, err
	}

	jsonbytes, err := json.Marshal(searchResults)
	if err != nil {
		return nil, err
	}
	var jsondata interface{}
	err = json.Unmarshal(jsonbytes, &jsondata)
	if err != nil {
		return nil, err
	}

	return jsondata, nil
}

func init() {
	rootCmd.AddCommand(findCmd)
}
