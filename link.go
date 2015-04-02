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

type link struct {
	Dst string
	Src string
}

func preparePath(loc string, flags int) error {
	clobber := flags&optClobber == optClobber
	force := flags&optForce == optForce
	verbose := flags&optVerbose == optVerbose

	if force {
		parentDir, _ := path.Split(loc)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			if verbose {
				log.Printf("Force creating path: '%s'", parentDir)
			}
			if err := os.MkdirAll(parentDir, 0777); err != nil {
				return err
			}
		}
	}

	if info, _ := os.Lstat(loc); info != nil {
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if verbose {
				log.Printf("Removing symlink: '%s'", loc)
			}
			if err := os.Remove(loc); err != nil {
				return err
			}
		} else if clobber {
			if verbose {
				log.Print("Clobbering path: '%s'", loc)
			}
			if err := os.RemoveAll(loc); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this link) install(srcDir, dstDir string, flags int) error {
	if len(this.Dst) == 0 {
		this.Dst = this.Src
	}

	srcPath := path.Join(srcDir, this.Src)
	dstPath := path.Join(dstDir, this.Dst)

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("Source path does not exist in filesystem: '%s'", srcPath)
	}

	if err := preparePath(dstPath, flags); err != nil {
		return err
	}

	if flags&optVerbose == optVerbose {
		log.Printf("Linking: '%s' => '%s'", srcPath, dstPath)
	}

	return os.Symlink(srcPath, dstPath)
}
