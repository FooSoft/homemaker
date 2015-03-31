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
	"github.com/naoina/toml"
	"io/ioutil"
	"log"
)

type link struct {
	Dst   string
	Force bool
	Src   string
	Stomp bool
}

type profile struct {
	Deps  []string
	Force bool
	Links []link
	Stomp bool
}

type config struct {
	Force bool
	Profs map[string]profile
	Stomp bool
}

func parse(filename string) (*config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &config{}
	if err := toml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (conf *config) prepare() {

}

func process(src, dst string) error {
	return nil
}

func main() {
	conf, err := parse("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	conf.prepare()

	if err := process("/mnt/storage/sync/Dropbox", "/mnt/storage/projects/blah"); err != nil {
		log.Fatal(err)
	}

	log.Print(conf)
}
