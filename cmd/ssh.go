package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/MinseokOh/toml-cli/toml"
	"github.com/fatih/color"
	lib "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH host management",
	Long:  "Manage SSH hosts, generate configurations, and connect to hosts stored in cmdb",
}

var sshAddCmd = &cobra.Command{
	Use:   "add [host-key]",
	Short: "Add a new SSH host to cmdb",
	Args:  cobra.ExactArgs(1),
	Run:   runSSHAdd,
}

var sshListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all SSH hosts in cmdb",
	Run:   runSSHList,
}

var sshConfigCmd = &cobra.Command{
	Use:   "config [host-key]",
	Short: "Generate SSH config entry for a host",
	Args:  cobra.ExactArgs(1),
	Run:   runSSHConfig,
}

var sshConnectCmd = &cobra.Command{
	Use:   "c [host-key]",
	Short: "Connect to a host using SSH",
	Args:  cobra.ExactArgs(1),
	Run:   runSSHConnect,
}

var sshSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Generate SSH config file for all hosts",
	Long:  "Generate a separate SSH config file (~/.ssh/cmdb_config) containing all hosts. This does not modify your existing ~/.ssh/config file.",
	Run:   runSSHSync,
}

var sshKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "SSH key management",
}

var sshKeyGenerateCmd = &cobra.Command{
	Use:   "generate [host-key]",
	Short: "Generate SSH key pair for a host",
	Args:  cobra.ExactArgs(1),
	Run:   runSSHKeyGenerate,
}

var sshKeyImportPublicCmd = &cobra.Command{
	Use:   "public [host-key] [public-key-file]",
	Short: "Import existing SSH public key for a host",
	Args:  cobra.ExactArgs(2),
	Run:   runSSHKeyImportPublic,
}

var sshKeyImportPrivateCmd = &cobra.Command{
	Use:   "private [host-key] [private-key-file]",
	Short: "Import existing SSH private key for a host",
	Args:  cobra.ExactArgs(2),
	Run:   runSSHKeyImportPrivate,
}

var sshKeyImportBothCmd = &cobra.Command{
	Use:   "both [host-key] [private-key-file]",
	Short: "Import existing SSH private key and generate corresponding public key",
	Args:  cobra.ExactArgs(2),
	Run:   runSSHKeyImportBoth,
}

func init() {
	sshCmd.AddCommand(sshAddCmd)
	sshCmd.AddCommand(sshListCmd)
	sshCmd.AddCommand(sshConfigCmd)
	sshCmd.AddCommand(sshConnectCmd)
	sshCmd.AddCommand(sshSyncCmd)
	sshKeyCmd.AddCommand(sshKeyGenerateCmd)
	sshKeyCmd.AddCommand(sshKeyImportPublicCmd)
	sshKeyCmd.AddCommand(sshKeyImportPrivateCmd)
	sshKeyCmd.AddCommand(sshKeyImportBothCmd)
	sshCmd.AddCommand(sshKeyCmd)
}

// SSHHost represents a host configuration
type SSHHost struct {
	Hostname     string `toml:"hostname"`
	User         string `toml:"user"`
	Port         int    `toml:"port"`
	Password     string `toml:"password"`
	KeyPath      string `toml:"key_path"`
	PrivateKey   string `toml:"private_key"`
	PublicKey    string `toml:"public_key"`
	Description  string `toml:"description"`
	Environment  string `toml:"environment"`
	Tags         string `toml:"tags"`
	ForwardAgent bool   `toml:"forward_agent"`
	ProxyJump    string `toml:"proxy_jump"`
}

