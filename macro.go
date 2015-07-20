/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
all
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
	"os/exec"
	"strings"
)

type macroDef struct {
	Prefix []string
	Suffix []string
}

func (m macroDef) process(dir string, params []string, flags int) error {
	args := appendExpEnv(nil, m.Prefix)
	args = appendExpEnv(args, params)
	args = appendExpEnv(args, m.Suffix)

	var cmd *exec.Cmd
	switch {
	case len(args) == 0:
		return fmt.Errorf("macro element is invalid")
	case len(args) == 1:
		cmd = exec.Command(args[0])
	default:
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if flags&flagVerbose == flagVerbose {
		log.Printf("executing macro %s", strings.Join(args, " "))
	}

	return cmd.Run()
}

func processMacro(args []string, dir string, conf *config, flags int) error {
	if len(args) == 0 {
		return fmt.Errorf("macro element is invalid")
	}

	macro, ok := conf.Macros[args[0]]
	if !ok {
		return fmt.Errorf("macro dependency not found %s", args[0])
	}

	return macro.process(dir, args[1:], flags)
}
