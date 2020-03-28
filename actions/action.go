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

package actions

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// The command that should be executed.
type Command struct {
	// The command
	Name string `json:"name"`
	// The arguments for the command
	Args []string `json:"args,omitempty"`
}

func (c *Command) String() string {
	if len(c.Args) == 0 {
		return c.Name
	}
	return fmt.Sprintf("%s %s", c.Name, strings.Join(c.Args, " "))
}
// An action is the intent by a plugin to execute a command.
type Action struct {
	// The command to execute
	Command Command `json:"command"`
	// The plugin that supposedly executed this action
	Plugin string `json:"plugin"`
}

// Gets the unique hash of this action. This is achieved by encoding the struct as json and getting the base64 of the sha256sum of this json.
// Note that as this uses the json parser, these hashes may not necessarily be the same across go versions. There should however never be collisions.
func (a Action) Hash() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	hash := sha256.New()
	return base64.StdEncoding.EncodeToString(hash.Sum(bytes))
}
