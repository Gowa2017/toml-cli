package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var path string

var rootCmd = GetRootCommand()

// GetRootCommand returns root command
func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "toml-cli",
		Short:        "toml-cli",
		SilenceUsage: true,
		Long: `A simple CLI for editing and querying TOML files.
	`,
	}

	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "config", "c", "", "配置文件路径")
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

func printAConfigure(k string, v any) {
	fmt.Println(k)
	fmt.Println(v)
	fmt.Println(strings.Repeat("-", 30))
}
