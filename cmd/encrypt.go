package cmd

import (
	"fmt"
	"os"

	"github.com/MinseokOh/toml-cli/encrypt"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptPassword string

var encryptCmd = &cobra.Command{
	Use:   "encrypt [file]",
	Short: "Encrypt a TOML file",
	Long: `Encrypt a TOML file using AES-256-GCM encryption.
The encrypted file can only be accessed with the correct password.

Example:
  cm encrypt config.toml
  cm encrypt --password mypassword config.toml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n", filePath)
			os.Exit(1)
		}

		// Check if file is already encrypted
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		if encrypt.IsEncrypted(data) {
			fmt.Println("File is already encrypted")
			os.Exit(1)
		}

		// Get password from flag or prompt
		var password string
		if encryptPassword != "" {
			password = encryptPassword
		} else {
			var promptErr error
			password, promptErr = encrypt.PromptPassword(true)
			if promptErr != nil {
				fmt.Printf("Error getting password: %v\n", promptErr)
				os.Exit(1)
			}
		}

		// Encrypt the file content
		encryptedContent, err := encrypt.Encrypt(data, password)
		if err != nil {
			fmt.Printf("Error encrypting file: %v\n", err)
			os.Exit(1)
		}

		// Save password for future use
		if err := encrypt.SavePassword(password); err != nil {
			fmt.Printf("Warning: failed to save password: %v\n", err)
		}

		// Write encrypted content back to file
		err = os.WriteFile(filePath, []byte(encryptedContent), 0644)
		if err != nil {
			fmt.Printf("Error writing encrypted file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("File '%s' has been encrypted successfully\n", filePath)
	},
}

// GetEncryptCommand returns the encrypt command
func GetEncryptCommand() *cobra.Command {
	return encryptCmd
}

func init() {
	encryptCmd.Flags().StringVarP(&encryptPassword, "password", "p", "", "Password for encryption")
}

