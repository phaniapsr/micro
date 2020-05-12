package store

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/go-micro/v2/store/service/proto"
	mcli "github.com/micro/micro/v2/cli"
	scli "github.com/micro/micro/v2/store/cli"
	"github.com/micro/micro/v2/store/handler"
	"github.com/pkg/errors"
)

var (
	// Name of the store service
	Name = "go.micro.store"
	// Address is the store address
	Address = ":8002"
)

// Run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "store"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
	)

	// the store handler
	storeHandler := &handler.Store{
		Default: *cmd.DefaultOptions().Store,
		Stores:  make(map[string]bool),
	}

	table := "store"
	if v := ctx.String("store_table"); len(v) > 0 {
		table = v
	}

	// set to store table
	storeHandler.Default.Init(
		store.Table(table),
	)

	backend := storeHandler.Default.String()
	options := storeHandler.Default.Options()

	log.Infof("Initialising the [%s] store with opts: %+v", backend, options)

	// set the new store initialiser
	storeHandler.New = func(database string, table string) (store.Store, error) {
		// Record the new database and table in the internal store
		if err := storeHandler.Default.Write(&store.Record{
			Key:   "databases/" + database,
			Value: []byte{},
		}, store.WriteTo("micro", "internal")); err != nil {
			return nil, errors.Wrap(err, "micro store couldn't store new database in internal table")
		}
		if err := storeHandler.Default.Write(&store.Record{
			Key:   "tables/" + database + "/" + table,
			Value: []byte{},
		}, store.WriteTo("micro", "internal")); err != nil {
			return nil, errors.Wrap(err, "micro store couldn't store new table in internal table")
		}

		return storeHandler.Default, nil
	}

	pb.RegisterStoreHandler(service.Server(), storeHandler)

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// Commands is the cli interface for the store service
func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "store",
		Usage: "Run the micro store service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro tunnel address :8002",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
			return nil
		},
		Subcommands: storeCommands(),
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	mcli.RegisterInteractiveCommands(&cli.Command{
		Name:        "store",
		Subcommands: command.Subcommands,
		HelpName:    "store",
	})

	return []*cli.Command{command}
}

//storeCommands for data storing
func storeCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "snapshot",
			Usage:  "Back up a store",
			Action: scli.Snapshot,
			Flags: append(scli.CommonFlags,
				&cli.StringFlag{
					Name:    "destination",
					Usage:   "Backup destination",
					Value:   "file:///tmp/store-snapshot",
					EnvVars: []string{"MICRO_SNAPSHOT_DESTINATION"},
				},
			),
		},
		{
			Name:   "sync",
			Usage:  "Copy all records of one store into another store",
			Action: scli.Sync,
			Flags:  scli.SyncFlags,
		},
		{
			Name:   "restore",
			Usage:  "restore a store snapshot",
			Action: scli.Restore,
			Flags: append(scli.CommonFlags,
				&cli.StringFlag{
					Name:  "source",
					Usage: "Backup source",
					Value: "file:///tmp/store-snapshot",
				},
			),
		},
		{
			Name:   "databases",
			Usage:  "List all databases known to the store service",
			Action: scli.Databases,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "store",
					Usage: "store service to call",
					Value: "go.micro.store",
				},
			},
		},
		{
			Name:   "tables",
			Usage:  "List all tables in the specified database known to the store service",
			Action: scli.Tables,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "store",
					Usage: "store service to call",
					Value: "go.micro.store",
				},
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d"},
					Usage:   "database to list tables of",
					Value:   "micro",
				},
			},
		},
		{
			Name:      "read",
			Usage:     "read a record from the store",
			UsageText: `micro store read [options] key`,
			Action:    scli.Read,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d"},
					Usage:   "database to write to",
					Value:   "micro",
				},
				&cli.StringFlag{
					Name:    "table",
					Aliases: []string{"t"},
					Usage:   "table to write to",
					Value:   "micro",
				},
				&cli.BoolFlag{
					Name:    "prefix",
					Aliases: []string{"p"},
					Usage:   "read prefix",
					Value:   false,
				},
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "show keys and headers (only values shown by default)",
					Value:   false,
				},
				&cli.StringFlag{
					Name:  "output",
					Usage: "output format (json, table)",
					Value: "table",
				},
			},
		},
		{
			Name:      "list",
			Usage:     "list all keys from a store",
			UsageText: `micro store list [options]`,
			Action:    scli.List,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d"},
					Usage:   "database to list from",
					Value:   "micro",
				},
				&cli.StringFlag{
					Name:    "table",
					Aliases: []string{"t"},
					Usage:   "table to write to",
					Value:   "micro",
				},
				&cli.StringFlag{
					Name:  "output",
					Usage: "output format (json)",
				},
				&cli.BoolFlag{
					Name:    "prefix",
					Aliases: []string{"p"},
					Usage:   "list prefix",
					Value:   false,
				},
				&cli.UintFlag{
					Name:    "limit",
					Aliases: []string{"l"},
					Usage:   "list limit",
				},
				&cli.UintFlag{
					Name:    "offset",
					Aliases: []string{"o"},
					Usage:   "list offset",
				},
			},
		},
		{
			Name:      "write",
			Usage:     "write a record to the store",
			UsageText: `micro store write [options] key value`,
			Action:    scli.Write,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "expiry",
					Aliases: []string{"e"},
					Usage:   "expiry in time.ParseDuration format",
					Value:   "",
				},
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d"},
					Usage:   "database to write to",
					Value:   "micro",
				},
				&cli.StringFlag{
					Name:    "table",
					Aliases: []string{"t"},
					Usage:   "table to write to",
					Value:   "micro",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "delete a key from the store",
			UsageText: `micro store delete [options] key`,
			Action:    scli.Delete,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "database",
					Usage: "database to delete from",
					Value: "micro",
				},
				&cli.StringFlag{
					Name:  "table",
					Usage: "table to delete from",
					Value: "micro",
				},
			},
		},
	}
}
