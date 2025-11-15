package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost bcrypt 哈希成本
	BcryptCost = 12
)

var (
	// ErrInvalidKey 无效的加密密钥
	ErrInvalidKey = errors.New("invalid encryption key")
	// ErrInvalidCiphertext 无效的密文
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	// ErrEncryptionFailed 加密失败
	ErrEncryptionFailed = errors.New("encryption failed")
	// ErrDecryptionFailed 解密失败
	ErrDecryptionFailed = errors.New("decryption failed")
)

// CryptoService 加密服务
type CryptoService struct {
	aesKey []byte
}

// NewCryptoService 创建加密服务实例
func NewCryptoService(key string) (*CryptoService, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}
	return &CryptoService{
		aesKey: []byte(key),
	}, nil
}

// EncryptAES 使用 AES-256-GCM 加密数据
func (c *CryptoService) EncryptAES(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrEncryptionFailed
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES 使用 AES-256-GCM 解密数据
func (c *CryptoService) DecryptAES(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// HashPassword 使用 bcrypt 哈希密码
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码是否匹配哈希值
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
