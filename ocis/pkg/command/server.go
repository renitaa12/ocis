//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/runtime"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "Start fullstack server",
		Category: "Fullstack",
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg)
		},
		Action: func(c *cli.Context) error {

			cfg.Commons = &shared.Commons{
				Log: &cfg.Log,
			}

			r := runtime.New(cfg)
			return r.Start()
		},
	}
}

func init() {
	register.AddCommand(Server)
}
