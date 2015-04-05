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

func cleanPath(loc string, flags int) error {
	verbose := flags&optVerbose == optVerbose

	if info, _ := os.Lstat(loc); info != nil {
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if verbose {
				log.Printf("Removing symlink: '%s'", loc)
			}
			if err := os.Remove(loc); err != nil {
				return err
			}
		} else {
			if flags&optClobber == optClobber {
				if verbose {
					log.Print("Clobbering path: '%s'", loc)
				}
				if err := os.RemoveAll(loc); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Cannot create link; target already exists: '%s'", loc)
			}
		}
	}

	return nil
}

func createPath(loc string, flags int) error {
	if flags&optForce == 0 {
		return nil
	}

	parentDir, _ := path.Split(loc)

	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if flags&optVerbose == optVerbose {
			log.Printf("Force creating path: '%s'", parentDir)
		}
		if err := os.MkdirAll(parentDir, 0777); err != nil {
			return err
		}
	}

	return nil
}

func (this *link) source() string {
	if len(*this) > 0 {
		return (*this)[0]
	}

	return ""
}

func (this *link) destination() string {
	if len(*this) > 1 {
		return (*this)[1]
	}

	return this.source()
}

func (this *link) valid() bool {
	length := len(*this)
	return length >= 1 && length <= 2
}

func (this *link) install(srcDir, dstDir string, flags int) error {
	if !this.valid() {
		return fmt.Errorf("Link element is invalid")
	}

	srcPath := path.Join(srcDir, this.source())
	dstPath := path.Join(dstDir, this.destination())

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("Source path does not exist in filesystem: '%s'", srcPath)
	}

	if err := createPath(dstPath, flags); err != nil {
		return err
	}

	if err := cleanPath(dstPath, flags); err != nil {
		return err
	}

	if flags&optVerbose == optVerbose {
		log.Printf("Linking: '%s' to '%s'", srcPath, dstPath)
	}

	return os.Symlink(srcPath, dstPath)
}
