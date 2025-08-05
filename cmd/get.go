package cmd

import (
	"fmt"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

// GetTomlCommand returns get command
func GetTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get query",
		Aliases: []string{"q"},
		Short:   "Get info of key `query`",
		Long: `
e.g.
cm get title
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			res := toml.Get(query)
			if res == nil {
				return fmt.Errorf("Key %v does not exist in %v", query, path)
			}

			printAConfigure(query, res)
			return nil
		},
	}

	return cmd
}
