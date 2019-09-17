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
	"strings"
)

func processEnv(env []string, conf *Config) error {
	args := appendExpEnv(nil, env)

	var value string
	switch {
	case len(args) == 0:
		return fmt.Errorf("invalid environment statement")
	case len(args) == 1:
		if conf.Verbose {
			log.Printf("unsetting variable: %s", args[0])
		}
		os.Unsetenv(args[0])
		return nil
	default:
		if strings.HasPrefix(args[1], "!") {
			var err error
			args[1] = strings.TrimLeft(args[1], "!")
			if value, err = processCmdWithReturn(args[1:], conf); err != nil {
				return err
			}
		} else {
			value = strings.Join(args[1:], ",")
		}
	}

	if conf.Verbose {
		log.Printf("setting variable %s to %s", args[0], value)
	}

	os.Setenv(args[0], value)
	return nil
}
