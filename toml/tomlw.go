package toml

import (
	"fmt"
	"os"

	"github.com/MinseokOh/toml-cli/encrypt"
)

func (t *Toml) readFile() error {
	var err error
	t.raw, err = os.ReadFile(t.path)
	if err != nil {
		return err
	}

	// Check if file is encrypted and decrypt if necessary
	if encrypt.IsEncrypted(t.raw) {
		var password string
		var prompt bool
		// Check if password file exists and is not expired
		passwordData, err := encrypt.ReadPasswordFile()
		if err != nil {
			// No password file or expired, prompt for new password
			fmt.Println("Please enter password to decrypt the cmdb file:")
			newPassword, promptErr := encrypt.PromptPassword(false)
			if promptErr != nil {
				return fmt.Errorf("failed to prompt for password: %w", promptErr)
			}
			password = newPassword
			prompt = true
		} else {
			password = passwordData.Password
		}
		// Decrypt the file content
		t.raw, err = encrypt.Decrypt(string(t.raw), password)
		if err != nil {
			return fmt.Errorf("failed to decrypt file: %w", err)
		}

		if prompt {
			encrypt.SavePassword(password)
		} else {
			// Update last access time
			encrypt.UpdateLastAccess()
		}

	}

	return nil
}

// isFileEncrypted checks if a file is encrypted
func isFileEncrypted(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return encrypt.IsEncrypted(data), nil
}

// Write edited toml tree given path.
// if dest is not setted, overwrite it.
func (t *Toml) Write() error {
	path := t.out
	if path == "" {
		path = t.path
	}

	toml, err := t.tree.ToTomlString()
	if err != nil {
		return err
	}

	// Check if the target file should be encrypted
	shouldEncrypt, err := isFileEncrypted(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file is encrypted: %w", err)
	}

	// If file doesn't exist, default to not encrypting
	if os.IsNotExist(err) {
		shouldEncrypt = false
	}

	var content []byte
	if shouldEncrypt {
		// Get password for encryption
		passwordData, err := encrypt.ReadPasswordFile()
		if err != nil {
			return fmt.Errorf("failed to get password for encryption: %w", err)
		}

		// Encrypt the content
		encryptedContent, err := encrypt.Encrypt([]byte(toml), passwordData.Password)
		if err != nil {
			return fmt.Errorf("failed to encrypt content: %w", err)
		}
		content = []byte(encryptedContent)
	} else {
		content = []byte(toml)
	}

	return os.WriteFile(path, content, 0644)
}
