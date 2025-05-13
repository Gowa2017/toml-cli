package cmd

import (
	"fmt"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

func RenameTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename oldkey newkey",
		Short: "Rename a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ok := args[0]
			nk := args[1]
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			v := toml.Get(ok)
			if v == nil {
				return fmt.Errorf("The key [%s] do not exist", ok)
			}
			if err := toml.Clear(ok); err != nil {
				return fmt.Errorf("Clear key [%s] failed: %s", ok, err)
			}
			if err := toml.Set(nk, "", v); err != nil {
				return fmt.Errorf("Write new key [%s] error: %s", nk, err)
			}
			if err := toml.Write(); err != nil {
				return fmt.Errorf("Save error: %v", err)
			}
			printAConfigure(nk, toml.Get(nk))
			return nil
		},
	}
	return cmd
}
