package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func keyBytes(key string) []byte {
	h := sha256.Sum256([]byte(key))
	return h[:]
}

func Encrypt(key, plain string) (string, error) {
	if key == "" {
		return "", errors.New("未配置加密密钥")
	}
	block, err := aes.NewCipher(keyBytes(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func Decrypt(key, enc string) (string, error) {
	if key == "" {
		return "", errors.New("未配置加密密钥")
	}
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(keyBytes(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("密文格式错误")
	}
	nonce, payload := data[:ns], data[ns:]
	plain, err := gcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return "", errors.New("解密失败")
	}
	return string(plain), nil
}
