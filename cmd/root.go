package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/fatih/color"
	lib "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"os"
)

var path string
var plain bool

var rootCmd = GetRootCommand()

// GetRootCommand returns root command
func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "cm",
		Short:        "cm",
		SilenceUsage: true,
		Long: `A simple CLI for editing and querying TOML files. We use it as a config manager.
	`,
	}

	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "config", "c", "", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&plain, "plain", "p", false, "是否解析密文信息")
	rootCmd.AddCommand(GetTomlCommand())
	rootCmd.AddCommand(SetTomlCommand())
	rootCmd.AddCommand(ListTomlCommand())
	rootCmd.AddCommand(DeleteTomlCommand())
	rootCmd.AddCommand(DumpTomlCommand())
	rootCmd.AddCommand(ImportTomlCommand())
	rootCmd.AddCommand(ClearTomlCommand())
	rootCmd.AddCommand(FingerTomlCommand())
	rootCmd.AddCommand(NamespaceTomlCommand())
	rootCmd.AddCommand(RenameTomlCommand())
	rootCmd.AddCommand(ScanTomlCommand())
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(GetEncryptCommand())
	rootCmd.AddCommand(GetDecryptCommand())
}

// Execute commands
func Execute() {
	home := os.Getenv("HOME")
	if path == "" {
		path = fmt.Sprintf("%s/.config/cmdb/cmdb.toml", home)
		fmt.Printf("配置文件未指定，使用默认文件: %s\n", path)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(fmt.Sprintf("%s/.config/cmdb", home), 0700); err != nil {
			log.Fatalf("Create cmdb dir %s/.config/cmdb failed: %s", home, err)
		}
		f, err := os.OpenFile(path, os.O_CREATE, 0700)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func printAConfigure(key string, v any) {
	color.New(color.FgRed).Add(color.Bold).Add(color.Underline).Printf("%s\n", key)
	switch v.(type) {
	case *lib.Tree:
		tree := v.(*lib.Tree)
		keys := make([]string, 0, len(tree.ToMap()))
		for k := range tree.ToMap() {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Calculate the maximum key length for consistent width
		maxKeyLength := 0
		treeMap := tree.ToMap()
		for _, k := range keys {
			if len(k) > maxKeyLength {
				maxKeyLength = len(k)
			}
			if m, ok := treeMap[k].(map[string]any); ok {
				for kk := range m {
					if len(kk) > maxKeyLength {
						maxKeyLength = len(kk)
					}
				}
			}
		}

		// Add some padding for better readability
		maxKeyLength += 2

		for _, k := range keys {
			v = treeMap[k]
			if s, ok := v.(string); ok {
				if k == "private_key" && !plain {
					fmt.Printf("%-*s = %s\n", maxKeyLength, k, "********************")
				} else {
					fmt.Printf("%-*s = %s\n", maxKeyLength, k, s)
				}
			} else if m, ok := v.(map[string]any); ok {
				fmt.Printf("%s:\n", k)
				for kk, vv := range m {
					fmt.Printf("  %-*s = %s\n", maxKeyLength, kk, vv)
				}
			} else {
				// Handle other types (bool, int, etc.)
				fmt.Printf("%-*s = %v\n", maxKeyLength, k, v)
			}
		}
	default:
		fmt.Println(v)
	}
	color.Blue(strings.Repeat("-", 50))
}
