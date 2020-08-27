package cmd

import (
	"encoding/json"
	"fmt"
	"octopiper/pkg/cli"
	"octopiper/pkg/octopus"

	"github.com/spf13/cobra"
)

// findCmd represents the findVariables command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		searchTerm := args[0]
		results, err := find(searchTerm)
		cli.WriteOutput(results)
		return err

	},
}

type variableSetSummary struct {
	VariableSetID string `json:"VariableSetId"`
	Name          string `json:"Name"`
}

func find(searchTerm string) ([]octopus.Variable, error) {
	// In-memory octopus information model
	octopusModel := octopus.NewModel()

	var variableSets []variableSetSummary

	// Get list of project variable sets
	jsontext, err := OctopusClient().Get("projects/all")
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
	jsontext, err = OctopusClient().Get("libraryvariablesets/all")
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
		jsonString, err := OctopusClient().Get(fmt.Sprintf("variables/%s", summary.VariableSetID))
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

	return searchResults, nil
}

func init() {
	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
