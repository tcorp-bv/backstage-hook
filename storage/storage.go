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
	"github.com/tcorp-bv/backstage-hook/sessions"
	"time"
)

// Instantiates a Store object with a Policies and Sessions backend.
// Example usage: storage.New(storage.NewMemoryPolicyStorage(), storage.NewMemorySessionStorage()))
func New(pol Policies, ses Sessions) Store {
	return &store{policyStore: pol, SessionStore: ses}
}

// The external storage interface to store policies and sessions.
type Store interface {
	// Inherit the PoliciesStore interface
	PoliciesStore
	// Inherit the SessionsStore interface
	SessionsStore
}

// External interface to get and set policies to the backend.
type PoliciesStore interface {
	// Returns the stored policy and true if the storage contains the policy
	Policy(a actions.Action) (p policies.Policy, contains bool)

	// Stores the given action to the storage. If p == nil, this will delete the policy.
	SetPolicy(a actions.Action, p policies.Policy)
}

// External interface to get and set sessions to the backend.
type SessionsStore interface {
	// Returns the stored session and true if the storage contains the session
	Session(id string) (s sessions.Session, contains bool)

	// Stores the given session to the storage
	SetSession(s sessions.Session)

	//Removes the session with the provided id if it exists
	DeleteSession(id string)
}

type store struct {
	policyStore  Policies
	SessionStore Sessions
}

func (s *store) Policy(a actions.Action) (policies.Policy, bool) {
	val, ok := s.policyStore.Get(a.Hash())
	if !ok || (val == StoredPolicy{}) || !policies.IdValid(val.PolicyId) {
		return nil, false
	}
	p, _ := policies.ById(val.PolicyId)
	return p, true
}

func (s *store) SetPolicy(a actions.Action, pol policies.Policy) {
	if pol == nil {
		s.policyStore.Store(a.Hash(), StoredPolicy{})
		return
	}
	s.policyStore.Store(a.Hash(), StoredPolicy{PolicyId: pol.Id(), Timestamp: time.Now()})
}

func (s *store) Session(id string) (session sessions.Session, contains bool) {
	val, ok := s.SessionStore.Get(id)
	if !ok || (val == StoredSession{}) {
		return nil, false
	}
	return sessions.New(id, val.Secret), true
}

func (s *store) SetSession(session sessions.Session) {
	if session == nil {
		panic("tried to store a nil session in the session storage, this should not happen.")
	}
	s.SessionStore.Store(session.Id(), StoredSession{Secret: session.Secret(), Created: time.Now()})
}

func (s *store) DeleteSession(id string) {
	s.SessionStore.Store(id, StoredSession{})
}
