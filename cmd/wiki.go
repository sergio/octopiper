package cmd

import (
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var wikiCmd = &cobra.Command{
	Use:   "wiki",
	Short: "Opens the Octopus API Wiki in the default browser",
	Long:  `Opens a new tab in the default browser with the wiki page for the Octopus Deploy API documentation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return browser.OpenURL(`https://github.com/OctopusDeploy/OctopusDeploy-Api/wiki`)
	},
}

func init() {
	rootCmd.AddCommand(wikiCmd)
}
