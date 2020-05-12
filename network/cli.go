package network

import (
	"os"

	"github.com/micro/cli/v2"
	mcli "github.com/micro/micro/v2/cli"
	clic "github.com/micro/micro/v2/internal/command/cli"
)

// networkCommands for network toplogy routing
func networkCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "connect",
			Usage:  "connect to the network. specify nodes e.g connect ip:port",
			Action: mcli.Print(networkConnect),
		},
		{
			Name:   "connections",
			Usage:  "List the immediate connections to the network",
			Action: mcli.Print(networkConnections),
		},
		{
			Name:   "graph",
			Usage:  "Get the network graph",
			Action: mcli.Print(networkGraph),
		},
		{
			Name:   "nodes",
			Usage:  "List nodes in the network",
			Action: mcli.Print(netNodes),
		},
		{
			Name:   "routes",
			Usage:  "List network routes",
			Action: mcli.Print(netRoutes),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "service",
					Usage: "Filter by service",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "Filter by address",
				},
				&cli.StringFlag{
					Name:  "gateway",
					Usage: "Filter by gateway",
				},
				&cli.StringFlag{
					Name:  "router",
					Usage: "Filter by router",
				},
				&cli.StringFlag{
					Name:  "network",
					Usage: "Filter by network",
				},
			},
		},
		{
			Name:   "services",
			Usage:  "Get the network services",
			Action: mcli.Print(networkServices),
		},
		// TODO: duplicates call. Move so we reuse same stuff.
		{
			Name:   "call",
			Usage:  "Call a service e.g micro call greeter Say.Hello '{\"name\": \"John\"}",
			Action: mcli.Print(netCall),
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
	}
}

// networkDNSCommands for networking routing
func networkDNSCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "advertise",
			Usage: "Advertise a new node to the network",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Address to register for the specified domain",
					EnvVars: []string{"MICRO_NETWORK_DNS_ADVERTISE_ADDRESS"},
				},
				&cli.StringFlag{
					Name:    "domain",
					Usage:   "Domain name to register",
					EnvVars: []string{"MICRO_NETWORK_DNS_ADVERTISE_DOMAIN"},
					Value:   "network.micro.mu",
				},
				&cli.StringFlag{
					Name:    "token",
					Usage:   "Bearer token for the go.micro.network.dns service",
					EnvVars: []string{"MICRO_NETWORK_DNS_ADVERTISE_TOKEN"},
				},
			},
			Action: mcli.Print(netDNSAdvertise),
		},
		{
			Name:  "remove",
			Usage: "Remove a node's record'",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Address to register for the specified domain",
					EnvVars: []string{"MICRO_NETWORK_DNS_REMOVE_ADDRESS"},
				},
				&cli.StringFlag{
					Name:    "domain",
					Usage:   "Domain name to remove",
					EnvVars: []string{"MICRO_NETWORK_DNS_REMOVE_DOMAIN"},
					Value:   "network.micro.mu",
				},
				&cli.StringFlag{
					Name:    "token",
					Usage:   "Bearer token for the go.micro.network.dns service",
					EnvVars: []string{"MICRO_NETWORK_DNS_REMOVE_TOKEN"},
				},
			},
			Action: mcli.Print(netDNSRemove),
		},
		{
			Name:  "resolve",
			Usage: "Remove a record'",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "domain",
					Usage:   "Domain name to resolve",
					EnvVars: []string{"MICRO_NETWORK_DNS_RESOLVE_DOMAIN"},
					Value:   "network.micro.mu",
				},
				&cli.StringFlag{
					Name:    "type",
					Usage:   "Domain name type to resolve",
					EnvVars: []string{"MICRO_NETWORK_DNS_RESOLVE_TYPE"},
					Value:   "A",
				},
				&cli.StringFlag{
					Name:    "token",
					Usage:   "Bearer token for the go.micro.network.dns service",
					EnvVars: []string{"MICRO_NETWORK_DNS_RESOLVE_TOKEN"},
				},
			},
			Action: mcli.Print(netDNSResolve),
		},
	}
}

func networkConnect(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkConnect(c, args)
}

func networkConnections(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkConnections(c)
}

func networkGraph(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkGraph(c)
}

func networkServices(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkServices(c)
}

// netCall calls services through the network
func netCall(c *cli.Context, args []string) ([]byte, error) {
	os.Setenv("MICRO_PROXY", "go.micro.network")
	return clic.CallService(c, args)
}

func netNodes(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkNodes(c)
}

func netRoutes(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkRoutes(c)
}

func netDNSAdvertise(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkDNSAdvertise(c)
}

func netDNSRemove(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkDNSRemove(c)
}

func netDNSResolve(c *cli.Context, args []string) ([]byte, error) {
	return clic.NetworkDNSResolve(c)
}
