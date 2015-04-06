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

func createPath(loc string, flags int) error {
	if flags&flagForce == flagForce {
		parentDir, _ := path.Split(loc)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			if flags&flagVerbose == flagVerbose {
				log.Printf("force creating path %s", parentDir)
			}
			if err := os.MkdirAll(parentDir, 0777); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *link) destination() string {
	if len(*this) > 0 {
		return (*this)[0]
	}

	return ""
}

func (this *link) source() string {
	if len(*this) > 1 {
		return (*this)[1]
	}

	return this.destination()
}

func (this *link) valid() bool {
	length := len(*this)
	return length >= 1 && length <= 2
}

func (this *link) process(srcDir, dstDir string, flags int) error {
	if !this.valid() {
		return fmt.Errorf("link element is invalid")
	}

	srcPath := path.Join(srcDir, this.source())
	dstPath := path.Join(dstDir, this.destination())

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source path %s does not exist in filesystem", srcPath)
	}

	if err := createPath(dstPath, flags); err != nil {
		return err
	}

	if err := cleanPath(dstPath, flags); err != nil {
		return err
	}

	if flags&flagVerbose == flagVerbose {
		log.Printf("linking %s to %s", srcPath, dstPath)
	}

	return os.Symlink(srcPath, dstPath)
}
