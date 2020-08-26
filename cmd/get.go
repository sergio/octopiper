package cmd

import (
	"octopiper/pkg/cli"
	"octopiper/pkg/jsonfilter"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		resource := args[0]
		jsondata, err := OctopusClient().GetJSON(resource)
		if err != nil {
			return err
		}

		jsondata, err = jsonfilter.Query(*query, jsondata)
		if err != nil {
			return err
		}

		cli.WriteOutput(jsondata)

		return nil
	},
}

var query *string

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	query = getCmd.PersistentFlags().StringP("query", "q", "", "JMESPath expression")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
