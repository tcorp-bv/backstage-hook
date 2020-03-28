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

package storage

import (
	"github.com/tcorp-bv/backstage-hook/actions"
	"github.com/tcorp-bv/backstage-hook/policies"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	// Runs a bunch of integration tests against the memory storage
	testPolicyStorage(t, New(NewMemoryPolicyStorage(), nil))
}

func testPolicyStorage(t *testing.T, s PoliciesStore) {
	acts := generateUniqueActions()
	// Make sure that no policies are in the store on the start
	for _, act := range acts {
		assertPolicyNotExist(t, s, act)
	}
	// Set the first action and ensure it is set
	s.SetPolicy(acts[0], policies.Allow())
	assertPolicyValue(t, s, acts[0], policies.Allow())

	// Ensure that the second action still does not exist
	assertPolicyNotExist(t, s, acts[1])

	// Ovverride the first policy and ensure it was overriden
	s.SetPolicy(acts[0], policies.Deny())
	assertPolicyValue(t, s, acts[0], policies.Deny())

	// Ensure that the second action still does not exist
	assertPolicyNotExist(t, s, acts[1])

	// Delete the first policy and ensure it is deleted
	s.SetPolicy(acts[0], nil)
	assertPolicyNotExist(t, s, acts[0])

	// Make sure that setting a bunch of values work
	for _, act := range acts {
		s.SetPolicy(act, policies.AllowAlways())
		assertPolicyValue(t, s, act, policies.AllowAlways())
	}
	// Make sure that deleting all values works
	for _, act := range acts {
		s.SetPolicy(act, nil)
	}
	for _, act := range acts {
		assertPolicyNotExist(t, s, act)
	}
}

func assertPolicyValue(t *testing.T, s PoliciesStore, act actions.Action, pol policies.Policy) {
	p, contains := s.Policy(act)
	if p == nil || p != pol || !contains {
		t.Error("Get did not return expected r1 values")
	}
}
func assertPolicyNotExist(t *testing.T, s PoliciesStore, act actions.Action) {
	p, contains := s.Policy(act)
	if p != nil || contains {
		t.Error("Get did not return expected empty values")
	}
}

var strings = []string{"", "test", "tests", "t", "{\"name\":\"test\"}", "Test", "*"}

// Generates a set of unique actions.
func generateUniqueActions() []actions.Action {
	var acts []actions.Action
	for _, plugin := range strings {
		for _, commandName := range strings {
			for argNum := 0; argNum < 3; argNum++ {
				var args []string
				for i := 0; i < argNum; i++ {
					args = append(args, strings...)
				}
				acts = append(acts, actions.Action{Plugin: plugin, Command: actions.Command{Name: commandName, Args: args}})
			}
		}
	}
	return acts
}
