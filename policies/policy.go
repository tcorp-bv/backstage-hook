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

package policies

import (
	"errors"
)

var (
	// Policy to allow an action once
	allow = policy{id: "ALLOW", name: "Allow", description: "Allow the action once", shortcut: "a"}
	// Policy to always allow the action in the future
	allowAlways = policy{id: "ALLOW_ALWAYS", name: "Always allow", description: "Always allow this action", shortcut: "s"}
	// Policy to deny this action this time
	deny = policy{id: "DENY", name: "Deny", description: "Deny this action this time", shortcut: "d"}
)

// Returns a list containing allow, allowAlways and deny.
func All() []Policy {
	return []Policy{Allow(), AllowAlways(), Deny()}
}

// Allow is the policy to allow an action once.
func Allow() Policy {
	return allow
}

// AllowAlways is the policy to always allow the action in the future.
func AllowAlways() Policy {
	return allowAlways
}

// Deny is the policy to deny the action this time.
func Deny() Policy {
	return deny
}

// Check if one of the three policies has the given shortcut (which is a single character id used in the cli).
func ShortcutValid(sc string) bool {
	for _, pol := range All() {
		if pol.Shortcut() == sc {
			return true
		}
	}
	return false
}

// Check if one of the three policies has the given Id.
func IdValid(id string) bool {
	for _, pol := range All() {
		if pol.Id() == id {
			return true
		}
	}
	return false
}

// Gets the policy by the given shortcut, returns an error if the shortcut is not found.
func ByShortcut(sc string) (Policy, error) {
	for _, p := range All() {
		if p.Shortcut() == sc {
			return p, nil
		}
	}
	return policy{}, errors.New("policy with shortcut not found")
}

// Gets the policy by the given id, returns an error if the id is not found.
func ById(id string) (Policy, error) {
	for _, pol := range All() {
		if pol.Id() == id {
			return pol, nil
		}
	}
	return policy{}, errors.New("policy with id not found")
}

// The decision a user can make when a new action comes in.
type Policy interface {
	// The identifier for database storage
	Id() string
	// The display name
	Name() string
	// The display description
	Description() string
	// One character shortcut that can be used for the command line as an id
	Shortcut() string
}

// The private implementation for Policy. This makes the fields immutable by other packages.
type policy struct {
	// The identifier for database storage
	id string
	// The display name
	name string
	// The display description
	description string
	// One letter shortcut that can be used for the command line
	shortcut string
}

// See the Policy interface.
func (p policy) Id() string {
	return p.id
}

// See the Policy interface.
func (p policy) Name() string {
	return p.name
}

// See the Policy interface.
func (p policy) Description() string {
	return p.description
}

// See the Policy interface.
func (p policy) Shortcut() string {
	return p.shortcut
}
