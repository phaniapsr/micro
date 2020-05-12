// Package cli is a command line interface
package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/micro/cli/v2"

	"github.com/chzyer/readline"
)

var (
	prompt = "micro> "

	interactiveCmds = map[string]*cli.Command{}
	// commands = map[string]*command{
	// 	"quit":       {"quit", "Exit the CLI", quit},
	// 	"exit":       {"exit", "Exit the CLI", quit},
	// 	"call":       {"call", "Call a service", callService},
	// 	"list":       {"list", "List services, peers or routes", list},
	// 	"get":        {"get", "Get service info", getService},
	// 	"stream":     {"stream", "Stream a call to a service", streamService},
	// 	"publish":    {"publish", "Publish a message to a topic", publish},
	// 	"health":     {"health", "Get service health", queryHealth},
	// 	"stats":      {"stats", "Get service stats", queryStats},
	// 	"register":   {"register", "Register a service", registerService},
	// 	"deregister": {"deregister", "Deregister a service", deregisterService},
	// }
)

// type command struct {
// 	name  string
// 	usage string
// 	exec  exec
// }

// RegisterInteractiveCommands registers a new command to be added to the list of available commands in interactive mode
func RegisterInteractiveCommands(cmds ...*cli.Command) error {
	for _, cmd := range cmds {
		// if _, ok := interactiveCmds[cmd.Name]; ok {
		// 	return fmt.Errorf("Command %s already registered", cmd.Name)
		// }
		interactiveCmds[cmd.Name] = cmd
	}
	return nil
}

func runc(c *cli.Context) error {
	// commands["help"] = &command{"help", "CLI usage", help}
	alias := map[string]string{
		"?":  "help",
		"ls": "list",
	}

	r, err := readline.New(prompt)
	if err != nil {
		// TODO return err
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}
	defer r.Close()

	for {
		args, err := r.Readline()
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			return err
		}

		args = strings.TrimSpace(args)

		// skip no args
		if len(args) == 0 {
			continue
		}

		parts := strings.Split(args, " ")
		if len(parts) == 0 {
			continue
		}

		name := parts[0]

		// get alias
		if n, ok := alias[name]; ok {
			name = n
		}

		if cmd, ok := interactiveCmds[name]; ok {
			fs := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
			for _, f := range cmd.Flags {
				if err := f.Apply(fs); err != nil {
					return err
				}
			}

			fs.Parse(parts)
			newCtx := cli.NewContext(c.App, fs, c)
			err := cmd.Run(newCtx)
			if err != nil {
				// TODO return err
				println(err.Error())
				continue
			}
		} else {
			// TODO return err
			println("unknown command")
		}
	}
	return nil
}

func RegistryCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "list",
			Usage: "List items in registry or network",
			Subcommands: []*cli.Command{
				{
					Name:   "nodes",
					Usage:  "List nodes in the network",
					Action: Print(netNodes),
				},
				{
					Name:   "routes",
					Usage:  "List network routes",
					Action: Print(netRoutes),
				},
				{
					Name:   "services",
					Usage:  "List services in registry",
					Action: Print(listServices),
				},
			},
		},
		{
			Name:  "register",
			Usage: "Register an item in the registry",
			Subcommands: []*cli.Command{
				{
					Name:   "service",
					Usage:  "Register a service with JSON definition",
					Action: Print(registerService),
				},
			},
		},
		{
			Name:  "deregister",
			Usage: "Deregister an item in the registry",
			Subcommands: []*cli.Command{
				{
					Name:   "service",
					Usage:  "Deregister a service with JSON definition",
					Action: Print(deregisterService),
				},
			},
		},
		{
			Name:  "get",
			Usage: "Get item from registry",
			Subcommands: []*cli.Command{
				{
					Name:   "service",
					Usage:  "Get service from registry",
					Action: Print(getService),
				},
			},
		},
	}
}

//Commands for micro calling action
func Commands() []*cli.Command {
	commands := []*cli.Command{
		{
			Name:   "cli",
			Usage:  "Run the interactive CLI",
			Action: runc,
		},
		{
			Name:   "call",
			Usage:  "Call a service e.g micro call greeter Say.Hello '{\"name\": \"John\"}",
			Action: Print(callService),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Set the address of the service instance to call",
					EnvVars: []string{"MICRO_ADDRESS"},
				},
				&cli.StringFlag{
					Name:    "output, o",
					Usage:   "Set the output format; json (default), raw",
					EnvVars: []string{"MICRO_OUTPUT"},
				},
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
			},
		},
		{
			Name:   "stream",
			Usage:  "Create a service stream",
			Action: Print(streamService),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output, o",
					Usage:   "Set the output format; json (default), raw",
					EnvVars: []string{"MICRO_OUTPUT"},
				},
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
			},
		},
		{
			Name:   "publish",
			Usage:  "Publish a message to a topic",
			Action: Print(publish),
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
			},
		},
		{
			Name:   "stats",
			Usage:  "Query the stats of a service",
			Action: Print(queryStats),
		},
		{
			Name:   "env",
			Usage:  "Get/set micro cli environment",
			Action: Print(listEnvs),
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Action: Print(getEnv),
				},
				{
					Name:   "set",
					Action: Print(setEnv),
				},
				{
					Name:   "add",
					Action: Print(addEnv),
				},
			},
		},
	}

	return append(commands, RegistryCommands()...)
}
