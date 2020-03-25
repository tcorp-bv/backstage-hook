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

import "testing"

// Check that there are no collisions between a bunch of unique actions.
func TestHashUniqueness(t *testing.T) {
	actions := generateUniqueActions()
	hashes := map[string]bool{}

	for _, action := range actions {
		hash1 := action.Hash()
		hash2 := action.Hash()

		if hash1 != hash2 {
			t.Error("Hash should equal itself")
		}
		hashes[hash1] = true
	}

	if len(hashes) != len(actions) {
		t.Error("action hash function does not work: expected ", len(actions), " hashes, hashed to ", len(hashes), " hashes")
	}
}

// Tests that an action with nil args with collide with []string{} args.
func TestCollisionArgs(t *testing.T) {
	h1 := Action{Plugin: "test", Command: Command{Name: "name", Args: nil}}.Hash()
	h2 := Action{Plugin: "test", Command: Command{Name: "name", Args: []string{}}}.Hash()

	if h1 != h2 {
		t.Error("hashes of nil args and []string{} args should be equal")
	}
}

var strings = []string{"", "test", "tests", "t", "{\"name\":\"test\"}", "Test", "*"}

// Generates a set of unique actions.
func generateUniqueActions() []Action {
	var actions []Action
	for _, plugin := range strings {
		for _, commandName := range strings {
			for argNum := 0; argNum < 3; argNum++ {
				args := []string{}
				for i := 0; i < argNum; i++ {
					args = append(args, strings...)
				}
				actions = append(actions, Action{Plugin: plugin, Command: Command{Name: commandName, Args: args}})
			}
		}
	}
	return actions
}
