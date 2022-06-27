package cmd

import (
	"fmt"
	"net"
	"os"
	"parser/config"
	"parser/generated"
	"parser/parser_server"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var cfgFile string
var cnfg config.Config

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
		if cnfg.Localhost {
			os.Setenv("HTTPS_PROXY", "http://127.0.0.1:8888")
		}
		///////////////////////////////////////////////////////////
		// p := platform.New(cnfg)

		// cond := &platform.Condition{
		// 	Keyword:  "NAME",
		// 	Type:     "фильм",
		// 	Genres:   []string{"боевик"},
		// 	YearFrom: "2021",
		// 	YearTo:   "2023",
		// 	Coutries: []string{"США"},
		// }
		// fmt.Println(p.SearchByCondition(cond, cnfg.Proxy))
		///////////////////////////////////////////////////////////
		grpcServ := grpc.NewServer()
		pServ := parser_server.New(cnfg)
		generated.RegisterParserServiceServer(grpcServ, pServ)

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cnfg.ServerHost, cnfg.ServerPort))
		if err != nil {
			logrus.Fatalf("failed to listen: %v", err)
		}

		if cnfg.WithReflection {
			reflection.Register(grpcServ)
		}

		logrus.Info("Starting gRPC listener on port " + cnfg.ServerPort)
		if err := grpcServ.Serve(lis); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
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
