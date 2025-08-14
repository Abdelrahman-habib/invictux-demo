package security

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"
)

func TestNewEncryptionManager(t *testing.T) {
	passphrase := "test-passphrase-123"
	em := NewEncryptionManager(passphrase)

	if em == nil {
		t.Fatal("Expected encryption manager to be created")
	}

	if len(em.key) != 32 {
		t.Errorf("Expected key length to be 32, got %d", len(em.key))
	}
}

func TestNewEncryptionManagerWithKey(t *testing.T) {
	// Test with valid 32-byte key
	validKey := make([]byte, 32)
	rand.Read(validKey)

	em, err := NewEncryptionManagerWithKey(validKey)
	if err != nil {
		t.Fatalf("Expected no error with valid key, got: %v", err)
	}

	if em == nil {
		t.Fatal("Expected encryption manager to be created")
	}

	// Test with invalid key size
	invalidKey := make([]byte, 16)
	_, err = NewEncryptionManagerWithKey(invalidKey)
	if err != ErrInvalidKeySize {
		t.Errorf("Expected ErrInvalidKeySize, got: %v", err)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	em := NewEncryptionManager("test-passphrase")

	testCases := []struct {
		name      string
		plaintext string
	}{
		{"empty string", ""},
		{"simple text", "hello world"},
		{"special characters", "!@#$%^&*()_+-=[]{}|;:,.<>?"},
		{"unicode", "Hello 世界 🌍"},
		{"long text", strings.Repeat("This is a long text for testing encryption. ", 100)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test encryption
			ciphertext, err := em.Encrypt(tc.plaintext)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Empty plaintext should return nil ciphertext
			if tc.plaintext == "" {
				if ciphertext != nil {
					t.Error("Expected nil ciphertext for empty plaintext")
				}
				return
			}

			// Ciphertext should be different from plaintext
			if bytes.Equal(ciphertext, []byte(tc.plaintext)) {
				t.Error("Ciphertext should not equal plaintext")
			}

			// Test decryption
			decrypted, err := em.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("Decrypted text doesn't match original. Expected: %s, Got: %s", tc.plaintext, decrypted)
			}
		})
	}
}

func TestEncryptDecryptConsistency(t *testing.T) {
	em := NewEncryptionManager("test-passphrase")
	plaintext := "consistent test data"

	// Encrypt the same data multiple times
	ciphertext1, err := em.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	ciphertext2, err := em.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Ciphertexts should be different (due to random nonce)
	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Error("Multiple encryptions of same data should produce different ciphertexts")
	}

	// Both should decrypt to the same plaintext
	decrypted1, err := em.Decrypt(ciphertext1)
	if err != nil {
		t.Fatalf("First decryption failed: %v", err)
	}

	decrypted2, err := em.Decrypt(ciphertext2)
	if err != nil {
		t.Fatalf("Second decryption failed: %v", err)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Error("Decrypted texts don't match original plaintext")
	}
}

func TestDecryptInvalidData(t *testing.T) {
	em := NewEncryptionManager("test-passphrase")

	testCases := []struct {
		name       string
		ciphertext []byte
		expectErr  error
	}{
		{"empty ciphertext", []byte{}, nil},
		{"too short ciphertext", []byte{1, 2, 3}, ErrInvalidCiphertext},
		{"invalid ciphertext", []byte("invalid data that's long enough"), ErrDecryptionFailed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decrypted, err := em.Decrypt(tc.ciphertext)

			if tc.name == "empty ciphertext" {
				if err != nil {
					t.Errorf("Expected no error for empty ciphertext, got: %v", err)
				}
				if decrypted != "" {
					t.Errorf("Expected empty string for empty ciphertext, got: %s", decrypted)
				}
				return
			}

			if err != tc.expectErr {
				t.Errorf("Expected error %v, got: %v", tc.expectErr, err)
			}
		})
	}
}

func TestDifferentKeys(t *testing.T) {
	em1 := NewEncryptionManager("passphrase1")
	em2 := NewEncryptionManager("passphrase2")

	plaintext := "test data"

	// Encrypt with first manager
	ciphertext, err := em1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt with second manager (different key)
	_, err = em2.Decrypt(ciphertext)
	if err != ErrDecryptionFailed {
		t.Errorf("Expected decryption to fail with different key, got: %v", err)
	}
}

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey()
	if err != nil {
		t.Fatalf("Key generation failed: %v", err)
	}

	if len(key1) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key1))
	}

	key2, err := GenerateKey()
	if err != nil {
		t.Fatalf("Second key generation failed: %v", err)
	}

	// Keys should be different
	if bytes.Equal(key1, key2) {
		t.Error("Generated keys should be different")
	}
}

func TestClearMemory(t *testing.T) {
	data := []byte("sensitive data")
	original := make([]byte, len(data))
	copy(original, data)

	ClearMemory(data)

	// Check that data is cleared
	for i, b := range data {
		if b != 0 {
			t.Errorf("Expected byte at index %d to be 0, got %d", i, b)
		}
	}

	// Original should be unchanged
	if bytes.Equal(data, original) {
		t.Error("Data should have been cleared")
	}
}

func BenchmarkEncrypt(b *testing.B) {
	em := NewEncryptionManager("benchmark-passphrase")
	plaintext := "This is a test string for benchmarking encryption performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := em.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	em := NewEncryptionManager("benchmark-passphrase")
	plaintext := "This is a test string for benchmarking decryption performance"

	ciphertext, err := em.Encrypt(plaintext)
	if err != nil {
		b.Fatalf("Setup encryption failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := em.Decrypt(ciphertext)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}
