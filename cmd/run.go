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
	"log"
	"os"

	"github.com/FooSoft/homemaker/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the desired task",
	Args:  cobra.MaximumNArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		taskName := "default"
		if len(args) == 1 {
			taskName = args[0]
		}

		os.Setenv("HM_TASK", taskName)

		conf, err := internal.GenerateConfigStruct()
		if err != nil {
			panic(err)
		}

		log.Printf("Execute task %#v", taskName)

		internal.ProcessTask(taskName, conf)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().String("variant", "", "execution variant for tasks and macros")
	runCmd.Flags().Bool("clobber", false, "delete files and directories at target")
	runCmd.Flags().Bool("force", true, "create parent directories to target")
	runCmd.Flags().Bool("nocmds", false, "don't execute commands")
	runCmd.Flags().Bool("nolinks", false, "don't create links")
	runCmd.Flags().Bool("notemplates", false, "don't process templates")
	runCmd.Flags().Bool("unlink", false, "remove existing links instead of creating them")

	viper.BindPFlags(runCmd.Flags())
}