func runSSHAdd(cmd *cobra.Command, args []string) {
	hostKey := args[0]

	// Parse namespace and host name
	parts := strings.Split(hostKey, ":")
	if len(parts) < 3 || parts[1] != "host" {
		color.Red("Invalid host key format. Use: namespace:host:name")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	color.Cyan("Adding SSH host: %s", hostKey)

	// Get hostname
	fmt.Print("Hostname (IP or domain): ")
	hostname, _ := reader.ReadString('\n')
	hostname = strings.TrimSpace(hostname)

	// Get user
	fmt.Print("Username (default: root): ")
	userInput, _ := reader.ReadString('\n')
	user := strings.TrimSpace(userInput)
	if user == "" {
		user = "root"
	}

	// Get port
	fmt.Print("Port (default: 22): ")
	portInput, _ := reader.ReadString('\n')
	portStr := strings.TrimSpace(portInput)
	port := 22
	if portStr != "" {
		if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
			color.Red("Invalid port number, using 22")
			port = 22
		}
	}

	// Get description
	fmt.Print("Description (optional): ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// Get environment
	fmt.Print("Environment (dev/staging/prod, optional): ")
	environment, _ := reader.ReadString('\n')
	environment = strings.TrimSpace(environment)

	// Get tags
	fmt.Print("Tags (comma-separated, optional): ")
	tags, _ := reader.ReadString('\n')
	tags = strings.TrimSpace(tags)

	// Ask about authentication method
	fmt.Print("Authentication method (key/password/both, default: key): ")
	authMethodInput, _ := reader.ReadString('\n')
	authMethod := strings.ToLower(strings.TrimSpace(authMethodInput))
	if authMethod == "" {
		authMethod = "key"
	}

	host := &SSHHost{
		Hostname:    hostname,
		User:        user,
		Port:        port,
		Description: description,
		Environment: environment,
		Tags:        tags,
	}

	switch authMethod {
	case "password":
		fmt.Print("Password: ")
		password, _ := reader.ReadString('\n')
		host.Password = strings.TrimSpace(password)
	case "both":
		fmt.Print("Password: ")
		password, _ := reader.ReadString('\n')
		host.Password = strings.TrimSpace(password)

		fmt.Print("Generate SSH key pair? (y/N): ")
		genKeyInput, _ := reader.ReadString('\n')
		genKey := strings.ToLower(strings.TrimSpace(genKeyInput)) == "y"

		if genKey {
			color.Yellow("Generating SSH key pair...")
			if err := generateSSHKeyPair(hostKey, host); err != nil {
				color.Red("Failed to generate SSH key pair: %v", err)
				return
			}
			color.Green("SSH key pair generated successfully")
		} else {
			fmt.Print("Existing private key path (optional): ")
			keyPathInput, _ := reader.ReadString('\n')
			keyPath := strings.TrimSpace(keyPathInput)
			if keyPath != "" {
				host.KeyPath = keyPath
			}
		}
	default: // key
		fmt.Print("Generate SSH key pair? (y/N): ")
		genKeyInput, _ := reader.ReadString('\n')
		genKey := strings.ToLower(strings.TrimSpace(genKeyInput)) == "y"

		if genKey {
			color.Yellow("Generating SSH key pair...")
			if err := generateSSHKeyPair(hostKey, host); err != nil {
				color.Red("Failed to generate SSH key pair: %v", err)
				return
			}
			color.Green("SSH key pair generated successfully")
		} else {
			fmt.Print("Existing private key path (optional): ")
			keyPathInput, _ := reader.ReadString('\n')
			keyPath := strings.TrimSpace(keyPathInput)
			if keyPath != "" {
				host.KeyPath = keyPath
			}
		}
	}

	// Save to cmdb
	if err := saveHostToCMDB(hostKey, *host); err != nil {
		color.Red("Failed to save host to cmdb: %v", err)
		return
	}

	color.Green("Host '%s' added successfully!", hostKey)

	// Provide usage information
	if host.Password != "" {
		color.Cyan("Connect using: cm ssh connect %s", hostKey)
	} else {
		hasPrivateKey := host.KeyPath != "" || host.PrivateKey != ""
		if hasPrivateKey {
			color.Cyan("Connect using: cm ssh connect %s", hostKey)
		}
	}
	fmt.Println("To generate SSH config file: cm ssh sync")
}

func runSSHList(cmd *cobra.Command, args []string) {
	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	color.Cyan("SSH Hosts in cmdb:")
	fmt.Println()

	count := 0
	for _, key := range tomlFile.Keys() {
		if strings.Contains(key, ":host:") {
			count++
			printHostInfo(key, &tomlFile)
		}
	}

	if count == 0 {
		color.Yellow("No SSH hosts found in cmdb")
		return
	}

	fmt.Printf("\nTotal: %d hosts\n", count)
}

func runSSHConfig(cmd *cobra.Command, args []string) {
	hostKey := args[0]

	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	host, err := getHostFromCMDB(hostKey, tomlFile)
	if err != nil {
		color.Red("Failed to get host '%s': %v", hostKey, err)
		return
	}

	config := generateSSHConfigEntry(hostKey, *host)
	fmt.Println(config)
}

func runSSHConnect(cmd *cobra.Command, args []string) {
	hostKey := args[0]

	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	host, err := getHostFromCMDB(hostKey, tomlFile)
	if err != nil {
		color.Red("Failed to get host '%s': %v", hostKey, err)
		return
	}

	// Determine authentication method
	hasPrivateKey := host.KeyPath != "" || host.PrivateKey != ""
	hasPassword := host.Password != ""

	if !hasPrivateKey && !hasPassword {
		color.Red("No authentication method configured for host '%s'", hostKey)
		return
	}

	// Build SSH command
	var cmdExec *exec.Cmd
	var sshArgs []string

	if hasPrivateKey {
		// Key-based authentication - always use -i flag
		sshArgs = []string{}

		// Determine key path
		var keyPath string
		if host.KeyPath != "" {
			keyPath = host.KeyPath
		} else if host.PrivateKey != "" {
			keyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "cm_tmp")
			// Ensure the key file exists
			if _, err := os.Stat(keyPath); os.IsNotExist(err) {
				// Write the key to file if it doesn't exist
				if err := os.WriteFile(keyPath, []byte(host.PrivateKey+"\r"), 0600); err != nil {
					color.Red("Failed to write private key file: %v", err)
					return
				}
			}
		}

		// Add -i flag for private key
		sshArgs = append(sshArgs, "-i", keyPath)

		if host.Port != 22 {
			sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", host.Port))
		}

		if host.ForwardAgent {
			sshArgs = append(sshArgs, "-A")
		}

		if host.ProxyJump != "" {
			sshArgs = append(sshArgs, "-J", host.ProxyJump)
		}

		// Add user@hostname
		sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", host.User, host.Hostname))

		color.Cyan("Connecting to %s using key authentication (%s)...", hostKey, keyPath)
		cmdExec = exec.Command("ssh", sshArgs...)

	} else if hasPassword {
		// Password-based authentication using sshpass
		sshArgs = []string{}

		if host.Port != 22 {
			sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", host.Port))
		}

		if host.ForwardAgent {
			sshArgs = append(sshArgs, "-A")
		}

		if host.ProxyJump != "" {
			sshArgs = append(sshArgs, "-J", host.ProxyJump)
		}

		// Disable strict host key checking for password auth
		sshArgs = append(sshArgs, "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null")

		// Add user@hostname
		sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", host.User, host.Hostname))

		// Check if sshpass is available
		if _, err := exec.LookPath("sshpass"); err != nil {
			color.Red("sshpass is required for password authentication. Install it with: brew install sshpass (macOS) or apt-get install sshpass (Ubuntu)")
			return
		}

		color.Cyan("Connecting to %s using password authentication...", hostKey)
		cmdExec = exec.Command("sshpass", "-p", host.Password, "ssh")
		cmdExec.Args = append(cmdExec.Args, sshArgs...)
	}

	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		color.Red("SSH connection failed: %v", err)
		os.Exit(1)
	}
}

