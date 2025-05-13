package cmd

import (
	"reflect"
	"strings"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/fatih/color"
	lib "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

// GetTomlCommand returns get command
func ScanTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scan [query]",
		Aliases: []string{"f"},
		Short:   "Scan entry which key or attr contains `query`",
		Long: `
e.g.
cm scan title
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}

			for _, k := range toml.Keys() {
				pos := strings.Index(strings.ToLower(k), strings.ToLower(query))
				if pos >= 0 {
					ori := k[pos : pos+len(query)]
					nk := strings.ReplaceAll(k, ori, color.RedString(ori))
					printAConfigure(nk, toml.Get(k))
					continue
				}

				v := toml.Get(k)
				if v == nil {
					continue
				}

				var vs string
				tp := reflect.TypeOf(v)
				if tp.String() == "*toml.Tree" {
					vs = v.(*lib.Tree).String()
				} else {
					vs = v.(string)
				}

				pos = strings.Index(strings.ToLower(vs), strings.ToLower(query))
				if pos >= 0 {
					printAConfigure(k, strings.ReplaceAll(vs, vs[pos:pos+len(query)], color.RedString(vs[pos:pos+len(query)])))
					continue
				}

			}
			return nil
		},
	}

	return cmd
}
