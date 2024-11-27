package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

func ImportTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <json file> <prefix>",
		Short: "Import json to toml",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := os.ReadFile(args[0])
			if err != nil {
				log.Fatal(err)
			}
			data := make(map[string]interface{})
			if err := json.Unmarshal([]byte(from), &data); err != nil {
				log.Fatal(err)
			}
			prefix := args[1]
			newdata := make(map[string]interface{})
			for k, v := range data {
				if prefix != "" {
					newdata[fmt.Sprintf("%s:%s", prefix, k)] = v
				} else {
					newdata[k] = v
				}
			}
			tree, err := toml.TreeFromMap(newdata)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(tree.String())

			return nil
		},
	}
	return cmd
}
