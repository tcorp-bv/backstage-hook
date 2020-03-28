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
	"testing"
	"time"
)

// See policies_test.go for integration testing

// Make sure that a policy is properly deleted when you store StoredPolicy{}
func TestPolicyDeletion(t *testing.T) {
	s := NewMemoryPolicyStorage()
	pol := StoredPolicy{PolicyId: "ALLOW", Timestamp: time.Now()}
	key := "somehash"

	val, contains := s.Get(key)
	if contains || (val != StoredPolicy{}) {
		t.Error("Policy should not be fetchable when it was not stored.")
	}

	s.Store(key, pol)
	val, contains = s.Get(key)
	if !contains || val != pol {
		t.Error("Get did not return expected results.")
	}

	s.Store(key, StoredPolicy{})
	val, contains = s.Get(key)
	if contains || (val != StoredPolicy{}) {
		t.Error("Policy should not be fetchable after it was removed.")
	}
}

// Make sure that a session is properly deleted when you store StoredSession{}
func TestSessionDeletion(t *testing.T) {
	s := NewMemorySessionStorage()
	session := StoredSession{"somesecret", time.Now()}
	key := "somehash"

	val, contains := s.Get(key)
	if contains || (val != StoredSession{}) {
		t.Error("Session should not be fetchable when it was not stored.")
	}

	s.Store(key, session)
	val, contains = s.Get(key)
	if !contains || val != session {
		t.Error("Get did not return expected results.")
	}

	s.Store(key, StoredSession{})
	val, contains = s.Get(key)
  if contains || (val != StoredSession{}) {
		t.Error("Session should not be fetchable after it was removed.")
	}
}
