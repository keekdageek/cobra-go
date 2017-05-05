// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
	"github.com/gogap/logrus_mate"
	"reflect"
	"strings"
)

var cfgFile string
var logLevel string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hello",
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

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hello.yaml)")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log", "l", "", "Package Log Level")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//func run(cmd *cobra.Command, args []string) {
//	config, err := config.LoadConfig(cmd)
//	if err != nil {
//		log.Fatal("Failed to load config: " + err.Error())
//	}
//
//	//logger, err := conf.ConfigureLogging(&config.LogConfig)
//	//if err != nil {
//	//	log.Fatal("Failed to configure logging: " + err.Error())
//	//}
//	//
//	//logger.Infof("Starting with config: %+v", config)
//}

// initConfig reads in config file and ENV variables if set.
func initConfig() {


	// from the environment
	viper.SetEnvPrefix("HELLO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// from the config files
	customSettings := "hello.local"
	if cfgFile != "" {
		customSettings = cfgFile
	}
	files := []string{"hello", customSettings}
	configFiles := [] string{}

	for _, file := range files {
		viper.SetConfigName(file) // name of config file (without extension)
		viper.AddConfigPath(".")  // adding home directory as first search path

		// If a config file is found, read it in.
		if err := viper.MergeInConfig(); err == nil {
			configFiles = append(configFiles, viper.ConfigFileUsed())
		}
	}

	mate, _ := logrus_mate.NewLogrusMate(
		logrus_mate.ConfigFile(
			"./config/logrus.conf",
		),
	)

	// Initialize log, merge logrus and log flag
	if logLevel != "" {
		mate.Hijack(
			log.StandardLogger(), "hello",
			logrus_mate.ConfigString(fmt.Sprintf(`{ hello { level = "%s"} }`, logLevel)),
		)
	} else {
		mate.Hijack(
			log.StandardLogger(), "hello",
		)
	}

	for _, configFile := range configFiles {
		log.Debug("Using config file: ", configFile)
	}
}

func CallFuncByName(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("Method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	out = m.Call(in)
	return
}
