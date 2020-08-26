package cmd

import (
	"fmt"
	"octopiper/pkg/cli"
	"octopiper/pkg/octopusapi"
	"octopiper/pkg/octopusclient"
	"octopiper/pkg/sqlitekvs"
	"os"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var globalconfig cli.Config

// OctopusClient //
func OctopusClient() *octopusclient.OctopusClient {
	client := &octopusclient.OctopusClient{
		LocalCache: &sqlitekvs.KVSCache{
			FilePath: globalconfig.LocalCache.FilePath,
			TTL:      time.Duration(globalconfig.LocalCache.TTLMinutes) * time.Minute,
		},
		OctopusServer: &octopusapi.OctopusAPI{
			BaseURL: globalconfig.OctopusServer.BaseURL,
			APIKey:  globalconfig.OctopusServer.APIKey,
		},
	}
	return client
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "octopiper",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.octopiper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".octopiper" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".octopiper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	viper.Unmarshal(&globalconfig)
}
