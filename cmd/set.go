package cmd

import (
	"strconv"

	"github.com/MinseokOh/toml-cli/toml"
	lib "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

const (
	flagOut = "out"
)

// SetTomlCommand returns set command
func SetTomlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set key attr value [attr1 value1]",
		Short:   "Edit the file to set some data",
		Aliases: []string{"s"},
		Long: `
e.g.
cm set  192.168.11.11 title 123456

e.g.
cm set  192.168.11.11 title 123456 comment 测试主机 -o out.toml
`,
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			outDir, err := cmd.Flags().GetString(flagOut)
			if err != nil {
				return err
			}
			toml, err := toml.NewToml(path)
			if err != nil {
				return err
			}

			toml.Out(outDir)

			for i := 1; i < len(args); i += 2 {
				if err := toml.Set(key, args[i], args[i+1]); err != nil {
					return err
				}
			}

			if err := toml.Write(); err != nil {
				return err
			}
            printAConfigure(key, toml.Get(key))

			return nil
		},
	}

	cmd.Flags().StringP(flagOut, "o", "", "set output directory")
	return cmd
}

func parseInput(str string) interface{} {
	if val, err := strconv.ParseBool(str); err == nil {
		return val
	}

	if val, err := strconv.ParseInt(str, 0, 64); err == nil {
		return val
	}

	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val
	}

	if val, err := lib.ParseLocalDate(str); err == nil {
		return val
	}

	if val, err := lib.ParseLocalDateTime(str); err == nil {
		return val
	}

	if val, err := lib.ParseLocalTime(str); err == nil {
		return val
	}

	return str
}
