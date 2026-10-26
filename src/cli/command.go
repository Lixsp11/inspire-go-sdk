// Package cli provides command routing, flags and machine-readable diagnostics.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type AssignedMode int

const (
	AssignedNone       AssignedMode = -1
	AssignedOnce       AssignedMode = 1
	AssignedRepeatable AssignedMode = 9
)

type Flag struct {
	Name         string       `json:"name"`
	Short        string       `json:"description"`
	AssignedMode AssignedMode `json:"assign_mode"`
	Persistent   bool         `json:"-"`
	Category     string       `json:"category,omitempty"`
	value        string
	values       []string
	assigned     bool
}

func (f *Flag) GetValue() (string, bool) { return f.value, f.assigned }
func (f *Flag) GetValues() []string      { return append([]string(nil), f.values...) }
func (f *Flag) IsAssigned() bool         { return f != nil && f.assigned }
func (f *Flag) GetStringOrDefault(def string) string {
	if f.IsAssigned() {
		return f.value
	}
	return def
}

type FlagSet struct{ flags map[string]*Flag }

func NewFlagSet() *FlagSet                { return &FlagSet{flags: map[string]*Flag{}} }
func (fs *FlagSet) Add(f *Flag)           { fs.flags[f.Name] = f }
func (fs *FlagSet) Get(name string) *Flag { return fs.flags[name] }
func (fs *FlagSet) GetValue(name string) (string, bool) {
	if f := fs.Get(name); f != nil {
		return f.GetValue()
	}
	return "", false
}
func (fs *FlagSet) Flags() []*Flag {
	names := make([]string, 0, len(fs.flags))
	for name := range fs.flags {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]*Flag, 0, len(names))
	for _, name := range names {
		result = append(result, fs.flags[name])
	}
	return result
}

type Context struct {
	Context             context.Context
	Stdin               io.Reader
	Getenv              func(string) string
	stdout, stderr      io.Writer
	command             *Command
	flags, unknownFlags *FlagSet
}

func NewCommandContext(stdout, stderr io.Writer) *Context {
	return &Context{Context: context.Background(), Stdin: os.Stdin, Getenv: os.Getenv, stdout: stdout, stderr: stderr}
}
func (ctx *Context) Stdout() io.Writer      { return ctx.stdout }
func (ctx *Context) Stderr() io.Writer      { return ctx.stderr }
func (ctx *Context) Command() *Command      { return ctx.command }
func (ctx *Context) Flags() *FlagSet        { return ctx.flags }
func (ctx *Context) UnknownFlags() *FlagSet { return ctx.unknownFlags }

type Command struct {
	Name              string
	Short             string
	Usage             string
	RequestSchema     interface{}
	EnableUnknownFlag bool
	Run               func(*Context, []string) error
	parent            *Command
	subCommands       []*Command
	flags             *FlagSet
}

func (c *Command) AddSubCommand(cmd *Command) {
	cmd.parent = c
	c.subCommands = append(c.subCommands, cmd)
}
func (c *Command) Flags() *FlagSet {
	if c.flags == nil {
		c.flags = NewFlagSet()
	}
	return c.flags
}
func (c *Command) GetSubCommand(name string) *Command {
	for _, child := range c.subCommands {
		if child.Name == name {
			return child
		}
	}
	return nil
}
func (c *Command) getName() string {
	if c.parent == nil {
		return c.Name
	}
	return c.parent.getName() + " " + c.Name
}

