/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
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

package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	homeDir := os.Getenv("HOME")
	app := cli.NewApp()

	app.Usage = "http://foosoft.net/projects/homemaker"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dest",
			Value: homeDir,
			Usage: "target directory for tasks",
		},
		cli.StringFlag{
			Name:  "task",
			Value: "default",
			Usage: "name of task to execute",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:  "variant",
			Value: "",
			Usage: "execution variant for tasks and macros",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "bootstrap",
			Aliases: []string{"b"},
			Usage:   "bootstrap a machine",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force",
					Usage: "create parent directories to target",
				},
				cli.BoolFlag{
					Name:  "clobber",
					Usage: "delete files and directories at target",
				},
				cli.BoolFlag{
					Name:  "nocmds",
					Usage: "don't execute commands",
				},
				cli.BoolFlag{
					Name:  "nolinks",
					Usage: "don't create links",
				},
			},
			Action: func(c *cli.Context) {
				checkNumberOfArgs(2, c)

				conf := makeConf(c)

				conf.srcDir = makeAbsPath(c.Args()[1])
				conf.dstDir = makeAbsPath(c.GlobalString("dest"))
				conf.task = c.GlobalString("task")
				conf.variant = c.GlobalString("variant")
				conf.force = c.Bool("force")
				conf.clobber = c.Bool("clobber")
				conf.verbose = c.GlobalBool("verbose")
				conf.nocmds = c.Bool("nocmds")
				conf.nolinks = c.Bool("nolinks")

				os.Setenv("HM_CONFIG", conf.file)
				os.Setenv("HM_TASK", c.GlobalString("task"))
				os.Setenv("HM_SRC", conf.srcDir)
				os.Setenv("HM_DEST", conf.dstDir)
				os.Setenv("HM_VARIANT", conf.variant)

				if err := processTask(c.GlobalString("task"), conf); err != nil {
					log.Println(err)
					cli.ShowAppHelp(c)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "encrypt task files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password",
					Value: "",
					Usage: "a password to encrypt a task",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "remove the original file after encryption",
				},
			},
			Action: func(c *cli.Context) {
				checkNumberOfArgs(1, c)

				conf := makeConf(c)

				conf.password = c.String("password")
				conf.remove = c.Bool("remove")

				if err := encryptTask(c.GlobalString("task"), conf); err != nil {
					log.Println(err)
					cli.ShowAppHelp(c)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "decrypt task files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password",
					Value: "",
					Usage: "a password to decrypt a task",
				},
			},
			Action: func(c *cli.Context) {
				checkNumberOfArgs(1, c)

				conf := makeConf(c)

				conf.password = c.String("password")

				if err := decryptTask(c.GlobalString("task"), conf); err != nil {
					log.Println(err)
					cli.ShowAppHelp(c)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "unlink",
			Aliases: []string{"u"},
			Usage:   "remove existing links instead of creating them",
			Action: func(c *cli.Context) {
				checkNumberOfArgs(1, c)

				conf := makeConf(c)

				conf.dstDir = makeAbsPath(c.GlobalString("dest"))
				conf.task = c.GlobalString("task")
				conf.variant = c.GlobalString("variant")
				conf.verbose = c.GlobalBool("verbose")

				if err := unlinkTask(c.GlobalString("task"), conf); err != nil {
					log.Println(err)
					cli.ShowAppHelp(c)
					os.Exit(1)
				}
			},
		},
	}

	app.Run(os.Args)
}

func checkNumberOfArgs(num int, c *cli.Context) {
	if len(c.Args()) != num {
		log.Println("Invalid number of arguments")
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
}

func makeConf(c *cli.Context) *config {
	confFile := makeAbsPath(c.Args()[0])

	conf, err := newConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}

	conf.file = confFile

	return conf
}
