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

package internal

import (
	"fmt"
	"log"
	"os"
)

type task struct {
	Deps      []string
	Links     [][]string
	CmdsPre   [][]string
	Cmds      [][]string
	CmdsPost  [][]string
	Envs      [][]string
	Accepts   [][]string
	Rejects   [][]string
	Templates [][]string
}

func (t *task) deps(conf *Config) []string {
	deps := t.Deps

	if !conf.Nocmds {
		for _, currCmd := range t.Cmds {
			deps = append(deps, findCmdDeps(currCmd, conf)...)
		}
	}

	return deps
}

func (t *task) process(conf *Config) error {
	for _, currTask := range t.deps(conf) {
		currTask = os.ExpandEnv(currTask)
		if err := ProcessTask(currTask, conf); err != nil {
			return err
		}
	}

	for _, currEnv := range t.Envs {
		if err := processEnv(currEnv, conf); err != nil {
			return err
		}
	}

	if !conf.Nocmds {
		for _, currCmd := range t.CmdsPre {
			if err := processCmd(currCmd, true, conf); err != nil {
				return err
			}
		}

		for _, currCmd := range t.Cmds {
			if err := processCmd(currCmd, true, conf); err != nil {
				return err
			}
		}
	}

	if !conf.Nolinks {
		for _, currLink := range t.Links {
			if err := processLink(currLink, conf); err != nil {
				return err
			}
		}
	}

	if !conf.Notemplates {
		for _, currTmpl := range t.Templates {
			if err := processTemplate(currTmpl, conf); err != nil {
				return err
			}
		}
	}

	if !conf.Nocmds {
		for _, currCmd := range t.CmdsPost {
			if err := processCmd(currCmd, true, conf); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *task) skippable(conf *Config) bool {
	for _, currCnd := range t.Accepts {
		if err := processCmd(currCnd, false, conf); err != nil {
			return true
		}
	}

	for _, currCnd := range t.Rejects {
		if err := processCmd(currCnd, false, conf); err == nil {
			return true
		}
	}

	return false
}

func ProcessTask(taskName string, conf *Config) error {
	for _, tn := range makeVariantNames(taskName, conf.Variant) {
		if conf.Verbose {
			log.Printf("starting task: %s", tn)
		}
		t, ok := conf.Tasks[tn]
		if !ok {
			continue
		}

		if conf.handled[tn] || t.skippable(conf) {
			if conf.Verbose {
				log.Printf("skipping task: %s", tn)
			}

			return nil
		}

		if conf.Verbose {
			log.Printf("processing task: %s", tn)
		}

		conf.handled[tn] = true
		return t.process(conf)
	}

	return fmt.Errorf("task or variant not found: %s", taskName)
}
