package cmd

import (
	"fmt"
	"strings"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

// GetTomlCommand returns get command
func NamespaceTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ns",
		Short: "show namespace info",
		Long: `
e.g.
toml-cli ns
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}

			ns := make(map[string]int)

			for _, k := range toml.Keys() {
				ns[strings.Split(k, ":")[0]]++
			}
			fmt.Printf("%10s: %s\n", "namespace", "config number")
			for k, v := range ns {
				fmt.Printf("%10s: %d\n", k, v)
			}
			return nil
		},
	}

	return cmd
}
