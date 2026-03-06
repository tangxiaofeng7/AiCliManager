package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/user"
)

// deriveKey 从机器指纹（主机名 + 用户名）派生 AES-256 密钥
func deriveKey() ([]byte, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}
	currentUser, err := user.Current()
	username := "unknown-user"
	if err == nil {
		username = currentUser.Username
	}
	// 用 sha256 对主机名+用户名做哈希，得到 32 字节密钥
	raw := fmt.Sprintf("%s|%s|aiclimgr-salt-v1", hostname, username)
	hash := sha256.Sum256([]byte(raw))
	return hash[:], nil
}

// Encrypt 使用 AES-256-GCM 加密明文，返回 base64 编码的密文
// 密文格式：nonce(12字节) + ciphertext，整体 base64 编码
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := deriveKey()
	if err != nil {
		return "", fmt.Errorf("派生加密密钥失败: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 失败: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成随机 nonce 失败: %w", err)
	}

	// 加密：nonce 拼接在密文前
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密 Encrypt 产生的 base64 密文，返回明文
func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	key, err := deriveKey()
	if err != nil {
		return "", fmt.Errorf("派生解密密钥失败: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 解码失败: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 失败: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("密文长度不足，数据可能已损坏")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", fmt.Errorf("AES-GCM 解密失败: %w", err)
	}

	return string(plaintext), nil
}

// MaskApiKey 对 API Key 进行脱敏处理，仅保留前4位和后4位
func MaskApiKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
