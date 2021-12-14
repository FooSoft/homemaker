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
	"path"
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

func cleanPath(loc string, flags int) (bool, error) {
	if info, _ := os.Lstat(loc); info != nil {
		if info.Mode()&os.ModeSymlink == 0 {
			shouldContinue := false
			if flags&flagClobber == 0 {
				shouldContinue = prompt("clobber path", loc)
			}
			if flags&flagClobber != 0 || shouldContinue {
				if flags&flagVerbose != 0 {
					log.Printf("clobbering path: %s", loc)
				}
				if err := try(func() error { return os.RemoveAll(loc) }) ; err != nil {
					return false, err
				}
			} else {
				return false, nil
			}
		} else {
			if flags&flagVerbose != 0 {
				log.Printf("removing symlink: %s", loc)
			}
			if err := try(func() error { return os.Remove(loc) }); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func createPath(loc string, flags int, mode os.FileMode) error {
	parentDir := path.Dir(loc)

	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if flags&flagForce != 0 || prompt("force create path", parentDir) {
			if flags&flagVerbose != 0 {
				log.Printf("force creating path: %s", parentDir)
			}
			if err := os.MkdirAll(parentDir, mode); err != nil {
				return err
			}
		}
	}

	return nil
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
