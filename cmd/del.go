package cmd

import (
	"fmt"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

func DeleteTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "del key attr",
		Aliases: []string{"d"},
		Short:   "Delete a key's attr",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("参数最少有两个")
			}
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			for _, attr := range args[1:] {
				if err := toml.Delete(args[0], attr); err != nil {
					return err
				}
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
