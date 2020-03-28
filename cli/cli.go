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

package cli

// Extremely simple command line interface handler
// Inspired by https://github.com/urfave/cli

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	// Ansi code to clear the whole terminal screen
	ClearScreen EscapeCode = "\033[2J"
	// Ansi code to clear the whole line of the cursor
	ClearLine EscapeCode = "\033[2K"
	// Ansi code to clear the whole line to the right of the cursor
	ClearLineRight EscapeCode = "\033[1K"
	// Ansi code for formatting text green fmt eg fmt.Sprintf(cli.GreenColor, "Hello")
	GreenColor Color = "\033[1;32m"
	// Ansi code for formatting text yellow with fmt
	YellowColor Color = "\033[1;33m"
	// Ansi code for formatting text blue with fmt
	BlueColor Color = "\033[1;34m"
	// Ansi code for formatting text Red with fmt
	RedColor Color = "\033[1;31m"
	// Ansi code for formatting text White with fmt
	WhiteColor Color = "\033[1;37m"
	// Save the cursor location
	CursorSave EscapeCode = "\033[s"
	// Move the cursor to the saved cursor location
	CursorRestore EscapeCode = "\033[u"
	// Move the cursor to the bottom of the terminal
	CursorBottom EscapeCode = "\033[100B"
	// Move the cursor to the first column of the terminal
	CursorLeft   EscapeCode = "\033[100D"
	// Scroll to the bottom of the terminal
	ScrollBottom EscapeCode = "\033[3J\033[30;40m"
)

type EscapeCode string

func (c EscapeCode) String() string {
	return string(c)
}

// Ansi color to highlight text in terminal.
type Color string

// Returns the ansi string that can be used inline.
func (c Color) String() string {
	return string(c)
}

// Ansi markup a string: Can be used like println(c.Format("Hello, World")).
func (c Color) Format(text string) string {
	return fmt.Sprintf(c.FmtString(), text)
}

// String for fmt formatted text, can be used as fmt.Sprintf(c.FmtString, "Hello, world")
func (c Color) FmtString() string {
	return string(c) + "%s\033[0m"
}

// Ansi code to move the cursor in the terminal to said row and column.
func CursorTo(row int, col int) EscapeCode {
	return EscapeCode(fmt.Sprintf("\033[%d;%dH", row, col))
}

// Ansi code to move the cursor up row rows.
func CursorUp(rows int) EscapeCode {
	return EscapeCode(fmt.Sprintf("\033[%dA", rows))
}

// Ansi code to scroll the cursor down n rows.
func CursorDown(rows int) EscapeCode {
	return EscapeCode(fmt.Sprintf("\033[%dB", rows))
}

// Ansi code to set the title of the terminal.
func Title(title string) EscapeCode {
	return EscapeCode(fmt.Sprintf("\x1b]0;%s\x07", title))
}

// App is the main structure of a cli application. It is recommended that a new App is created through cli.newApp().
type App struct {
	// The name of the program
	Name string
	// Description of the program
	Usage string
	// Description of the program argument format
	ArgsUsage string
	// Version of the program
	Version string
	// List of the commands to execute
	Commands []*Command
	// Author of this application
	Author string

	// Reader can read the user's input
	Reader io.Reader

	// Writer to write output to
	Writer io.Writer

	// ErrWritter to write error output to
	ErrWriter io.Writer
}

// Handler is the actual handler when the user executes a command in the cli
type Handler func(a *App, args []string) error

// Represents a command that the user executes in the cli (eg. when the user types backstage-hook start, start is the command)
type Command struct {
	// The command string (eg. "start")
	Name string
	// The usage explanation
	Usage string
	// The command handler that manages the actual behavior of the command
	Handler Handler
}

// Prints the help menu to the app's writer (terminal)
func (a *App) Help() {
	fmt.Fprintf(a.Writer, "NAME:\n    %s - %s\n\n", a.Name, a.Usage)
	fmt.Fprintf(a.Writer, "USAGE:\n    %s\n\n", a.ArgsUsage)
	fmt.Fprintf(a.Writer, "COMMANDS:\n")
	for _, c := range a.Commands {
		fmt.Fprintf(a.Writer, "    %s  %s\n", c.Name, c.Usage)
	}
}

// Executes the cli app, parses the arguments to the relevant command.
func (a *App) Run(arguments []string) (err error) {
	if len(arguments) == 0 {
		a.Help()
		return nil
	}
	cmd := arguments[0]
	if cmd == "-h" || cmd == "--help" {
		a.Help()
		return nil
	}

	cmdArgs := []string{}
	if len(arguments) > 1 {
		cmdArgs = arguments[1:]
	}

	isCommand, err := a.executeIfCommand(cmd, cmdArgs)
	if isCommand || err != nil { // return if the command executed or the error is not nil
		return err
	}
	// If no command was executed, display the help menu.
	a.Help()
	return nil
}

func (a *App) executeIfCommand(cmd string, args []string) (isCommand bool, err error) {
	for _, c := range a.Commands {
		if c.Name == cmd {
			return true, a.executeCommand(c, args)
		}
	}
	return false, nil
}

// Executes a command's Handler
func (a *App) executeCommand(c *Command, args []string) error {
	if c.Handler == nil {
		return errors.New(fmt.Sprintf("%v does not have a handler.", c.Name))
	}
	return c.Handler(a, args)
}

// Waits for input (until the user types and presses enter) and returns it (including the newline)
func (a *App) GetInput() string {
	var in string
	_, err := fmt.Fscanln(a.Reader, &in)
	if err != nil {
		log.Fatal(err)
	}
	return in
}

// Creates a new instance of App with some default values for the name and writer.
func NewApp() *App {
	return &App{
		Name:   filepath.Base(os.Args[0]),
		Usage:  "A new cli application",
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
}
