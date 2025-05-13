package cmd

import (
	"fmt"
	"sort"
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
cm ns
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
			keys := make([]string, 0, len(ns))
			for k := range ns {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			fmt.Printf("%10s: %s\n", "namespace", "config number")
			for _, k := range keys {
				fmt.Printf("%10s: %d\n", k, ns[k])
			}
			return nil
		},
	}

	return cmd
}
