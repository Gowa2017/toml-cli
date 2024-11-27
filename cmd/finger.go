package cmd

import (
	"fmt"
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
			query := args[0]

			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}

			for _, k := range toml.Keys() {
				if strings.Index(k, query) == 0 {
					fmt.Println(k)
					fmt.Println(toml.Get(k))
					fmt.Println(strings.Repeat("-", 30))
				}
			}
			return nil
		},
	}

	return cmd
}
