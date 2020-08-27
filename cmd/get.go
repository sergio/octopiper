package cmd

import (
	"octopiper/pkg/cli"
	"octopiper/pkg/jsonfilter"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get json from Octopus Deploy API resources",
	Long: `Get downloads json as result of accesing any Octopus Deploy resource.
	Combined with the global --query option allows easy json parsing and
	processing using JMESPath expressions.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		resource := args[0]
		jsondata, err := octopusClient().GetJSON(resource)
		if err != nil {
			return err
		}

		jsondata, err = jsonfilter.Query(query, jsondata)
		if err != nil {
			return err
		}

		cli.WriteOutput(jsondata)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
