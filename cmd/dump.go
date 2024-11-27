package cmd

import (
	"log"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/spf13/cobra"
)

func DumpTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump <json | toml>",
		Short: "Convert to json or toml string",
		RunE: func(cmd *cobra.Command, args []string) error {
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}
			var r string
			if len(args) == 0 {
				r, err = toml.ToToml()
			} else {
				if args[0] == "json" {
					r, err = toml.ToJson()
				} else {
					r, err = toml.ToToml()
				}
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(r))
			return nil
		},
	}
	return cmd
}
