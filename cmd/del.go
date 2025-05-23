package cmd

import (
	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

func DeleteTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "del key attr",
		Aliases: []string{"d"},
		Short:   "Delete a key's attr",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			if err := toml.Delete(args[0], args[1]); err != nil {
				return err
			}
			if err := toml.Write(); err != nil {
				return err
			}
			printAConfigure(args[0], toml.Get(args[0]))
			return nil
		},
	}
	return cmd
}
