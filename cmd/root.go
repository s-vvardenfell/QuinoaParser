package cmd

import (
	"fmt"
	"os"
	"parser/kinoafisha"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cnfg Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quinoa",
	Short: "Quinoa is a telegram bot server and parser",
	Long: `Usage: quinoa --config <config_file_name.yaml>
	Parsed platforms:
	-kp
	-im
	-zn

	API methods:
	-1
	-2
	-3
	`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Parser started...")
		// fmt.Println(cnfg)
		// for p := cnfg.Platforms{

		// }

		ka := kinoafisha.New(cnfg.Proxy)
		res := ka.ParseSeriesCalendar()

		fmt.Println(res)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatalf("execute fails: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", "", "config file (default is resources/config.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(wd, "resources"))
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}

	if err := viper.Unmarshal(&cnfg); err != nil {
		cobra.CheckErr(err)
	}
}