package cmd

import (
	"fmt"
	"octopiper/pkg/cli"
	"octopiper/pkg/jsonfilter"
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

func find(searchTerm string) ([]octopus.Variable, error) {
	// In-memory octopus information model
	octopusModel := octopus.NewModel()

	var variableSetIds []string

	// Get list of project variable sets

	jsondata, err := OctopusClient().GetJSON("projects/all")
	if err != nil {
		return nil, err
	}

	jsondata, err = jsonfilter.Query("[].VariableSetId", jsondata)
	if err != nil {
		return nil, err
	}

	switch t := jsondata.(type) {
	case []interface{}:
		for _, s := range t {
			variableSetIds = append(variableSetIds, s.(string))
		}
	default:
		return nil, fmt.Errorf("Unexpected json structure: %#v", jsondata)
	}

	// Get list of library variable sets
	jsondata, err = OctopusClient().GetJSON("libraryvariablesets/all")

	if err != nil {
		return nil, err
	}

	jsondata, err = jsonfilter.Query("[].VariableSetId", jsondata)
	if err != nil {
		return nil, err
	}

	switch t := jsondata.(type) {
	case []interface{}:
		for _, s := range t {
			variableSetIds = append(variableSetIds, s.(string))
		}
	default:
		return nil, fmt.Errorf("Unexpected json structure: %#v", jsondata)
	}

	// Download all variable sets
	var searchResults []octopus.Variable
	for _, setID := range variableSetIds {
		jsonString, err := OctopusClient().Get(fmt.Sprintf("variables/%s", setID))
		if err != nil {
			return nil, err
		}

		// add to in-memory information model
		octopusModel.Add(jsonString)
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
