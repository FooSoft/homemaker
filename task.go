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
	"fmt"
	"log"
)

type task struct {
	Deps   []string
	Links  [][]string
	Cmds   [][]string
	Envs   [][]string
	WTrue  [][]string
	WFalse [][]string
}

func (t *task) deps(conf *config) []string {
	deps := t.Deps

	if conf.flags&flagNoCmds == 0 {
		for _, currCmd := range t.Cmds {
			deps = append(deps, findCmdDeps(currCmd, conf)...)
		}
	}

	return deps
}

func (t *task) process(conf *config) error {
	for _, currTask := range t.deps(conf) {
		if err := processTask(currTask, conf); err != nil {
			return err
		}
	}

	for _, currEnv := range t.Envs {
		if err := processEnv(currEnv, conf); err != nil {
			return err
		}
	}

	if conf.flags&flagNoCmds == 0 {
	CmdLoop:
		for {
			for _, currCnd := range t.WTrue {
				if err := processCmd(currCnd, false, conf); err != nil {
					break CmdLoop
				}
			}

			for _, currCnd := range t.WFalse {
				if err := processCmd(currCnd, false, conf); err == nil {
					break CmdLoop
				}
			}

			for _, currCmd := range t.Cmds {
				if err := processCmd(currCmd, true, conf); err != nil {
					return err
				}
			}

			break
		}
	}

	if conf.flags&flagNoLinks == 0 {
		for _, currLink := range t.Links {
			if err := processLink(currLink, conf); err != nil {
				return err
			}
		}
	}

	return nil
}

func processTask(taskName string, conf *config) error {
	for _, tn := range makeVariantNames(taskName, conf.variant) {
		t, ok := conf.Tasks[tn]
		if !ok {
			continue
		}

		if conf.handled[tn] {
			if conf.flags&flagVerbose != 0 {
				log.Printf("skipping processed task: %s", tn)
			}

			return nil
		}

		if conf.flags&flagVerbose != 0 {
			log.Printf("processing task: %s", tn)
		}

		conf.handled[tn] = true
		return t.process(conf)
	}

	return fmt.Errorf("task or variant not found: %s", taskName)
}
