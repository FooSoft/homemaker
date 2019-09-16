/*
 * Copyright (c) 2018 Metalblueberry <metalblueberry@gmail.com>
 * Author: Metalblueberry <metalblueberry@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "homemaker",
	Short: "A lightweight tool for dot-files managment",
	Long: `Homemaker is a lightweight tool for straightforward and efficient 
management of *nix configuration files found in the user's home directory,
commonly known as dot-files. It can also be readily used for general purpose
system bootstrapping, including installing packages, cloning repositories,
etc. This tool is written in Go, requires no installation, has no dependencies 
and makes use of simple configuration file structure inspired by make to
generate symlinks and execute system commands to aid in configuring a new
system for use.

Full docs at https://foosoft.net/projects/homemaker/
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "conf", "", "config file (default is homemaker.yaml)")

	rootCmd.PersistentFlags().String("home-dst", os.Getenv("HOME"), "target directory for tasks")
	rootCmd.PersistentFlags().String("home-src", ".", "source directory for tasks")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")

	viper.BindPFlags(rootCmd.PersistentFlags())

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("homemaker")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