func runSSHSync(cmd *cobra.Command, args []string) {
	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	// Create separate cmdb SSH config file
	sshConfigPath := filepath.Join(os.Getenv("HOME"), ".ssh", "cmdb_config")

	// Create cmdb config file
	file, err := os.OpenFile(sshConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		color.Red("Failed to create cmdb SSH config: %v", err)
		return
	}
	defer file.Close()

	// Write header
	file.WriteString("# Generated by cm - " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	file.WriteString("# This file is automatically included by SSH config\n\n")

	count := 0
	for _, key := range tomlFile.Keys() {
		if strings.Contains(key, ":host:") {
			host, err := getHostFromCMDB(key, tomlFile)
			if err != nil {
				color.Red("Failed to get host '%s': %v", key, err)
				continue
			}

			config := generateSSHConfigEntry(key, *host)
			if _, err := file.WriteString(config + "\n"); err != nil {
				color.Red("Failed to write host '%s' to config: %v", key, err)
				continue
			}
			count++
		}
	}

	// Add Include directive to main SSH config if not already present
	mainSSHConfigPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	includeLine := fmt.Sprintf("Include %s", sshConfigPath)

	if err := addIncludeToSSHConfig(mainSSHConfigPath, includeLine); err != nil {
		color.Red("Failed to add Include directive to SSH config: %v", err)
		return
	}

	color.Green("Synced %d hosts and updated SSH config", count)
	fmt.Printf("Config file: %s\n", sshConfigPath)
	fmt.Println("SSH config now includes this file automatically")
}

func runSSHKeyGenerate(cmd *cobra.Command, args []string) {
	hostKey := args[0]

	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	host, err := getHostFromCMDB(hostKey, tomlFile)
	if err != nil {
		color.Red("Failed to get host '%s': %v", hostKey, err)
		return
	}

	if err := generateSSHKeyPair(hostKey, host); err != nil {
		color.Red("Failed to generate SSH key pair: %v", err)
		return
	}

	// Update host in cmdb
	if err := saveHostToCMDB(hostKey, *host); err != nil {
		color.Red("Failed to update host in cmdb: %v", err)
		return
	}

	color.Green("SSH key pair generated and saved for '%s'", hostKey)
	fmt.Printf("Private key: ~/.ssh/cm_%s\n", hostKey)
	fmt.Printf("Public key:  ~/.ssh/cm_%s.pub\n", hostKey)
}

func runSSHKeyImportPublic(cmd *cobra.Command, args []string) {
	hostKey := args[0]
	publicKeyFile := args[1]

	// Read public key
	content, err := os.ReadFile(publicKeyFile)
	if err != nil {
		color.Red("Failed to read public key file: %v", err)
		return
	}

	publicKey := strings.TrimSpace(string(content))

	// Load existing host
	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	host, err := getHostFromCMDB(hostKey, tomlFile)
	if err != nil {
		color.Red("Failed to get host '%s': %v", hostKey, err)
		return
	}

	// Update public key
	host.PublicKey = publicKey

	// Save back to cmdb
	if err := saveHostToCMDB(hostKey, *host); err != nil {
		color.Red("Failed to update host in cmdb: %v", err)
		return
	}

	color.Green("Public key imported for '%s'", hostKey)
}

func runSSHKeyImportPrivate(cmd *cobra.Command, args []string) {
	hostKey := args[0]
	privateKeyFile := args[1]

	// Read private key
	content, err := os.ReadFile(privateKeyFile)
	if err != nil {
		color.Red("Failed to read private key file: %v", err)
		return
	}

	privateKey := strings.TrimSpace(string(content))

	// Load existing host
	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	host, err := getHostFromCMDB(hostKey, tomlFile)
	if err != nil {
		color.Red("Failed to get host '%s': %v", hostKey, err)
		return
	}

	// Update private key
	host.PrivateKey = privateKey

	// Save back to cmdb
	if err := saveHostToCMDB(hostKey, *host); err != nil {
		color.Red("Failed to update host in cmdb: %v", err)
		return
	}

	color.Green("Private key imported for '%s'", hostKey)
}

func runSSHKeyImportBoth(cmd *cobra.Command, args []string) {
	hostKey := args[0]
	privateKeyFile := args[1]

	// Read private key
	content, err := os.ReadFile(privateKeyFile)
	if err != nil {
		color.Red("Failed to read private key file: %v", err)
		return
	}

	privateKey := strings.TrimSpace(string(content))

	// Generate public key from private key
	publicKey, err := generatePublicKeyFromPrivate(privateKeyFile)
	if err != nil {
		color.Red("Failed to generate public key from private key: %v", err)
		return
	}

	// Load existing host
	tomlFile, err := toml.NewToml(path)
	if err != nil {
		color.Red("Failed to load cmdb file: %v", err)
		return
	}

	host, err := getHostFromCMDB(hostKey, tomlFile)
	if err != nil {
		color.Red("Failed to get host '%s': %v", hostKey, err)
		return
	}

	// Update both keys
	host.PrivateKey = privateKey
	host.PublicKey = publicKey

	// Save back to cmdb
	if err := saveHostToCMDB(hostKey, *host); err != nil {
		color.Red("Failed to update host in cmdb: %v", err)
		return
	}

	color.Green("Private and public keys imported for '%s'", hostKey)
}

// Helper functions

func printHostInfo(hostKey string, tomlFile *toml.Toml) {
	host, err := getHostFromCMDB(hostKey, *tomlFile)
	if err != nil {
		color.Red("Failed to load host '%s': %v", hostKey, err)
		return
	}

	color.New(color.FgCyan).Add(color.Bold).Printf("%s\n", hostKey)
	fmt.Printf("  Host:     %s@%s:%d\n", host.User, host.Hostname, host.Port)

	if host.Description != "" {
		fmt.Printf("  Desc:     %s\n", host.Description)
	}

	if host.Environment != "" {
		fmt.Printf("  Env:      %s\n", host.Environment)
	}

	if host.Tags != "" {
		fmt.Printf("  Tags:     %s\n", host.Tags)
	}

	// Show authentication method
	hasPrivateKey := host.KeyPath != "" || host.PrivateKey != ""
	hasPassword := host.Password != ""

	if hasPrivateKey {
		if host.KeyPath != "" {
			fmt.Printf("  Auth:     Key (%s)\n", host.KeyPath)
		} else {
			fmt.Printf("  Auth:     Key (Inline)\n")
		}
	}

	if hasPassword {
		fmt.Printf("  Auth:     Password\n")
	}

	fmt.Println()
}

func generateSSHConfigEntry(hostKey string, host SSHHost) string {
	var config strings.Builder

	config.WriteString(fmt.Sprintf("Host %s\n", hostKey))
	config.WriteString(fmt.Sprintf("    HostName %s\n", host.Hostname))
	config.WriteString(fmt.Sprintf("    User %s\n", host.User))
	config.WriteString(fmt.Sprintf("    Port %d\n", host.Port))

	// Handle key-based authentication
	hasPrivateKey := host.KeyPath != "" || host.PrivateKey != ""
	hasPassword := host.Password != ""

	if hasPrivateKey {
		if host.KeyPath != "" {
			config.WriteString(fmt.Sprintf("    IdentityFile %s\n", host.KeyPath))
		} else {
			config.WriteString(fmt.Sprintf("    IdentityFile ~/.ssh/cm_%s\n", hostKey))
		}
	}

	// For password authentication, SSH config doesn't support storing passwords
	// Password auth will be handled by sshpass during connection
	if hasPassword && !hasPrivateKey {
		config.WriteString("    # Password authentication (use sshpass or manual entry)\n")
	}

	if host.ForwardAgent {
		config.WriteString("    ForwardAgent yes\n")
	}

	if host.ProxyJump != "" {
		config.WriteString(fmt.Sprintf("    ProxyJump %s\n", host.ProxyJump))
	}

	// Add some nice defaults
	config.WriteString("    StrictHostKeyChecking no\n")
	config.WriteString("    UserKnownHostsFile /dev/null\n")

	return config.String()
}

func generateSSHKeyPair(hostKey string, host *SSHHost) error {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return err
	}

	privateKeyPath := filepath.Join(sshDir, "cm_"+hostKey)
	publicKeyPath := privateKeyPath + ".pub"

	// Generate key pair
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", privateKeyPath, "-N", "")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Read generated keys
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	host.PrivateKey = string(privateKey)
	host.PublicKey = string(publicKey)

	return nil
}

func copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0600)
}

func getHostFromCMDB(hostKey string, tomlFile toml.Toml) (*SSHHost, error) {
	hostData := tomlFile.Get(hostKey)
	if hostData == nil {
		return nil, fmt.Errorf("host '%s' not found in cmdb", hostKey)
	}

	// The Get method returns a *lib.Tree for complex objects
	var hostMap map[string]interface{}

	// Try to convert to map[string]interface{}
	switch v := hostData.(type) {
	case *lib.Tree:
		hostMap = v.ToMap()
	case map[string]interface{}:
		hostMap = v
	default:
		return nil, fmt.Errorf("invalid host data format for '%s': got %T", hostKey, hostData)
	}

	host := &SSHHost{}
	if hostname, ok := hostMap["hostname"].(string); ok {
		host.Hostname = hostname
	}
	if user, ok := hostMap["user"].(string); ok {
		host.User = user
	}
	if port, ok := hostMap["port"].(int64); ok {
		host.Port = int(port)
	}
	if password, ok := hostMap["password"].(string); ok {
		host.Password = password
	}
	if keyPath, ok := hostMap["key_path"].(string); ok {
		host.KeyPath = keyPath
	}
	if privateKey, ok := hostMap["private_key"].(string); ok {
		host.PrivateKey = privateKey
	}
	if publicKey, ok := hostMap["public_key"].(string); ok {
		host.PublicKey = publicKey
	}
	if description, ok := hostMap["description"].(string); ok {
		host.Description = description
	}
	if environment, ok := hostMap["environment"].(string); ok {
		host.Environment = environment
	}
	if tags, ok := hostMap["tags"].(string); ok {
		host.Tags = tags
	}
	if forwardAgent, ok := hostMap["forward_agent"].(bool); ok {
		host.ForwardAgent = forwardAgent
	}
	if proxyJump, ok := hostMap["proxy_jump"].(string); ok {
		host.ProxyJump = proxyJump
	}

	if host.Hostname == "" {
		return nil, fmt.Errorf("hostname is required for host '%s'", hostKey)
	}

	if host.User == "" {
		host.User = "root"
	}

	if host.Port == 0 {
		host.Port = 22
	}

	return host, nil
}

