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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

const (
	FLAG_CLOBBER = 1 << iota
	FLAG_FORCE
	FLAG_VERBOSE
	FLAG_NO_CMD
	FLAG_NO_LINK
	FLAG_NO_MACRO
)

func parseCfg(filename string) (*config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &config{}
	switch path.Ext(filename) {
	case ".json":
		if err := json.Unmarshal(bytes, &conf); err != nil {
			return nil, err
		}
	case ".toml", ".tml":
		if err := toml.Unmarshal(bytes, &conf); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(bytes, &conf); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported configuration file format")
	}

	return conf, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] conf src\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "http://foosoft.net/projects/homemaker/\n\n")
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	currUsr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	taskName := flag.String("task", "default", "name of task to execute")
	dstDir := flag.String("dest", currUsr.HomeDir, "target directory for tasks")
	force := flag.Bool("force", true, "create parent directories to target")
	clobber := flag.Bool("clobber", false, "delete files and directories at target")
	verbose := flag.Bool("verbose", false, "verbose output")
	nocmd := flag.Bool("nocmd", false, "don't execute commands")
	nolink := flag.Bool("nolink", false, "don't create links")

	flag.Usage = usage
	flag.Parse()

	flags := 0
	if *clobber {
		flags |= FLAG_CLOBBER
	}
	if *force {
		flags |= FLAG_FORCE
	}
	if *verbose {
		flags |= FLAG_VERBOSE
	}
	if *nocmd {
		flags |= FLAG_NO_CMD
	}
	if *nolink {
		flags |= FLAG_NO_LINK
	}

	if flag.NArg() == 2 {
		confDirAbs := makeAbsPath(flag.Arg(0))
		srcDirAbs := makeAbsPath(flag.Arg(1))
		dstDirAbs := makeAbsPath(*dstDir)

		os.Setenv("HM_CONFIG", confDirAbs)
		os.Setenv("HM_TASK", *taskName)
		os.Setenv("HM_SRC", srcDirAbs)
		os.Setenv("HM_DEST", dstDirAbs)

		conf, err := parseCfg(confDirAbs)
		if err != nil {
			log.Fatal(err)
		}

		if err := conf.process(srcDirAbs, dstDirAbs, *taskName, flags); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
		os.Exit(2)
	}
}
