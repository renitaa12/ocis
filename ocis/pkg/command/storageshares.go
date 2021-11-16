//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/urfave/cli/v2"
)

// StorageSharesCommand is the entrypoint for the storage-shares command.
func StorageSharesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-shares",
		Usage:    "Start storage and data provider for /home/Shares mount",
		Category: "Extensions",
		Flags:    flagset.StorageSharesWithConfig(cfg.Storage),
		Action: func(c *cli.Context) error {
			origCmd := command.StorageShares(configureStorageShares(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureStorageShares(cfg *config.Config) *svcconfig.Config {
	cfg.Storage.Log.Level = cfg.Log.Level
	cfg.Storage.Log.Pretty = cfg.Log.Pretty
	cfg.Storage.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Storage.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Storage.Tracing.Type = cfg.Tracing.Type
		cfg.Storage.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Storage.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Storage
}

func init() {
	register.AddCommand(StorageSharesCommand)
}
