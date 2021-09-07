package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GoPomodoro",
	Short: "Pomodoro in CLI",
	Long:  `Go pomodoro app in CLI the best app in the world`,
	//Run: func(cmd *cobra.Command, args []string) {
	//},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.GoPomodoro.env)")

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

		// Search config in home directory with name ".GoPomodoro.env".
		viper.AddConfigPath(home)
		viper.SetConfigName(".GoPomodoro")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, using default values")
		} else {
			fmt.Println("Config file found but err when reading file:", viper.ConfigFileUsed(), err)
		}
	}

	viper.SetDefault("SLACK_TOKEN", "")
	viper.SetDefault("ENABLE_SLACK_DND", false)
	viper.SetDefault("ENABLE_WORK_CONTINUE", false)
	viper.SetDefault("ENABLE_ECHO_PROGRESS_TO_FILES", false)
	viper.SetDefault("WORK_DURATION_MINUTES", 25)
	viper.SetDefault("REST_DURATION_MINUTES", 5)
	viper.SetDefault("LONG_REST_DURATION_MINUTES", 20)
	viper.SetDefault("MAX_CYCLES", 4)

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
