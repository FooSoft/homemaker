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

func processCmd(params []string, dir string, conf *config, flags int) error {
	args := appendExpEnv(nil, params)
	if len(args) == 0 {
		return fmt.Errorf("invalid command statement")
	}

	cmdName := args[0]
	if strings.HasPrefix(cmdName, "@") {
		macroName := strings.TrimPrefix(cmdName, "@")

		m, ok := conf.Macros[macroName]
		if !ok {
			return fmt.Errorf("macro not found: %s", macroName)
		}

		margs := appendExpEnv(nil, m.Prefix)
		if len(args) > 1 {
			margs = appendExpEnv(margs, args[1:])
		}
		margs = appendExpEnv(margs, m.Suffix)

		if flags&flagVerbose != 0 {
			log.Printf("using macro: %s", macroName)
		}

		return processCmd(margs, dir, conf, flags)
	}

	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if flags&flagVerbose != 0 {
		log.Printf("executing command: %s %s", cmdName, strings.Join(cmdArgs, " "))
	}

	return cmd.Run()
}
