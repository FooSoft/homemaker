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
	"path/filepath"
	"strconv"
)

func parseLink(params []string) (srcPath, dstPath string, mode os.FileMode, err error) {
	length := len(params)
	if length < 1 || length > 3 {
		err = fmt.Errorf("invalid link statement")
		return
	}

	if length > 2 {
		var parsed uint64
		parsed, err = strconv.ParseUint(params[2], 0, 64)
		if err != nil {
			return
		}

		mode = os.FileMode(parsed)
	} else {
		mode = 0755
	}

	dstPath = os.ExpandEnv(params[0])
	srcPath = dstPath
	if length > 1 {
		srcPath = os.ExpandEnv(params[1])
	}

	return
}

func processLink(params []string, conf *config) error {
	srcPath, dstPath, mode, err := parseLink(params)
	if err != nil {
		return err
	}

	srcPathAbs := srcPath
	if !filepath.IsAbs(srcPathAbs) {
		srcPathAbs = filepath.Join(conf.srcDir, srcPath)
	}

	dstPathAbs := dstPath
	if !filepath.IsAbs(dstPathAbs) {
		dstPathAbs = filepath.Join(conf.dstDir, dstPath)
	}

	if conf.flags&flagUnlink != flagUnlink {
		if _, err := os.Stat(srcPathAbs); os.IsNotExist(err) {
			return fmt.Errorf("source path %s does not exist in filesystem", srcPathAbs)
		}

		if err := try(func() error { return createPath(dstPathAbs, conf.flags, mode) }); err != nil {
			return err
		}

		pathCleaned, err := cleanPath(dstPathAbs, conf.flags)
		if err != nil {
			return err
		}
		if !pathCleaned {
			return nil
		}

		if conf.flags&flagVerbose != 0 {
			log.Printf("linking %s to %s", srcPathAbs, dstPathAbs)
		}

		return try(func() error {
			return os.Symlink(srcPathAbs, dstPathAbs)
		})
	} else {
		stat, err := os.Lstat(dstPathAbs)
		if os.IsNotExist(err) || stat.Mode()&os.ModeSymlink == 0 {
			return nil
		}

		_, err = cleanPath(dstPathAbs, conf.flags)
		return err
	}
}
