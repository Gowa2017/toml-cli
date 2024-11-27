package cmd

import (
	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

func ClearTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear key ",
		Short: "remove key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			if err := toml.Clear(args[0]); err != nil {
				return err
			}
			if err := toml.Write(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
