package cmd

import (
	"fmt"
	"sort"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

func ListTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [query]",
		Aliases: []string{"l", "ls"},
		Short:   "List keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			res := toml.List(query)
			sort.Strings(res)
			for i, k := range res {
				fmt.Printf("%4d\t%s\n", i+1, k)
			}
			return nil
		},
	}
	return cmd
}
