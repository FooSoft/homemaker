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
	"github.com/naoina/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	flagClobber = 1 << iota
	flagForce
	flagVerbose
)

func parse(filename string) (*config, error) {
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
	case ".toml":
		if err := toml.Unmarshal(bytes, &conf); err != nil {
			return nil, err
		}
	case ".yaml":
		if err := yaml.Unmarshal(bytes, &conf); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported configuration file format")
	}

	return conf, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] conf [src]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "http://foosoft.net/projects/homemaker/\n\n")
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func makeAbsPath(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	return path
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

	flag.Usage = usage
	flag.Parse()

	flags := 0
	if *clobber {
		flags |= flagClobber
	}
	if *force {
		flags |= flagForce
	}
	if *verbose {
		flags |= flagVerbose
	}

	if flag.NArg() == 2 {
		conf, err := parse(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		if err := conf.process(makeAbsPath(flag.Arg(1)), makeAbsPath(*dstDir), *taskName, flags); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
		os.Exit(2)
	}
}
