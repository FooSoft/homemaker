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
	"strconv"
)

type link []string

func cleanPath(loc string, flags int) error {
	if info, _ := os.Lstat(loc); info != nil {
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if flags&flagVerbose == flagVerbose {
				log.Printf("removing symlink %s", loc)
			}
			if err := os.Remove(loc); err != nil {
				return err
			}
		} else {
			if flags&flagClobber == flagClobber {
				if flags&flagVerbose == flagVerbose {
					log.Printf("clobbering path %s", loc)
				}
				if err := os.RemoveAll(loc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func createPath(loc string, flags int, mode os.FileMode) error {
	if flags&flagForce == flagForce {
		parentDir, _ := path.Split(loc)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			if flags&flagVerbose == flagVerbose {
				log.Printf("force creating path %s", parentDir)
			}
			if err := os.MkdirAll(parentDir, mode); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *link) parse() (string, string, os.FileMode, error) {
	length := len(*l)
	if length < 1 || length > 3 {
		return "", "", 0, fmt.Errorf("link element is invalid")
	}

	dstPath := (*l)[0]
	srcPath := dstPath
	if length > 1 {
		srcPath = (*l)[1]
	}

	var mode os.FileMode = 0755
	if length > 2 {
		parsed, err := strconv.ParseUint((*l)[2], 0, 64)
		if err != nil {
			return "", "", 0, err
		}

		mode = os.FileMode(parsed)
	}

	return srcPath, dstPath, mode, nil
}

func (l *link) process(srcDir, dstDir string, flags int) error {
	srcPath, dstPath, mode, err := l.parse()
	if err != nil {
		return err
	}

	srcPathAbs := srcPath
	if !path.IsAbs(srcPathAbs) {
		srcPathAbs = path.Join(srcDir, srcPath)
	}

	dstPathAbs := dstPath
	if !path.IsAbs(dstPathAbs) {
		dstPathAbs = path.Join(dstDir, dstPath)
	}

	if _, err := os.Stat(srcPathAbs); os.IsNotExist(err) {
		return fmt.Errorf("source path %s does not exist in filesystem", srcPathAbs)
	}

	if err := createPath(dstPathAbs, flags, mode); err != nil {
		return err
	}

	if err := cleanPath(dstPathAbs, flags); err != nil {
		return err
	}

	if flags&flagVerbose == flagVerbose {
		log.Printf("linking %s to %s", srcPathAbs, dstPathAbs)
	}

	if err := os.Symlink(srcPathAbs, dstPathAbs); err != nil {
		return err
	}

	return nil
}
