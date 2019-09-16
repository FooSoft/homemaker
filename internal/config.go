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
	"os"

	"github.com/spf13/viper"
)

const (
	flagClobber = 1 << iota
	flagForce
	flagVerbose
	flagNoCmds
	flagNoLinks
	flagNoTemplates
	flagNoMacro
	flagUnlink = flagNoCmds | (1 << iota)
)

type Config struct {
	Tasks   map[string]task
	Macros  map[string]macro
	File    string
	SrcDir  string `mapstructure:"home-src"`
	DstDir  string `mapstructure:"home-dst"`
	Variant string

	Clobber     bool
	Force       bool
	Verbose     bool
	Nocmds      bool
	Nolinks     bool
	Notemplates bool
	Unlink      bool

	handled map[string]bool
	flags   int
}

func GenerateConfigStruct() (*Config, error) {
	conf := &Config{handled: make(map[string]bool)}
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	conf.digest()

	return conf, nil
}

func (c *Config) digest() {
	c.SrcDir = makeAbsPath(c.SrcDir)
	c.DstDir = makeAbsPath(c.DstDir)

	flags := 0
	if c.Clobber {
		flags |= flagClobber
	}
	if c.Force {
		flags |= flagForce
	}
	if c.Verbose {
		flags |= flagVerbose
	}
	if c.Nocmds {
		flags |= flagNoCmds
	}
	if c.Nolinks {
		flags |= flagNoLinks
	}
	if c.Notemplates {
		flags |= flagNoTemplates
	}
	if c.Unlink {
		flags |= flagUnlink
	}
	c.flags = flags
}

func (c *Config) SetEnv() {
	os.Setenv("HM_CONFIG", c.File)
	os.Setenv("HM_SRC", c.SrcDir)
	os.Setenv("HM_DEST", c.DstDir)
	os.Setenv("HM_VARIANT", c.Variant)
}
