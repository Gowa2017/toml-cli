package cmd

import (
	"strings"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

// GetTomlCommand returns get command
func FingerTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "finger [query]",
		Aliases: []string{"f"},
		Short:   "show config info",
		Long: `
e.g.
toml-cli f  title
TOML Example
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var query = ""
			if len(args) > 0 {
				query = args[0]
			}

			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}

			for _, k := range toml.Keys() {
				if strings.Contains(strings.ToLower(k), strings.ToLower(query)) {
					printAConfigure(k, toml.Get(k))
				}
			}
			return nil
		},
	}

	return cmd
}
