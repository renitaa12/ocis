//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageSharingCommand is the entrypoint for the reva-sharing command.
func StorageSharingCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-sharing",
		Usage:    "Start storage sharing service",
		Category: "Extensions",
		//Flags:    flagset.SharingWithConfig(cfg.Storage),
		Before: func(ctx *cli.Context) error {
			return ParseStorageCommon(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Sharing(cfg.Storage)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageSharingCommand)
}
