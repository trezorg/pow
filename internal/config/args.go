package config

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

type CommandType string

const (
	Client CommandType = "client"
	Server CommandType = "server"
)

type Config struct {
	Address  string
	Port     int
	Zeroes   int
	LogLevel string
	Period   time.Duration
	Command  CommandType
}

func New(version string) Config {

	config := Config{}

	app := cli.NewApp()
	app.Version = version
	app.HideHelp = false
	app.HideVersion = false
	app.Authors = []*cli.Author{{
		Name:  "Igor Nemilentsev",
		Email: "trezorg@gmail.com",
	}}
	app.Usage = "POW"
	app.EnableBashCompletion = true
	app.ArgsUsage = "Start client or server"
	app.Description = `
	Start tcp pow server or client
	`
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "address",
			Aliases:     []string{"a"},
			Value:       "",
			Usage:       "Address to connect to or address to listen to",
			Destination: &config.Address,
		},
		&cli.StringFlag{
			Name:        "log",
			Aliases:     []string{"l"},
			Value:       "INFO",
			Usage:       "Log level",
			Destination: &config.LogLevel,
		},
		&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			Value:       8081,
			Usage:       "port to connect to or port to listen to",
			Destination: &config.Port,
		},
		&cli.IntFlag{
			Name:        "zeroes",
			Aliases:     []string{"z"},
			Value:       3,
			Usage:       "Zero bits for hashCash algorithm",
			Destination: &config.Zeroes,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "Start tcp pow client",
			Action: func(c *cli.Context) error {
				period := c.Int("period")
				config.Period = time.Duration(period) * time.Second
				config.Command = Client
				return nil
			},
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "period",
					Usage:   "Period client send requests in seconds",
					Value:   3,
					Aliases: []string{"p"},
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start tcp pow server",
			Action: func(c *cli.Context) error {
				config.Command = Server
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Cannot read cli arguments: %v\n", err)
		os.Exit(1)
	}
	return config

}
