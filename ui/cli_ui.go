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
	"fmt"
	"github.com/tcorp-bv/backstage-hook/actions"
	"github.com/tcorp-bv/backstage-hook/cli"
	"github.com/tcorp-bv/backstage-hook/policies"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// TODO: See https://github.com/tcorp-bv/backstage-hook/issues/1: Make the CLI compatible with small terminals (in #characters)

const (
	promptHeight = 6
)

// Contains all relevant properties of an action request.
type requestResponse struct {
	// The actual request
	Req actions.Action
	// The response (eg. AllowAlways)
	Res chan policies.Policy
	// The command is stored in a temporary file for the user to view. This is an extra measurement against injecting a command that hides itself through console properties.
	File *os.File
}

// Generates the file in the temporary directory if nonexistent and returns its url
func (r *requestResponse) FileURI() string {
	if r.File == nil {
		r.generateFile()
	}
	path, err := filepath.Abs(r.File.Name())
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("file://%s", path)
}

// Generates a new file and sets the command as its contents.
func (r *requestResponse) generateFile() {
	f, err := ioutil.TempFile(os.TempDir(), "*-command.txt")
	defer func() {	
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	r.File = f
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(r.Req.Command.String())
	if err != nil {
		log.Fatal(err)
	}
}

// Creates the command-line interface UI frontend for the hook.
func NewCli(app *cli.App) UI {
	return &cliUI{App: app, queue: []requestResponse{}}
}

// Interactive CLI frontend for backstage-hook
type cliUI struct {
	sync.Mutex
	App   *cli.App
	queue []requestResponse // Queued actions
}

// Writes the header to the top of the output
func (c *cliUI) writeHeader() {
	fmt.Fprintf(c.App.Writer, "%s", cli.CursorTo(1, 1))
	fmt.Fprintf(c.App.Writer, "%s%s%s", cli.ClearLine, cli.YellowColor.Format(c.App.Name+"\n"), cli.ClearLine)
	fmt.Fprintf(c.App.Writer, "%s%d actions are waiting for your approval\n%s\n", cli.ClearLine.String(), len(c.queue), cli.ClearLine)
	fmt.Fprintf(c.App.Writer, "%s%s%s", cli.ClearLine, cli.YellowColor.Format("QUEUED actions:\n"), cli.ClearLine)
}

// Writes the queued actions to the body of the output
func (c *cliUI) writeQueue() {
	if len(c.queue) <= 1 {
		return
	}

	//Move cursor just above the prompt
	fmt.Fprintf(c.App.Writer, "%s%s", cli.CursorBottom, cli.CursorUp(promptHeight))
	for i := 1; i < len(c.queue); i++ {
		req := c.queue[i].Req
		// Move a line up and print a truncated version of the command and plugin
		fmt.Fprintf(c.App.Writer, "%s%s%d. %.20q by %q", cli.CursorUp(1).String(), cli.CursorLeft, i+1, req.Command.Name+"...", req.Plugin)
	}
}

// Writes the command prompt (the actual request for a policy) to the bottom of the prompt.
func (c *cliUI) writePrompt() {
	if len(c.queue) == 0 {
		return
	}
	fmt.Fprintf(c.App.Writer,
		"%s%s%s", cli.CursorBottom, cli.CursorLeft, cli.CursorUp(promptHeight-1))
	fmt.Fprintf(c.App.Writer, "%s By %q:\n\n", cli.GreenColor.Format("1. NEW REQUEST"+string(len(c.queue))), c.queue[0].Req.Plugin)
	fmt.Fprintf(c.App.Writer, "     %s\n", cli.WhiteColor.Format(fmt.Sprintf("%.100q", c.queue[0].Req.Command.String())))
	fmt.Fprintf(c.App.Writer, "Full command at %s\n\n", c.queue[0].FileURI()) // Todo: check behavior of this when previous line overflows
	fmt.Fprintf(c.App.Writer, "%s/%s/%s: ", decisionString(policies.Deny()), decisionString(policies.Allow()), decisionString(policies.AllowAlways()))
}

// Sets up the command line interface. This includes setting the title and clearing the screen.
func (c *cliUI) Setup() {
	c.Lock()
	defer c.Unlock()
	fmt.Fprintf(c.App.Writer, "%s", cli.Title(c.App.Name).String())
	fmt.Fprintf(c.App.Writer, "%s", cli.ClearScreen.String())

	c.writeHeader()

}

// Handles an incoming action request by adding it to the queue and updating the display.
func (c *cliUI) Handle(req actions.Action, res chan policies.Policy) {
	c.queue = append(c.queue, requestResponse{Req: req, Res: res})
	c.handleQueue()
}

// Updates the display and starts the display prompt if a new item requires approval.
func (c *cliUI) handleQueue() {
	c.Lock()
	defer c.Unlock()
	c.render()

	if len(c.queue) == 0 {
		return
	}
	if len(c.queue) > 1 {
		c.updateDisplayedQueue()
	} else {
		c.updateDisplayedPrompt()
	}
}

// Updates the command queue in the center of the commandline
func (c *cliUI) updateDisplayedQueue() {
	fmt.Fprint(c.App.Writer, cli.CursorSave)
	c.writeQueue()
	c.writeHeader()
	fmt.Fprint(c.App.Writer, cli.CursorRestore)
}

// Updates the approval prompt at the bottom of the command line and handles the approval flow.
func (c *cliUI) updateDisplayedPrompt() {
	if len(c.queue) == 0 {
		return
	}

	c.writePrompt()
	go c.handlePrompt()
}

// Handles the actual prompt input
func (c *cliUI) handlePrompt() {
	shortcut := c.getShortcutInput()
	policy, err := policies.ByShortcut(shortcut)
	if err != nil {
		log.Fatal(err)
	}
	c.Lock()
	defer c.Unlock()
	c.queue[0].Res <- policy // Send the policy to the response channel
	err = os.Remove(c.queue[0].File.Name())
	if err != nil {
		log.Fatal(err)
	}
	c.queue = c.queue[1:] // Pop
	if len(c.queue) > 0 {
		c.render()
		go c.updateDisplayedPrompt()
	}
}

// Re-renders the complete command line interface view.
func (c *cliUI) render() {
	fmt.Fprintf(c.App.Writer, "%s%s", cli.ClearScreen.String(), cli.CursorTo(1, 1))
	c.writeQueue()
	c.writeHeader()
	c.writePrompt()
}

// Wait for the user to enter an input and return it. It will retry and re-render if the input is not a shortcut.
func (c *cliUI) getShortcutInput() string {
	var in string
	for {
		in = c.App.GetInput()
		if policies.ShortcutValid(in) {
			return in
		}
		c.render()
	}
}

// Maps the policies to a color for display purposes.
var colorMap = map[policies.Policy]cli.Color{
	policies.Allow():       cli.GreenColor,
	policies.AllowAlways(): cli.GreenColor,
	policies.Deny():        cli.RedColor,
}

// Gets the colorized decision string of a policy. This is used in the prompt, eg. "Allow (a)" or "Deny (d)".
func decisionString(pol policies.Policy) string {
	color, hasColor := colorMap[pol]
	if !hasColor {
		color = cli.WhiteColor
	}
	return color.Format(fmt.Sprintf("%s (%s)", pol.Name(), pol.Shortcut()))
}
