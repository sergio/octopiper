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

var (
	cfgFile string
	query   string

	globalconfig cli.Config
)

func octopusClient() *octopusclient.OctopusClient {

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

var rootCmd = &cobra.Command{
	Use:   "octopiper",
	Short: "Octopus Deploy command line utility",
	Long:  `Command line interface to Octopus Deploy REST API.`,
}

// Execute :
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.octopiper.yaml)")
	rootCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "JMESPath expression")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".octopiper")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	viper.Unmarshal(&globalconfig)
}
