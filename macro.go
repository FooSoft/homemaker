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

type macro struct {
	Prefix []string
	Suffix []string
}

func (m macro) process(dir string, params []string, flags int) error {
	var args []string
	args = appendExpEnv(args, m.Prefix)
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
