package cmd

import (
	"fmt"
	"os"

	"github.com/MinseokOh/toml-cli/encrypt"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptPassword string

var decryptCmd = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypt an encrypted TOML file",
	Long: `Decrypt an encrypted TOML file using the provided password.
This will convert the file back to plain TOML format.

Example:
  cm decrypt config.toml
  cm decrypt --password mypassword config.toml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n", filePath)
			os.Exit(1)
		}

		// Read the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		// Check if file is encrypted
		if !encrypt.IsEncrypted(data) {
			fmt.Println("File is not encrypted")
			os.Exit(1)
		}

		// Get password from flag or prompt
		var password string
		if decryptPassword != "" {
			password = decryptPassword
		} else {
			var promptErr error
			password, promptErr = encrypt.PromptPassword(false)
			if promptErr != nil {
				fmt.Printf("Error getting password: %v\n", promptErr)
				os.Exit(1)
			}
		}

		// Decrypt the file content
		decryptedContent, err := encrypt.Decrypt(string(data), password)
		if err != nil {
			fmt.Printf("Error decrypting file: %v\n", err)
			os.Exit(1)
		}

		// Write decrypted content back to file
		err = os.WriteFile(filePath, decryptedContent, 0644)
		if err != nil {
			fmt.Printf("Error writing decrypted file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("File '%s' has been decrypted successfully\n", filePath)
	},
}

// GetDecryptCommand returns the decrypt command
func GetDecryptCommand() *cobra.Command {
	return decryptCmd
}

func init() {
	decryptCmd.Flags().StringVarP(&decryptPassword, "password", "p", "", "Password for decryption")
}

