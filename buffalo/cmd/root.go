package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string

// RootCmd is the hook for all of the other commands in the buffalo binary.
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "buffalo",
	Short:         "Helps you build your Buffalo applications that much easier!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Buffalo version %s\n\n", Version)

		anywhereCommands := []string{"new", "version", "info"}
		isFreeCommand := false
		for _, freeCmd := range anywhereCommands {
			if freeCmd == cmd.Name() {
				isFreeCommand = true
			}
		}

		if isFreeCommand {
			return nil
		}

		if insideBuffaloProject() == false {
			return errors.New("you need to be inside your buffalo project path to run this command")
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(-1)
	}
}

func init() {
	decorate("root", RootCmd)
}

func insideBuffaloProject() bool {
	if _, err := os.Stat(".buffalo.dev.yml"); err != nil {
		return false
	}

	return true
}

// func init() {
// cobra.OnInitialize(initConfig)
// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.buffalo.yaml)")
// }

// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}
//
// 	viper.SetConfigName(".buffalo") // name of config file (without extension)
// 	viper.AddConfigPath("$HOME")    // adding home directory as first search path
// 	viper.AutomaticEnv()            // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
