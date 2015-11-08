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
)

type task struct {
	Deps    []string
	Links   [][]string
	Cmds    [][]string
	Envs    [][]string
	Secrets []string
}

func (t *task) deps(conf *config) []string {
	deps := t.Deps

	if conf.nocmds {
		for _, currCmd := range t.Cmds {
			deps = append(deps, findCmdDeps(currCmd, conf)...)
		}
	}

	return deps
}

func (t *task) process(conf *config, key string) error {
	for _, currTask := range t.deps(conf) {
		if err := processTask(currTask, conf); err != nil {
			return err
		}
	}

	for _, currSecret := range t.Secrets {
		if err := processEnc(currSecret, conf, key); err != nil {
			return err
		}
	}

	for _, currEnv := range t.Envs {
		if err := processEnv(currEnv, conf); err != nil {
			return err
		}
	}

	if !conf.nocmds {
		for _, currCmd := range t.Cmds {
			if err := processCmd(currCmd, conf); err != nil {
				return err
			}
		}
	}

	if !conf.nolinks {
		for _, currLink := range t.Links {
			if err := processLink(currLink, conf); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *task) hashPassword(task string, conf *config) (string, error) {
	var key string

	if len(t.Secrets) != 0 {
		if conf.password != "" {
			key = keyFromPasssword([]byte(conf.password))
		} else {
			fmt.Printf("Enter your password for %s: \n", task)
			k, err := getMaskedInput()
			if err != nil {
				return "", err
			}
			key = keyFromPasssword(k)
		}
	}

	return key, nil
}

func processTask(taskName string, conf *config) error {
	for _, tn := range makeVariantNames(taskName, conf.variant) {
		t, ok := conf.Tasks[tn]
		if !ok {
			continue
		}

		if conf.handled[tn] {
			if conf.verbose {
				log.Printf("skipping processed task: %s", tn)
			}

			return nil
		}

		if conf.verbose {
			log.Printf("processing task: %s", tn)
		}

		key, err := t.hashPassword(tn, conf)
		if err != nil {
			return err
		}

		conf.handled[tn] = true
		return t.process(conf, key)
	}

	return fmt.Errorf("task or variant not found: %s", taskName)
}

func (t *task) encrypt(conf *config, key string) error {
	for _, currTask := range t.deps(conf) {
		if err := encryptTask(currTask, conf); err != nil {
			return err
		}
	}

	for _, currSecret := range t.Secrets {
		if err := encryptSecret(currSecret, conf, key); err != nil {
			return err
		}
	}

	return nil
}

func (t *task) decrypt(conf *config, key string) error {
	for _, currTask := range t.deps(conf) {
		if err := decryptTask(currTask, conf); err != nil {
			return err
		}
	}

	for _, currSecret := range t.Secrets {
		if err := decryptSecret(currSecret, conf, key); err != nil {
			return err
		}
	}

	return nil
}

func encryptTask(taskName string, conf *config) error {
	for _, tn := range makeVariantNames(taskName, conf.variant) {
		t, ok := conf.Tasks[tn]
		if !ok {
			continue
		}

		if conf.handled[tn] {
			if conf.verbose {
				log.Printf("skipping processed task: %s", tn)
			}

			return nil
		}

		if conf.verbose {
			log.Printf("encrypting task: %s", tn)
		}

		key, err := t.hashPassword(tn, conf)
		if err != nil {
			return err
		}

		conf.handled[tn] = true
		return t.encrypt(conf, key)
	}

	return fmt.Errorf("task or variant not found: %s", taskName)
}

func decryptTask(taskName string, conf *config) error {
	for _, tn := range makeVariantNames(taskName, conf.variant) {
		t, ok := conf.Tasks[tn]
		if !ok {
			continue
		}

		if conf.handled[tn] {
			if conf.verbose {
				log.Printf("skipping processed task: %s", tn)
			}

			return nil
		}

		if conf.verbose {
			log.Printf("encrypting task: %s", tn)
		}

		key, err := t.hashPassword(tn, conf)
		if err != nil {
			return err
		}

		conf.handled[tn] = true
		return t.decrypt(conf, key)
	}

	return fmt.Errorf("task or variant not found: %s", taskName)
}

func (t *task) unlink(conf *config) error {
	for _, currLink := range t.deps(conf) {
		if err := unlinkTask(currLink, conf); err != nil {
			return err
		}
	}

	for _, currLink := range t.Links {
		if err := removeLink(currLink, conf); err != nil {
			return err
		}
	}

	return nil
}

func unlinkTask(taskName string, conf *config) error {
	for _, tn := range makeVariantNames(taskName, conf.variant) {
		t, ok := conf.Tasks[tn]
		if !ok {
			continue
		}

		if conf.handled[tn] {
			if conf.verbose {
				log.Printf("skipping processed task: %s", tn)
			}

			return nil
		}

		if conf.verbose {
			log.Printf("unlinking task: %s", tn)
		}

		return t.unlink(conf)
	}

	return fmt.Errorf("task or variant not found: %s", taskName)
}
