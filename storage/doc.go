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

/*
	The storage package manages all persistent data for the backstage Handler.
    Currently this data consists of policies (eg. ALLOW_ALWAYS) and sessions. A
    policy is identified by an Action (That is the combination of the command,
    its arguments and the plugin). A session is identified by its Id.
*/

/*
	How to use the storage package?

	Currently the storage package provides access to Policy storage and Session
    storage. To get started, you must call New(...) with the implementations for
    these you wish to use, the returned object allows you to get and set the
    relevant data.
*/

/*
	How to implement a new storage backend?

	To store policies, your backend must implement the Policies interface in
    policies.go. To store sessions, it has to implement the Sessions interface
    in sessions.go.

    The storage package contains some other interfaces (Store, PoliciesStore...)
    but these are meant for external access and provide the external interface.
    You should thus not re-implement these.
*/
