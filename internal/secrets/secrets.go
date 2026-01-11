package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// SecretStore interface for secret storage backends
type SecretStore interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	List() ([]string, error)
}

// KeychainStore uses the OS keychain for secret storage
type KeychainStore struct {
	service string
}

// FileStore uses an encrypted file for secret storage
type FileStore struct {
	filePath   string
	passphrase string
	secrets    map[string]string
}

// SecretsManager manages API keys and other secrets
type SecretsManager struct {
	store    SecretStore
	cacheDir string
}

// NewSecretsManager creates a new secrets manager
func NewSecretsManager() (*SecretsManager, error) {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".viki", "secrets")

	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create secrets directory: %w", err)
	}

	// Try to use keychain first, fall back to encrypted file
	var store SecretStore
	var err error

	if isKeychainAvailable() {
		store = NewKeychainStore("viki")
	} else {
		store, err = NewFileStore(filepath.Join(cacheDir, "vault.enc"))
		if err != nil {
			return nil, err
		}
	}

	return &SecretsManager{
		store:    store,
		cacheDir: cacheDir,
	}, nil
}

// SetAPIKey stores an API key for a provider
func (sm *SecretsManager) SetAPIKey(provider, apiKey string) error {
	key := fmt.Sprintf("apikey_%s", provider)
	return sm.store.Set(key, apiKey)
}

// GetAPIKey retrieves an API key for a provider
func (sm *SecretsManager) GetAPIKey(provider string) (string, error) {
	// Check environment first
	envKey := fmt.Sprintf("VIKI_%s_API_KEY", strings.ToUpper(provider))
	if key := os.Getenv(envKey); key != "" {
		return key, nil
	}

	// Fall back to SDD_API_KEY for compatibility
	if key := os.Getenv("SDD_API_KEY"); key != "" {
		return key, nil
	}

	// Check store
	key := fmt.Sprintf("apikey_%s", provider)
	return sm.store.Get(key)
}

// DeleteAPIKey removes an API key
func (sm *SecretsManager) DeleteAPIKey(provider string) error {
	key := fmt.Sprintf("apikey_%s", provider)
	return sm.store.Delete(key)
}

// ListProviders returns all providers with stored API keys
func (sm *SecretsManager) ListProviders() ([]string, error) {
	keys, err := sm.store.List()
	if err != nil {
		return nil, err
	}

	var providers []string
	for _, key := range keys {
		if strings.HasPrefix(key, "apikey_") {
			providers = append(providers, strings.TrimPrefix(key, "apikey_"))
		}
	}

	return providers, nil
}

// SetSecret stores a generic secret
func (sm *SecretsManager) SetSecret(key, value string) error {
	return sm.store.Set(key, value)
}

// GetSecret retrieves a generic secret
func (sm *SecretsManager) GetSecret(key string) (string, error) {
	return sm.store.Get(key)
}

// Keychain implementation

// NewKeychainStore creates a new keychain store
func NewKeychainStore(service string) *KeychainStore {
	return &KeychainStore{service: service}
}

func (k *KeychainStore) Set(key, value string) error {
	switch runtime.GOOS {
	case "darwin":
		return k.setMacOS(key, value)
	case "linux":
		return k.setLinux(key, value)
	case "windows":
		return k.setWindows(key, value)
	default:
		return fmt.Errorf("unsupported OS for keychain: %s", runtime.GOOS)
	}
}

func (k *KeychainStore) Get(key string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return k.getMacOS(key)
	case "linux":
		return k.getLinux(key)
	case "windows":
		return k.getWindows(key)
	default:
		return "", fmt.Errorf("unsupported OS for keychain: %s", runtime.GOOS)
	}
}

func (k *KeychainStore) Delete(key string) error {
	switch runtime.GOOS {
	case "darwin":
		return k.deleteMacOS(key)
	case "linux":
		return k.deleteLinux(key)
	case "windows":
		return k.deleteWindows(key)
	default:
		return fmt.Errorf("unsupported OS for keychain: %s", runtime.GOOS)
	}
}

func (k *KeychainStore) List() ([]string, error) {
	// Listing is not easily supported across all keychains
	// Return empty list and let caller try known keys
	return []string{}, nil
}

