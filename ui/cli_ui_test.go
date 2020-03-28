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

package ui

import (
	"bytes"
	"github.com/tcorp-bv/backstage-hook/cli"
	"strings"
	"testing"
)

//Todo: Test that no ANSI escape codes can be injected
func TestInjectionProtection(t *testing.T) {
	// text = "kubectl get pods -o yaml \033[30m; echo some injection || \033[0m\n"
}

func TestPromptHeight(t *testing.T) {
	var buf bytes.Buffer
	ui := cliUI{App: &cli.App{Writer: &buf}}
	ui.queue = []requestResponse{requestResponse{Req: actions.Action{Command:actions.Command{Name:"test"}, Plugin:"testplugin"}}}

	ui.writePrompt()

	lines := strings.Count(buf.String()), "\n") + 1
	if lines != promptHeight {
		t.Error("Prompt is ", lines, " lines, expected ", promptHeight, " lines.")
	}

}
