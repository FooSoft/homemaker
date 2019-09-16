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
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type macro struct {
	Deps   []string
	Prefix []string
	Suffix []string
}

func findCmdMacro(macroName string, conf *Config) (*macro, string) {
	if strings.HasPrefix(macroName, "@") {
		mn := strings.TrimPrefix(macroName, "@")
		for _, mn := range makeVariantNames(mn, conf.Variant) {
			if m, ok := conf.Macros[mn]; ok {
				return &m, mn
			}
		}
	}

	return nil, ""
}

func findCmdDeps(params []string, conf *Config) []string {
	if len(params) == 0 {
		return nil
	}

	if m, _ := findCmdMacro(params[0], conf); m != nil {
		return m.Deps
	}

	return nil
}

func processCmdMacro(macroName string, args []string, interact bool, conf *Config) error {
	m, mn := findCmdMacro(macroName, conf)
	if m == nil {
		return fmt.Errorf("macro or variant not found: %s", macroName)
	}

	margs := appendExpEnv(nil, m.Prefix)
	margs = appendExpEnv(margs, args)
	margs = appendExpEnv(margs, m.Suffix)

	if conf.flags&flagVerbose != 0 {
		log.Printf("expanding macro: %s", mn)
	}

	return processCmd(margs, interact, conf)
}

func processCmd(params []string, interact bool, conf *Config) error {
	args := appendExpEnv(nil, params)
	if len(args) == 0 {
		return fmt.Errorf("invalid command statement")
	}

	cmdName := args[0]
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	if strings.HasPrefix(cmdName, "@") {
		return processCmdMacro(cmdName, cmdArgs, interact, conf)
	}

	if conf.flags&flagVerbose != 0 {
		log.Printf("executing command: %s %s", cmdName, strings.Join(cmdArgs, " "))
	}

	exec := func() error {
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Dir = conf.DstDir
		if interact {
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
		}

		return cmd.Run()
	}

	if interact {
		return try(exec)
	}

	return exec()
}

func processCmdWithReturn(params []string, conf *Config) (string, error) {
	args := appendExpEnv(nil, params)
	if len(args) == 0 {
		return "", fmt.Errorf("invalid command statement")
	}

	cmdName := args[0]
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	if strings.HasPrefix(cmdName, "@") {
		return "", processCmdMacro(cmdName, cmdArgs, false, conf)
	}

	if conf.flags&flagVerbose != 0 {
		log.Printf("executing command (with return): %s %s", cmdName, strings.Join(cmdArgs, " "))
	}

	exec := func() (string, error) {
		var stdout bytes.Buffer
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Dir = conf.DstDir
		cmd.Stderr = os.Stderr
		cmd.Stdout = &stdout
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		return strings.Trim(stdout.String(), "\r\n"), err
	}

	return exec()
}