// macOS keychain using security command
func (k *KeychainStore) setMacOS(key, value string) error {
	// Delete existing entry first
	k.deleteMacOS(key)

	cmd := exec.Command("security", "add-generic-password",
		"-a", key,
		"-s", k.service,
		"-w", value)
	return cmd.Run()
}

func (k *KeychainStore) getMacOS(key string) (string, error) {
	cmd := exec.Command("security", "find-generic-password",
		"-a", key,
		"-s", k.service,
		"-w")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return strings.TrimSpace(string(out)), nil
}

func (k *KeychainStore) deleteMacOS(key string) error {
	cmd := exec.Command("security", "delete-generic-password",
		"-a", key,
		"-s", k.service)
	cmd.Run() // Ignore errors (key might not exist)
	return nil
}

// Linux using secret-tool (libsecret)
func (k *KeychainStore) setLinux(key, value string) error {
	cmd := exec.Command("secret-tool", "store",
		"--label", k.service+" "+key,
		"service", k.service,
		"key", key)
	cmd.Stdin = strings.NewReader(value)
	return cmd.Run()
}

func (k *KeychainStore) getLinux(key string) (string, error) {
	cmd := exec.Command("secret-tool", "lookup",
		"service", k.service,
		"key", key)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return strings.TrimSpace(string(out)), nil
}

func (k *KeychainStore) deleteLinux(key string) error {
	cmd := exec.Command("secret-tool", "clear",
		"service", k.service,
		"key", key)
	return cmd.Run()
}

// Windows using Windows Credential Manager
func (k *KeychainStore) setWindows(key, value string) error {
	target := k.service + ":" + key
	cmd := exec.Command("cmdkey", "/generic:"+target, "/user:"+key, "/pass:"+value)
	return cmd.Run()
}

func (k *KeychainStore) getWindows(key string) (string, error) {
	// Windows credential manager doesn't have a simple CLI to retrieve passwords
	// Would need to use PowerShell or Windows APIs
	return "", fmt.Errorf("Windows credential retrieval not fully implemented")
}

func (k *KeychainStore) deleteWindows(key string) error {
	target := k.service + ":" + key
	cmd := exec.Command("cmdkey", "/delete:"+target)
	return cmd.Run()
}

// File-based encrypted storage

// NewFileStore creates a new encrypted file store
func NewFileStore(filePath string) (*FileStore, error) {
	fs := &FileStore{
		filePath: filePath,
		secrets:  make(map[string]string),
	}

	// Get passphrase from environment or generate one
	passphrase := os.Getenv("VIKI_VAULT_PASSPHRASE")
	if passphrase == "" {
		// Use machine-specific key
		passphrase = getMachineKey()
	}
	fs.passphrase = passphrase

	// Load existing secrets
	if err := fs.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStore) Set(key, value string) error {
	fs.secrets[key] = value
	return fs.save()
}

func (fs *FileStore) Get(key string) (string, error) {
	value, ok := fs.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return value, nil
}

func (fs *FileStore) Delete(key string) error {
	delete(fs.secrets, key)
	return fs.save()
}

func (fs *FileStore) List() ([]string, error) {
	var keys []string
	for key := range fs.secrets {
		keys = append(keys, key)
	}
	return keys, nil
}

func (fs *FileStore) load() error {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return err
	}

	decrypted, err := decrypt(data, fs.passphrase)
	if err != nil {
		return fmt.Errorf("failed to decrypt vault: %w", err)
	}

	return json.Unmarshal(decrypted, &fs.secrets)
}

func (fs *FileStore) save() error {
	data, err := json.Marshal(fs.secrets)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data, fs.passphrase)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	return os.WriteFile(fs.filePath, encrypted, 0600)
}

// Encryption helpers

func encrypt(data []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func getMachineKey() string {
	// Generate a machine-specific key
	hostname, _ := os.Hostname()
	homeDir, _ := os.UserHomeDir()
	return fmt.Sprintf("%s:%s:viki-vault", hostname, homeDir)
}

func isKeychainAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		_, err := exec.LookPath("security")
		return err == nil
	case "linux":
		_, err := exec.LookPath("secret-tool")
		return err == nil
	case "windows":
		_, err := exec.LookPath("cmdkey")
		return err == nil
	default:
		return false
	}
}
