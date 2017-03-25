package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo/buffalo/cmd/generate"
	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd is the hook for all of the other commands in the buffalo binary.
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "buffalo",
	Short:         "Helps you build your Buffalo applications that much easier!",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Buffalo version %s\n\n", Version)
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

func goInstall(pkg string) *exec.Cmd {
	return generate.GoInstall(pkg)
}

func goGet(pkg string, buildFlags ...string) *exec.Cmd {
	return generate.GoGet(pkg, buildFlags...)
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
