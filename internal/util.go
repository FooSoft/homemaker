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
	"path/filepath"
	"strings"
)

func appendExpEnv(dst, src []string) []string {
	for _, value := range src {
		dst = append(dst, os.ExpandEnv(value))
	}

	return dst
}

func makeAbsPath(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	return path
}

func makeVariantNames(name, variant string) []string {
	if nameParts := strings.Split(name, "__"); len(nameParts) > 1 {
		variant = nameParts[len(nameParts)-1]
		name = strings.Join(nameParts[:len(nameParts)-1], "")
	}

	names := []string{name}
	if len(variant) > 0 && !strings.HasSuffix(name, "__") {
		names = []string{fmt.Sprint(name, "__", variant), name}
	}

	return names
}

func prompt(prompts ...string) bool {
	for {
		fmt.Printf("%s: [y]es, [n]o? ", strings.Join(prompts, " "))

		var ans string
		fmt.Scanln(&ans)

		switch strings.ToLower(ans) {
		case "y":
			return true
		case "n":
			return false
		}
	}
}

func try(task func() error) error {
	for {
		err := task()
		if err == nil {
			return nil
		}

	loop:
		for {
			fmt.Printf("%s: [a]bort, [r]etry, [c]ancel? ", err)

			var ans string
			fmt.Scanln(&ans)

			switch strings.ToLower(ans) {
			case "a":
				return err
			case "r":
				break loop
			case "c":
				return nil
			}
		}
	}
}
