package cmd

import (
	"fmt"
	"octopiper/pkg/cli"
	"octopiper/pkg/jsonfilter"

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

		var variableSetIds []string

		// Get list of project variable sets

		jsondata, err := globalconfig.OctopusServer.GetJSON("projects/all")
		if err != nil {
			return err
		}

		jsondata, err = jsonfilter.Query("[].VariableSetId", jsondata)
		if err != nil {
			return err
		}

		switch t := jsondata.(type) {
		case []interface{}:
			for _, s := range t {
				variableSetIds = append(variableSetIds, s.(string))
			}
		default:
			return fmt.Errorf("Unexpected json structure: %#v", jsondata)
		}

		// Get list of library variable sets
		jsondata, err = globalconfig.OctopusServer.GetJSON("libraryvariablesets/all")

		if err != nil {
			return err
		}

		jsondata, err = jsonfilter.Query("[].VariableSetId", jsondata)
		if err != nil {
			return err
		}

		switch t := jsondata.(type) {
		case []interface{}:
			for _, s := range t {
				variableSetIds = append(variableSetIds, s.(string))
			}
		default:
			return fmt.Errorf("Unexpected json structure: %#v", jsondata)
		}

		var searchResults []interface{}
		for _, setID := range variableSetIds {
			jsondata, err = globalconfig.OctopusServer.GetJSON(fmt.Sprintf("variables/%s", setID))
			if err != nil {
				return err
			}
			expression := fmt.Sprintf("Variables[? Value != null]|[?contains(Value, `\"%s\"`)].{VariableSetId: `\"%s\"`,Name: Name, Value: Value}", searchTerm, setID)
			jsondata, err = jsonfilter.Query(expression, jsondata)
			if err != nil {
				return err
			}
			items := jsondata.([]interface{})
			if items != nil {
				searchResults = append(searchResults, items...)
			}
		}
		cli.WriteOutput(searchResults)

		return nil
	},
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
