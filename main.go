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
	"github.com/naoina/toml"
	"io/ioutil"
	"log"
	"os/user"
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

func install(conf config, name, target, source string, force, clobber bool) error {
	return nil
}

func uninstall(conf config, name, target string, force, clobber bool) error {
	return nil
}

func main() {
	currUsr, err := user.Current()
	if err == nil {
		log.Fatal(err)
	}

	action := flag.String("action", "install", "'install' or 'uninstall' symlinks")
	clobber := flag.Bool("clobber", false, "delete files and directories at target")
	force := flag.Bool("force", true, "force creation of parent directories for target")
	profile := flag.String("profile", "default", "name of profile to execute")
	target := flag.String("target", currUsr.HomeDir, "target directory for symlinks")

	flag.Parse()

	confPath := flag.Arg(0)
	source := flag.Arg(1)

	conf, err := parse(confPath)
	if err != nil {
		log.Fatal(err)
	}

	switch *action {
	case "install":
		install(*conf, *profile, *target, source, *force, *clobber)
	case "uninstall":
		uninstall(*conf, *profile, *target, *force, *clobber)
	default:
		log.Fatalf("Unrecognized action: '%s'", action)
	}
}
