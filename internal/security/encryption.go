package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidKeySize    = errors.New("invalid key size")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrDecryptionFailed  = errors.New("decryption failed")
)

// EncryptionManager handles AES-256 encryption and decryption
type EncryptionManager struct {
	key []byte
}

// NewEncryptionManager creates a new encryption manager with a derived key
func NewEncryptionManager(passphrase string) *EncryptionManager {
	// Derive a 32-byte key from the passphrase using SHA-256
	hash := sha256.Sum256([]byte(passphrase))
	return &EncryptionManager{
		key: hash[:],
	}
}

// NewEncryptionManagerWithKey creates a new encryption manager with a provided key
func NewEncryptionManagerWithKey(key []byte) (*EncryptionManager, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	keyCopy := make([]byte, 32)
	copy(keyCopy, key)

	return &EncryptionManager{
		key: keyCopy,
	}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (em *EncryptionManager) Encrypt(plaintext string) ([]byte, error) {
	if plaintext == "" {
		return nil, nil
	}

	block, err := aes.NewCipher(em.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (em *EncryptionManager) Decrypt(ciphertext []byte) (string, error) {
	if len(ciphertext) == 0 {
		return "", nil
	}

	block, err := aes.NewCipher(em.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// GenerateKey generates a new 32-byte encryption key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// ClearMemory securely clears sensitive data from memory
func ClearMemory(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