func (c *Command) Execute(ctx *Context, args []string) error {
	// Find the route without treating option values as command names.
	cmd := c
	positionals := []string{}
	values := map[string][]string{}
	allFlags := map[string]*Flag{}
	var collect func(*Command)
	collect = func(node *Command) {
		for _, f := range node.Flags().Flags() {
			allFlags[f.Name] = f
		}
		for _, child := range node.subCommands {
			collect(child)
		}
	}
	collect(c)
	help := false
	if len(args) > 0 && args[0] == "help" {
		help = true
		args = args[1:]
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-h" || arg == "--help" {
			help = true
			continue
		}
		if !strings.HasPrefix(arg, "--") {
			if strings.HasPrefix(arg, "-") {
				return UsageError("InvalidFlag", "Use long option names beginning with --.")
			}
			positionals = append(positionals, arg)
			continue
		}
		name, value, inline := strings.Cut(strings.TrimPrefix(arg, "--"), "=")
		if name == "" {
			return UsageError("InvalidFlag", "An option name is required.")
		}
		f := allFlags[name]
		mode := AssignedOnce
		if f != nil {
			mode = f.AssignedMode
		}
		if _, duplicate := values[name]; duplicate && mode != AssignedRepeatable {
			return UsageError("InvalidParameter", "Duplicate option --"+name+".")
		}
		if mode == AssignedNone {
			if inline {
				return UsageError("InvalidParameter", "Option --"+name+" does not take a value.")
			}
			values[name] = []string{"true"}
			continue
		}
		if mode == AssignedRepeatable {
			items := []string{}
			if inline {
				items = append(items, value)
			}
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") && args[i+1] != "-h" {
				i++
				items = append(items, args[i])
			}
			if len(items) == 0 {
				return UsageError("InvalidParameter", "Option --"+name+" needs key=value arguments.")
			}
			values[name] = append(values[name], items...)
			continue
		}
		if !inline {
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") || args[i+1] == "-h" {
				return UsageError("InvalidParameter", "Option --"+name+" needs a value.")
			}
			i++
			value = args[i]
		}
		values[name] = []string{value}
	}
	for _, name := range positionals {
		next := cmd.GetSubCommand(name)
		if next == nil {
			return UsageError("InvalidCommand", "Unknown command; use inspire --help to list commands.")
		}
		cmd = next
	}
	ctx.command = cmd
	ctx.flags = NewFlagSet()
	ctx.unknownFlags = NewFlagSet()
	for node := cmd; node != nil; node = node.parent {
		for _, f := range node.Flags().Flags() {
			if node == cmd || f.Persistent {
				copy := *f
				copy.assigned = false
				copy.value = ""
				copy.values = nil
				if ctx.flags.Get(f.Name) == nil {
					ctx.flags.Add(&copy)
				}
			}
		}
	}
	for name, items := range values {
		f := ctx.flags.Get(name)
		if f == nil {
			if !cmd.EnableUnknownFlag {
				return UsageError("InvalidFlag", "Unknown option --"+name+".")
			}
			f = &Flag{Name: name, AssignedMode: AssignedOnce, Category: "parameter"}
			ctx.unknownFlags.Add(f)
		}
		f.assigned = true
		f.values = items
		f.value = items[0]
	}
	if format, ok := ctx.flags.GetValue("cli-output"); ok && format != "json" {
		return UsageError("InvalidParameter", "--cli-output supports json.")
	}
	if help || cmd.Run == nil {
		return cmd.showHelp(ctx)
	}
	return cmd.Run(ctx, nil)
}

func (c *Command) showHelp(ctx *Context) error {
	type metadata struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Usage       string      `json:"usage"`
		Commands    []string    `json:"commands,omitempty"`
		Flags       []*Flag     `json:"flags,omitempty"`
		Request     interface{} `json:"request,omitempty"`
	}
	m := metadata{Name: c.getName(), Description: c.Short, Usage: c.Usage, Flags: ctx.Flags().Flags(), Request: c.RequestSchema}
	for _, child := range c.subCommands {
		m.Commands = append(m.Commands, child.Name)
	}
	sort.Strings(m.Commands)
	format, _ := ctx.Flags().GetValue("cli-output")
	if format == "json" || ctx.Getenv("INSPIRE_CLI_AI_MODE") == "1" {
		return WriteJSON(ctx.Stdout(), m)
	}
	if _, err := fmt.Fprintf(ctx.Stdout(), "%s\n\nUsage: %s\n", m.Description, m.Usage); err != nil {
		return err
	}
	if len(m.Commands) > 0 {
		if _, err := fmt.Fprintf(ctx.Stdout(), "\nCommands: %s\n", strings.Join(m.Commands, ", ")); err != nil {
			return err
		}
	}
	for _, f := range m.Flags {
		if _, err := fmt.Fprintf(ctx.Stdout(), "  --%-24s %s\n", f.Name, f.Short); err != nil {
			return err
		}
	}
	return nil
}

func WriteJSON(w io.Writer, value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return UsageError("InvalidOutput", "Result cannot be encoded as JSON.")
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}
