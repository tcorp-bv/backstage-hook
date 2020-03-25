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
	"testing"
)

// Tests that byId and byShortcut return the correct policies for all policies in All().
func TestGetters(t *testing.T) {
	for _, pol := range All() {
		byId, _ := ById(pol.Id())
		byShortcut, _ := ByShortcut(pol.Shortcut())
		if byId != pol {
			t.Error("method byId does not work")
		}
		if byShortcut != pol {
			t.Error("method ByShortcut does not work")
		}
	}
	_, err := ById("Some nonexistent id")
	if err == nil {
		t.Error("byId did not throw error")
	}
	_, err2 := ByShortcut("Some nonexistent shortcut")
	if err2 == nil {
		t.Error("byShortcut did not throw error")
	}
}

// Tests that existing policies are considered valid and non-existing policies are not.
func TestValidity(t *testing.T) {
	for _, pol := range All() {
		if !ShortcutValid(pol.Shortcut()) {
			t.Error("shortcutValid does not work")
		}
		if !IdValid(pol.Id()) {
			t.Error("idValid does not work")
		}
	}
	if IdValid("Some nonexistent id") {
		t.Error("idValid should not return true here")
	}
	if ShortcutValid("Some nonexistent shortcut") {
		t.Error("shortcutValid should not return true here")
	}
}

// Tests that All() contains Allow(), Deny() and AllowAlways() exactly once.
func TestPoliciesContents(t *testing.T) {
	var numAllow, numDeny, numAllowAlways int
	for _, pol := range All() {
		switch pol {
		case Allow():
			numAllow += 1
		case Deny():
			numDeny += 1
		case AllowAlways():
			numAllowAlways += 1
		}
	}
	if numAllow != 1 || numDeny != 1 || numAllowAlways != 1 {
		t.Error("All() does not contain the 3 base policies exactly once.")
	}
}

// Tests that the Id and Shortcut of all policies in All() are unique.
func TestUniqueIDS(t *testing.T) {
	ids, shortcuts := map[string]bool{}, map[string]bool{}
	for _, pol := range All() {
		ids[pol.Id()] = true
		shortcuts[pol.Shortcut()] = true
	}
	if len(ids) != len(All()) {
		t.Error("Policies with duplicate id")
	}
	if len(shortcuts) != len(All()) {
		t.Error("Policies with duplicate shortcuts")
	}
}

// Tests that the Policy interface getters match the actual values defined in policy.
func TestPolicyGetMethods(t *testing.T) {
	policies := map[Policy]policy{Allow(): allow, Deny(): deny, AllowAlways(): allowAlways}
	for P, p := range policies {
		if P.Name() != p.name || P.Shortcut() != p.shortcut || P.Id() != p.id || P.Description() != p.description {
			t.Error("Policy interface get method does not match policy value")
		}
	}
}
