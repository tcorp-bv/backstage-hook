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

import "time"

// Key value store to remember the last policy (used for allow always).
type Policies interface {
	// Store some value in key. If val == StoredPolicy{}, this will delete the policy.
	Store(key string, val StoredPolicy)

	// Get the value of the key, second argument is false if nonexistent
	Get(key string) (StoredPolicy, bool)
}

// The storage representation of a policy.
type StoredPolicy struct {
	// The Id value of the policy
	PolicyId string
	// The timestamp of when this policy was set (the time the user selected "AllowAlways")
	Timestamp time.Time
}
