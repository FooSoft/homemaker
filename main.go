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
		return nil, fmt.Errorf("Unsupported configuration file format")
	}

	return conf, nil
}

func fatalUsage() {
	_, executable := path.Split(os.Args[0])
	fmt.Printf("Usage: %s [options] config_file [source_dir]\n", executable)
	flag.PrintDefaults()
	os.Exit(1)
}

func absPath(path string) string {
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
	action := flag.String("action", "install", "'install' or 'uninstall' symlinks")
	dstDir := flag.String("dest", currUsr.HomeDir, "target directory for symlinks")
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
		fatalUsage()
	}

	conf, err := parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	switch *action {
	case "install":
		if flag.NArg() >= 2 {
			if err := conf.install(absPath(flag.Arg(1)), absPath(*dstDir), *taskName, flags); err != nil {
				log.Fatal(err)
			}
		} else {
			fatalUsage()
		}
	case "uninstall":
		if err := conf.uninstall(absPath(*dstDir), *taskName, flags); err != nil {
			log.Fatal(err)
		}
	default:
		fatalUsage()
	}
}
