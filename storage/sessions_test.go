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
	"github.com/tcorp-bv/backstage-hook/sessions"
	"testing"
)

func TestMemorySessionStorage(t *testing.T) {
	// Runs a bunch of integration tests against the memory storage
	testSessionStorage(t, New(nil, NewMemorySessionStorage()))
}

func testSessionStorage(t *testing.T, s SessionsStore) {
	sess := generateUniqueSessions()
	// Make sure that no sessions are in the store on the start
	for _, ses := range sess {
		assertSessionNotExist(t, s, ses)
	}
	// Set the first session and ensure it is set
	s.SetSession(sess[0])
	assertSessionValue(t, s, sess[0])

	// Ensure that the second session still does not exist
	assertSessionNotExist(t, s, sess[1])

	sess[0] = sessions.New(sess[0].Id(), "Some new secret")
	// Override the first session and ensure it was overridden
	s.SetSession(sess[0])
	assertSessionValue(t, s, sess[0])

	// Ensure that the second session still does not exist
	assertSessionNotExist(t, s, sess[1])

	// Delete the first session and ensure it is deleted
	s.DeleteSession(sess[0].Id())
	assertSessionNotExist(t, s, sess[0])

	// Make sure that setting a bunch of sessions work
	for _, ses := range sess {
		s.SetSession(ses)
		assertSessionValue(t, s, ses)
	}
	// Make sure that deleting all sessions works
	for _, ses := range sess {
		s.DeleteSession(ses.Id())
	}
	for _, ses := range sess {
		assertSessionNotExist(t, s, ses)
	}
}

func assertSessionValue(t *testing.T, s SessionsStore, ses sessions.Session) {
	storedSes, contains := s.Session(ses.Id())
	if storedSes == nil || storedSes.Id() != ses.Id() || storedSes.Secret() != ses.Secret() || !contains {
		t.Error("Get did not return expected session values")
	}
}
func assertSessionNotExist(t *testing.T, s SessionsStore, ses sessions.Session) {
	p, contains := s.Session(ses.Id())
	if p != nil || contains {
		t.Error("Get did not return expected empty session values")
	}
}

var (
	sesSecrets = []string{"", "test", "tests", "t", "{\"name\":\"test\"}", "Test", "*"}
	sesIds     = []string{"test", "t", "", "{\"name\":\"test\"}", "*", "Test"}
)

// Generates a set of unique sessions.
func generateUniqueSessions() []sessions.Session {
	var sess []sessions.Session
	for i, id := range sesIds {
		sess = append(sess, sessions.New(id, sesSecrets[i]))
	}
	return sess
}