func saveHostToCMDB(hostKey string, host SSHHost) error {
	tomlFile, err := toml.NewToml(path)
	if err != nil {
		return err
	}

	// Prepare host data for saving
	hostMap := make(map[string]interface{})
	hostMap["hostname"] = host.Hostname
	hostMap["user"] = host.User
	hostMap["port"] = host.Port

	if host.Password != "" {
		hostMap["password"] = host.Password
	}
	if host.KeyPath != "" {
		hostMap["key_path"] = host.KeyPath
	}
	if host.PrivateKey != "" {
		hostMap["private_key"] = host.PrivateKey
	}
	if host.PublicKey != "" {
		hostMap["public_key"] = host.PublicKey
	}
	if host.Description != "" {
		hostMap["description"] = host.Description
	}
	if host.Environment != "" {
		hostMap["environment"] = host.Environment
	}
	if host.Tags != "" {
		hostMap["tags"] = host.Tags
	}
	if host.ForwardAgent {
		hostMap["forward_agent"] = true
	}
	if host.ProxyJump != "" {
		hostMap["proxy_jump"] = host.ProxyJump
	}

	// Set the host data - need to handle this differently based on the toml package API
	// Since toml.Set expects (key, attr, value), we'll set each attribute individually
	for attr, value := range hostMap {
		// Convert int values to int64 for TOML compatibility
		if attr == "port" {
			value = int64(value.(int))
		}
		if err := tomlFile.Set(hostKey, attr, value); err != nil {
			return err
		}
	}

	// Save to file
	return tomlFile.Write()
}

func generatePublicKeyFromPrivate(privateKeyFile string) (string, error) {
	// Use ssh-keygen to extract public key from private key
	cmd := exec.Command("ssh-keygen", "-y", "-f", privateKeyFile)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate public key: %v", err)
	}

	publicKey := strings.TrimSpace(string(output))
	if publicKey == "" {
		return "", fmt.Errorf("generated public key is empty")
	}

	return publicKey, nil
}

func addIncludeToSSHConfig(sshConfigPath, includeLine string) error {
	// Read existing config or create new
	var content string
	var fileExists bool

	if data, err := os.ReadFile(sshConfigPath); err == nil {
		content = string(data)
		fileExists = true
	}

	// Check if include line already exists
	if strings.Contains(content, includeLine) {
		// Include already exists, nothing to do
		return nil
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(sshConfigPath), 0700); err != nil {
		return err
	}

	// Add comment and include directive at the top of the file
	newContent := fmt.Sprintf("# Added by cm - SSH host management\n%s\n\n", includeLine)
	if fileExists {
		newContent += content
	}

	return os.WriteFile(sshConfigPath, []byte(newContent), 0600)
}
