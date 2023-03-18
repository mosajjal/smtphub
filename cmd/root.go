package cmd

import (
	"fmt"
	"log"
	"os"

	smtphub "github.com/mosajjal/smtphub/pkg"
	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config smtphub.Config

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "smtphub",
		Short: "smtphub is a CLI tool for creating SMTP hooks",
		Long:  `smtphub provides a simple lexer to parse SMTP messages (emails) and decide which script to run on them.`,
		Run: func(cmd *cobra.Command, args []string) {
			smtphub.Run(config)
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smtphub.yaml)")

	// set up logging
	// set log level
	if l, err := zerolog.ParseLevel("debug"); err == nil {
		zerolog.SetGlobalLevel(l)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".smtphub")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s\n", err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
