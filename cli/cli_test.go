/*
 * MIT License
 *
 * Copyright (c) 2020 TCorp BV
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cli

import (
	"io/ioutil"
	"testing"
)

func TestCommandRouting(t *testing.T) {
	app := NewApp()
	app.Writer = ioutil.Discard // Discard output

	firstCalled := 0
	firstArgsMatchExpected := false
	secondCalled := 0
	secondArgsMatchExpected := false

	//Setup some command handlers to test
	app.Commands = []*Command{
		{Name: "mycommand", Usage: "usage", Handler: func(a *App, args []string) error {
			firstCalled += 1
			if len(args) == 0 {
				firstArgsMatchExpected = true
			}
			return nil
		}},
		{Name: "mycommand2", Usage: "usage", Handler: func(a *App, args []string) error {
			secondCalled += 1
			if len(args) == 2 && args[0] == "Hello" && args[1] == "Second argument" {
				secondArgsMatchExpected = true
			}
			return nil
		}}}
	// Make sure that nothing gets called initially
	if firstCalled != 0 || secondCalled != 0 || firstArgsMatchExpected || secondArgsMatchExpected {
		t.Error("Routing called commands before the Handle was called")
	}

	// Make sure that nothing is still called on an nonexisting command
	err := app.Run([]string{"hi"})
	if err != nil || firstCalled != 0 || secondCalled != 0 || firstArgsMatchExpected || secondArgsMatchExpected {
		t.Error("Routing called commands before the Handle was called")
	}

	// Make sure that only the first handler is called
	err = app.Run([]string{"mycommand"})
	if err != nil || firstCalled != 1 || !firstArgsMatchExpected || secondCalled != 0 || secondArgsMatchExpected {
		t.Error("First handler call was not propagated properly")
	}

	err = app.Run([]string{"mycommand2", "Hello", "Second argument"})
	if err != nil || firstCalled != 1 || !firstArgsMatchExpected || secondCalled != 1 || !secondArgsMatchExpected {
		t.Error("Second handler call was not propagated properly")
	}
}
