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

import "sync"

// Simple in-memory Policies storage implementation.
type memoryPolicyStorage struct {
	// Ensures that multi-threaded access is synchronized through locking and unlocking the storage
	sync.Mutex
	// The map where all the policies are stored
	storage map[string]StoredPolicy
}

// Simple in-memory Sessions storage implementation.
type memorySessionStorage struct {
	// Ensures that multi-threaded access is synchronized through locking and unlocking the storage
	sync.Mutex
	// The map where all the sessions are stored
	storage map[string]StoredSession
}

func (m *memoryPolicyStorage) Store(key string, value StoredPolicy) {
	m.Lock()
	defer m.Unlock()
	if (value == StoredPolicy{}) { // In case the StoredPolicy is empty, the policy should be deleted.
		delete(m.storage, key)
		return
	}
	m.storage[key] = value
}
func (m *memoryPolicyStorage) Get(key string) (StoredPolicy, bool) {
	m.Lock()
	defer m.Unlock()
	value, ok := m.storage[key]
	return value, ok
}

func (m *memorySessionStorage) Store(key string, value StoredSession) {
	m.Lock()
	defer m.Unlock()
	if (value == StoredSession{}) { // In case the StoredSession is empty, the session should be deleted.
		delete(m.storage, key)
		return
	}
	m.storage[key] = value
}

func (m *memorySessionStorage) Get(key string) (StoredSession, bool) {
	m.Lock()
	defer m.Unlock()
	value, ok := m.storage[key]
	return value, ok
}

// Returns a simple, in-memory Policies implementation. Data will not persist on restart.
func NewMemoryPolicyStorage() Policies {
	return &memoryPolicyStorage{storage: map[string]StoredPolicy{}}
}

// Returns a simple, in-memory Sessions implementation. Data will not persist on restart.
func NewMemorySessionStorage() Sessions {
	return &memorySessionStorage{storage: map[string]StoredSession{}}
}
