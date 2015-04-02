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
	"flag"
	"fmt"
	"github.com/naoina/toml"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

const (
	optClobber = 1 << 0
	optForce   = 1 << 1
	optVerbose = 1 << 2
)

func parse(filename string) (*config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &config{}
	if err := toml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func printUsageAndExit() {
	_, executable := path.Split(os.Args[0])
	fmt.Errorf("Usage: %s [options] config_file [target_path]", executable)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	currUsr, err := user.Current()
	if err == nil {
		log.Fatal(err)
	}

	taskName := flag.String("task", "default", "name of task to execute")
	action := flag.String("action", "install", "'install' or 'uninstall' symlinks")
	dstDir := flag.String("target", currUsr.HomeDir, "target directory for symlinks")
	force := flag.Bool("force", true, "create parent directories to target")
	clobber := flag.Bool("clobber", false, "delete files and directories at target")
	verbose := flag.Bool("verbose", false, "verbose output")

	flag.Parse()

	flags := 0
	if *clobber {
		flags |= optClobber
	}
	if *force {
		flags |= optForce
	}
	if *verbose {
		flags |= optVerbose
	}

	if flag.NArg() == 0 {
		printUsageAndExit()
	}

	conf, err := parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	switch *action {
	case "install":
		if flag.NArg() >= 2 {
			if err := conf.install(flag.Arg(1), *dstDir, *taskName, flags); err != nil {
				log.Fatal(err)
			}
		} else {
			printUsageAndExit()
		}
	case "uninstall":
		if err := conf.uninstall(*dstDir, *taskName, flags); err != nil {
			log.Fatal(err)
		}
	default:
		printUsageAndExit()
	}
}
