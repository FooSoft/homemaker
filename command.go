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
	"os"
	"os/exec"
	"strings"
)

type macro struct {
	Prefix []string
	Suffix []string
}

func processCmdMacro(macroName string, args []string, conf *config) error {
	for _, mn := range makeVariantNames(macroName, conf.variant) {
		m, ok := conf.Macros[mn]
		if !ok {
			continue
		}

		margs := appendExpEnv(nil, m.Prefix)
		margs = appendExpEnv(margs, args)
		margs = appendExpEnv(margs, m.Suffix)

		if conf.flags&flagVerbose != 0 {
			log.Printf("using macro: %s", mn)
		}

		return processCmd(margs, conf)
	}

	return fmt.Errorf("macro or variant not found: %s", macroName)
}

func processCmd(params []string, conf *config) error {
	args := appendExpEnv(nil, params)
	if len(args) == 0 {
		return fmt.Errorf("invalid command statement")
	}

	cmdName := args[0]
	if strings.HasPrefix(cmdName, "@") {
		var (
			macroName = strings.TrimPrefix(cmdName, "@")
			macroArgs []string
		)

		if len(args) > 1 {
			macroArgs = args[1:]
		}

		return processCmdMacro(macroName, macroArgs, conf)
	}

	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	if conf.flags&flagVerbose != 0 {
		log.Printf("executing command: %s %s", cmdName, strings.Join(cmdArgs, " "))
	}

	return try(func() error {
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Dir = conf.dstDir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()
	})
}
