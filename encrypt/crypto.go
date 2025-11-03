package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
)

const (
	CMDBRCFile    = ".cmdbrc"
	SaltSize      = 16
	NonceSize     = 12
	KeyIterations = 100000
	ExpiryMinutes = 10
)

// PasswordData stores password hash with timestamp
type PasswordData struct {
	Password   string `json:"password_hash"`
	LastAccess int64  `json:"last_access"`
}

// EncryptData represents encrypted file content
type EncryptData struct {
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
	Salt       string `json:"salt"`
}

// deriveKey derives encryption key from password using PBKDF2
func deriveKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, KeyIterations, 32, sha256.New)
}


// generateSalt generates a random salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	return salt, err
}

// generateNonce generates a random nonce for AES-GCM
func generateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	_, err := rand.Read(nonce)
	return nonce, err
}

// Encrypt encrypts data using AES-256-GCM
func Encrypt(data []byte, password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	nonce, err := generateNonce()
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	encryptData := EncryptData{
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Salt:       base64.StdEncoding.EncodeToString(salt),
	}

	jsonData, err := json.Marshal(encryptData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal encrypted data: %w", err)
	}

	return string(jsonData), nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(encryptedData string, password string) ([]byte, error) {
	var encryptData EncryptData
	err := json.Unmarshal([]byte(encryptedData), &encryptData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted data: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptData.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptData.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(encryptData.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GetCMDBRCPath returns the path to .cmdbrc file
func GetCMDBRCPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return CMDBRCFile
	}
	return filepath.Join(home, CMDBRCFile)
}

// SavePassword saves password hash to .cmdbrc file
func SavePassword(password string) error {
	passwordData := PasswordData{
		Password:   password,
		LastAccess: time.Now().Unix(),
	}

	jsonData, err := json.Marshal(passwordData)
	if err != nil {
		return fmt.Errorf("failed to marshal password data: %w", err)
	}

	path := GetCMDBRCPath()
	err = os.WriteFile(path, jsonData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write password file: %w", err)
	}

	return nil
}


// ReadPasswordFile reads password data if file exists and is not expired
func ReadPasswordFile() (*PasswordData, error) {
	path := GetCMDBRCPath()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("password file not found")
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read password file: %w", err)
	}

	var passwordData PasswordData
	err = json.Unmarshal(data, &passwordData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal password data: %w", err)
	}
	lastAccessTime := time.Unix(passwordData.LastAccess, 0)
	if time.Since(lastAccessTime) > time.Duration(10)*time.Minute {
		// Remove expired file
		os.Remove(path)
		return nil, errors.New("password expired")
	}

	return &passwordData, nil
}

// UpdateLastAccess updates the last access time in the password file
func UpdateLastAccess() {
	path := GetCMDBRCPath()

	// Read current data
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var passwordData PasswordData
	if err := json.Unmarshal(data, &passwordData); err != nil {
		return
	}

	// Update last access time
	passwordData.LastAccess = time.Now().Unix()
	jsonData, err := json.Marshal(passwordData)
	if err != nil {
		return
	}

	os.WriteFile(path, jsonData, 0600)
}

// ClearPassword removes the .cmdbrc file
func ClearPassword() error {
	path := GetCMDBRCPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

// IsEncrypted checks if data is encrypted by looking for JSON structure
func IsEncrypted(data []byte) bool {
	trimmed := strings.TrimSpace(string(data))
	return strings.HasPrefix(trimmed, "{") && strings.Contains(trimmed, "ciphertext")
}

// PromptPassword prompts user to enter password
func PromptPassword(confirm bool) (string, error) {
	fmt.Print("Enter password for cmdb file: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // New line after password input

	if confirm {
		fmt.Print("Confirm password: ")
		confirmPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println() // New line after password input

		if string(password) != string(confirmPassword) {
			return "", errors.New("passwords do not match")
		}
	}

	return string(password), nil
}
