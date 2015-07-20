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
	"strings"
)

type env []string

func (e env) process(flags int) error {
	items := appendExpEnv(nil, e)

	var value string
	switch {
	case len(items) == 0:
		return fmt.Errorf("enviornment element is invalid")
	case len(items) == 1:
		if flags&flagVerbose == flagVerbose {
			log.Printf("unsetting variable %s", items[0])
		}
		os.Unsetenv(items[0])
		return nil
	case len(items) == 2:
		value = items[1]
	default:
		value = strings.Join(items[1:], ",")
	}

	if flags&flagVerbose == flagVerbose {
		log.Printf("setting variable %s to %s", items[0], value)
	}

	os.Setenv(items[0], value)
	return nil
}
